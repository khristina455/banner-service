package main

import (
	"time"

	"github.com/sirupsen/logrus"

	"banner-service/internal/app"
)

func main() {
	logger := logrus.New()
	formatter := &logrus.TextFormatter{
		TimestampFormat: time.DateTime,
		FullTimestamp:   true,
	}
	logger.SetFormatter(formatter)

	app := app.NewApp(logger)

	if err := app.Run(); err != nil {
		logger.Fatal(err)
	}
}
