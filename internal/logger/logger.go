package logger

import "go.uber.org/zap"

func NewSugarLogger() *zap.SugaredLogger {
	log, err := zap.NewDevelopment()

	if err != nil {
		panic(err)
	}

	sugar := log.Sugar()

	return sugar
}
