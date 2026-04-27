package main

import (
	"log"

	"github.com/mgiks/gopher-social/internal/env"
	"github.com/mgiks/gopher-social/internal/store"
)

func main() {
	cfg := config{
		addr: env.GetString("PORT", ":8080"),
	}

	store := store.NewStore(nil)

	app := application{
		config: cfg,
		store:  store,
	}

	mux := app.mount()

	log.Fatal(app.run(mux))
}
