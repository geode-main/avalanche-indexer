package rewards

import (
	"context"

	"github.com/sirupsen/logrus"
)

type Worker struct {
	log *logrus.Logger
}

func NewWorker() Worker {
	return Worker{}
}

func (w *Worker) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		}
	}
}
