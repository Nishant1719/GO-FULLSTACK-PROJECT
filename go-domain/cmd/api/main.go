package main

import (
	"log/slog"
	"os"
)

func main() {
	cfg := config{
		addr: ":8080",
		db:   dbConfig{},
	}

	api := application{
		config: cfg,
	}

	// h := api.mount()
	// api.run(h)

	//stuctured logging
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil)) // create a new logger instance
	slog.SetDefault(logger)                                 // set it as the default logger

	if err := api.run(api.mount()); err != nil {
		slog.Error("Server Failed to start", "Error", err) // sets level to Error while logging
		os.Exit(1)
	}
}
