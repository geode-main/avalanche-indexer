package indexer

import (
	"github.com/sirupsen/logrus"

	"github.com/figment-networks/indexing-engine/pipeline"
)

type logger struct {
	l *logrus.Logger
}

func NewLogger(log *logrus.Logger) pipeline.Logger {
	return logger{log}
}

func (log logger) Info(msg string) {
	log.l.Info(msg)
}

func (log logger) Debug(msg string) {
	log.l.Debug(msg)
}
