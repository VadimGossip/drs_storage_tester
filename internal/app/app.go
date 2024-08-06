package app

import (
	"context"
	"fmt"
	"github.com/VadimGossip/tj-drs-storage/internal/closer"
	"os"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/VadimGossip/tj-drs-storage/internal/config"
	"github.com/VadimGossip/tj-drs-storage/internal/domain"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)
}

type App struct {
	serviceProvider *serviceProvider
	name            string
	configDir       string
	appStartedAt    time.Time
	cfg             *domain.Config
}

func NewApp(ctx context.Context, name, configDir string, appStartedAt time.Time) (*App, error) {
	a := &App{
		name:         name,
		configDir:    configDir,
		appStartedAt: appStartedAt,
	}

	if err := a.initDeps(ctx); err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
	}

	for _, f := range inits {
		if err := f(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	cfg, err := config.Init(a.configDir)
	if err != nil {
		return fmt.Errorf("[%s] config initialization error: %s", a.name, err)
	}
	a.cfg = cfg
	logrus.Infof("[%s] got config: [%+v]", a.name, *a.cfg)
	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider(a.cfg)
	return nil
}

func (a *App) Run(ctx context.Context) error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
		logrus.Infof("[%s] stopped", a.name)
	}()
	logrus.Infof("[%s] started", a.name)
	if err := a.serviceProvider.ImitatorService(ctx).RunTests(ctx, a.cfg.Task); err != nil {
		logrus.Errorf("[%s] fail to run tests: %s", a.name, err)
		return err
	}

	return nil
}
