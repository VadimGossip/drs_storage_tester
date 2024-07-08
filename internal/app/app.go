package app

import (
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)
}

type App struct {
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
	//ctx, cancel := context.WithCancel(context.Background())
	//defer cancel()
	dbAdapter := NewDBAdapter()
	if err := dbAdapter.Connect(); err != nil {
		logrus.Fatalf("Fail to connect db %s", err)
	}

	logrus.Infof("[%s] started", app.name)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	logrus.Infof("[%s] got signal: [%s]", app.name, <-c)

	logrus.Infof("[%s] stopped", app.name)
}
