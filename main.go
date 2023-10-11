package main

import (
	"log/slog"
	"os"
	"runtime/debug"
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
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	info, ok := debug.ReadBuildInfo()
	if !ok {
		slog.Info("unable to read build info")
	}

	commit, ok := getVCSCommit(info)
	if !ok {
		slog.Info("commit information not found in binary")
	}

	slog.Info("Starting petinfoservice", "commit", commit)
}
