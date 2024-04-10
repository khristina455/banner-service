package main

import (
	"banner-service/internal/app"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	formatter := &logrus.TextFormatter{
		TimestampFormat: time.DateTime,
		FullTimestamp:   true,
	}
	logger.SetFormatter(formatter)

	app := app.NewApp(logger) //все соднанное сверху предаем сюда

	if err := app.Run(); err != nil {
		logger.Fatal(err)
	}
}
