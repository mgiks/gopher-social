package main

import (
	"log"

	"github.com/mgiks/gopher-social/internal/db"
	"github.com/mgiks/gopher-social/internal/env"
	"github.com/mgiks/gopher-social/internal/store"
)

func main() {
	cfg := config{
		port: env.GetString("PORT", ":8080"),
		db: dbConfig{
			url:          env.GetString("DB_URL", "postgres://admin:adminpassword@localhost:5434/gopher-social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
	}

	db, err := db.New(
		cfg.db.url,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)
	if err != nil {
		log.Fatal(err)
	}

	store := store.NewStore(db)

	app := application{
		config: cfg,
		store:  store,
	}

	mux := app.mount()

	log.Fatal(app.run(mux))
}
