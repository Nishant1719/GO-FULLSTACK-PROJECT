package main

import (
	"log"
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

	if err := api.run(api.mount()); err != nil {
		log.Printf("Server Failed to start %s", err)
		os.Exit(1)
	}
}
