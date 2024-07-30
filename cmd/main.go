package main

import (
	"time"

	"github.com/VadimGossip/tj-drs-storage/internal/app"
)

var configDir = "config"

func main() {
	storage := app.NewApp("Db test stand", configDir, time.Now())
	storage.Run()
}
