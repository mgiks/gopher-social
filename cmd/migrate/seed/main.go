package main

import (
	"log"

	"github.com/mgiks/gopher-social/internal/db"
	"github.com/mgiks/gopher-social/internal/env"
	store "github.com/mgiks/gopher-social/internal/store/db"
)

func main() {
	dbUrl := env.GetString("DB_URL", "postgres://admin:adminpassword@localhost:5434/gopher-social?sslmode=disable")
	conn, err := db.New(dbUrl, 3, 3, "15m")
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	store := store.NewStore(conn)
	db.Seed(store, conn)
}
