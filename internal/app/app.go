package app

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/VadimGossip/platform_common/pkg/closer"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
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
}

func NewApp(ctx context.Context, name string, appStartedAt time.Time) (*App, error) {
	a := &App{
		name:         name,
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
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

func (a *App) Run(ctx context.Context) error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
		logrus.Infof("[%s] stopped", a.name)
	}()
	logrus.Infof("[%s] started", a.name)
	if err := a.serviceProvider.ImitatorService(ctx).RunTests(ctx); err != nil {
		logrus.Errorf("[%s] fail to run tests: %s", a.name, err)
		return err
	}

	return nil
}
