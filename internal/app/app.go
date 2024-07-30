package app

import (
	"context"
	"github.com/VadimGossip/tj-drs-storage/internal/domain"
	"github.com/VadimGossip/tj-drs-storage/pkg/util"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)
}

type App struct {
	*Factory
	name         string
	configDir    string
	appStartedAt time.Time
}

func NewApp(name, configDir string, appStartedAt time.Time) *App {
	return &App{
		name:         name,
		configDir:    configDir,
		appStartedAt: appStartedAt,
	}
}

func setLogFile(filepath string) *os.File {
	if filepath == "" {
		logrus.Info("Empty path to log to file, using default stdout")
		return nil
	}
	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logrus.Info("Failed to log to file, using default stdout")
		return nil
	} else {
		logrus.SetOutput(file)
	}
	return file
}

func (app *App) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	dbAdapter := NewDBAdapter(&domain.Config{TargetDb: domain.TargetDbConfig{
		Host:     "",
		Port:     29130,
		Username: "",
		Password: "",
		Db:       0,
	}},
	)
	if err := dbAdapter.Connect(ctx); err != nil {
		logrus.Fatalf("Fail to connect db %s", err)
	}
	app.Factory = newFactory(dbAdapter)

	logrus.Infof("[%s] started", app.name)
	if err := app.Factory.imitator.RunTests(ctx, &domain.Task{
		RequestsPerSec: 100,
		PackPerSec:     10,
		Summary: &domain.TaskSummary{
			Total:    100000,
			Duration: &domain.DurationSummary{EMA: util.NewEMA(0.01), Histogram: make(map[float64]int)},
		},
	}); err != nil {
		logrus.Errorf("[%s] fail to run tests: %s", app.name, err)
		return
	}

	//c := make(chan os.Signal, 1)
	//signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	//logrus.Infof("[%s] got signal: [%s]", app.name, <-c)
	if err := dbAdapter.Disconnect(); err != nil {
		logrus.Errorf("[%s] fail to diconnect db: %s", app.name, err)
		return
	}

	logrus.Infof("[%s] stopped", app.name)
}
