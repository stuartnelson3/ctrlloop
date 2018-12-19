package ctrlloop

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Loop(ctx context.Context, interval time.Duration, run, cleanup func() error) error {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		<-sigc
		cancel()
	}()

	t := time.NewTicker(interval)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			if err := run(); err != nil {
				return err
			}
		case <-ctx.Done():
			return cleanup()
		}
	}
}
