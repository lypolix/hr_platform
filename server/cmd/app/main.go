package main

import (
	"log/slog"
	"os"

	"github.com/hr-platform-mosprom/internal"
)

func main() {
	env, err := internal.LoadEnv()
	if err != nil {
		slog.Error("error loading env", "err", err)
		os.Exit(1)
	}
	err = internal.Run(env)
	if err != nil {
		slog.Error("error running app", "err", err)
		os.Exit(1)
	}
}
