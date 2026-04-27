package main

import (
	"log"

	"github.com/mgiks/gopher-social/internal/env"
)

func main() {
	cfg := config{
		addr: env.GetString("PORT", ":8080"),
	}

	app := application{config: cfg}

	mux := app.mount()

	log.Fatal(app.run(mux))
}
