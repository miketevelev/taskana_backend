package recurring_worker

import (
	"context"
	"time"

	core_logger "github.com/miketevelev/taskana_backend/internal/core/logger"
	"go.uber.org/zap"
)

type RecurrenceProcessor interface {
	ProcessFixedRecurrences(ctx context.Context, asOf time.Time) error
}

type Worker struct {
	processor RecurrenceProcessor
	log       *core_logger.Logger
	interval  time.Duration
}

func NewWorker(
	processor RecurrenceProcessor,
	log *core_logger.Logger,
	interval time.Duration,
) *Worker {
	return &Worker{
		processor: processor,
		log:       log,
		interval:  interval,
	}
}

func (w *Worker) Run(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	w.runOnce(ctx)

	for {
		select {
		case <-ctx.Done():
			w.log.Info("recurring worker stopped")
			return
		case <-ticker.C:
			w.runOnce(ctx)
		}
	}
}

func (w *Worker) runOnce(ctx context.Context) {
	asOf := time.Now().UTC()
	if err := w.processor.ProcessFixedRecurrences(ctx, asOf); err != nil {
		w.log.Error("process fixed recurrences failed", zap.Error(err))
	}
}
