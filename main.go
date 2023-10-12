package main

import (
	"context"
	"flag"
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

	flagset := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	addrFlag := flagset.String("addr", ":8080", "address server will listen on")
	if err := flagset.Parse(os.Args[1:]); err != nil {
		slog.ErrorContext(ctx, "could not parse command line flags", "error", err)
		return
	}

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

	serv := service.Service{
		Addr: *addrFlag,
	}

	servErr := make(chan error, 1)
	go func() {
		servErr <- serv.Open()
	}()

	select {
	case <-ctx.Done():
		slog.Info("closing server due to interrupt")
		serv.Close()
	case err := <-servErr:
		slog.ErrorContext(ctx, "could not start petinfoservice", "error", err)
	}
}
