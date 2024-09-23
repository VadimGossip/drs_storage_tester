package main

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/VadimGossip/drs_storage_tester/internal/app"
)

var (
	appName = "Db test stand"
)

func main() {
	ctx := context.Background()
	a, err := app.NewApp(ctx, appName, time.Now())
	if err != nil {
		logrus.Fatalf("failed to init app[%s]: %s", appName, err)
	}

	if err = a.Run(ctx); err != nil {
		logrus.Infof("app[%s] run process finished with error: %s", appName, err)
	}
}
