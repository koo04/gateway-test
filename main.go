package main

import (
	"log/slog"

	"github.com/koo04/gateway-test/api"
)

func main() {
	if err := api.Start(); err != nil {
		slog.Error("Failed to start the api", "error", err)
		return
	}
}
