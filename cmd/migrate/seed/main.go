package main

import (
	"AwesomeProject/internal/db"
	"AwesomeProject/internal/env"
	"AwesomeProject/internal/store"
	"strings"
)

func main() {
	addr := env.GetString("DB_ADDRESS", "postgres://admin:adminpassword@localhost/socialnetwork?sslmode=disable")
	addr = strings.TrimSpace(addr)
	addr = strings.Trim(addr, `"'`)
	conn, err := db.New(addr, 30, 30, "15m")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	_store := store.NewStorage(conn)
	db.Seed(_store)
}
