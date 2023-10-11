package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	"github.com/lucyquest/petinfoservice/service"
)

func getVCSCommit(b *debug.BuildInfo) (string, bool) {
	for _, s := range b.Settings {
		if s.Key == "vcs.revision" {
			return s.Value, true
		}
	}
	return "", false
}

func main() {
	// Setup a context that is cancelled when we get an interupt
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	go func() {
		osSignal := make(chan os.Signal, 1)
		signal.Notify(osSignal, syscall.SIGTERM, os.Interrupt)

		<-osSignal
		cancelFunc()
	}()

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	info, ok := debug.ReadBuildInfo()
	if !ok {
		slog.InfoContext(ctx, "unable to read build info")
	}

	commit, ok := getVCSCommit(info)
	if !ok {
		slog.InfoContext(ctx, "commit information not found in binary")
	}

	slog.InfoContext(ctx, "Starting petinfoservice", "commit", commit)

	serv := service.Service{}
	if err := serv.Open(); err != nil {
		slog.ErrorContext(ctx, "could not start petinfoservice", "error", err)
		return
	}

	<-ctx.Done()
	serv.Close()
}
