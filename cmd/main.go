package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"time"

	"github.com/VadimGossip/tj-drs-storage/internal/app"
)

var (
	configDir = "config"
	appName   = "Db test stand"
)

func main() {
	ctx := context.Background()
	a, err := app.NewApp(ctx, appName, configDir, time.Now())
	if err != nil {
		logrus.Fatalf("failed to init app[%s]: %s", appName, err)
	}

	if err = a.Run(ctx); err != nil {
		logrus.Infof("app[%s] run process finished with error: %s", appName, err)
	}
}
