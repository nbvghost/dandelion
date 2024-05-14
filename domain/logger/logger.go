package logger

import (
	"github.com/nbvghost/dandelion/library/environments"
	"go.uber.org/zap"
)

func CreateLogger(pName string, traceID string) (*zap.Logger, error) {
	var logger *zap.Logger
	var err error
	if environments.Release() {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment()
	}
	if err != nil {
		return nil, err
	}
	logger = logger.Named(pName).With(zap.String("TraceID", traceID)) //.With(zap.String("DomainName", domainName))
	return logger, err
}
