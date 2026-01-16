package main

import (
	"AwesomeProject/internal/db"
	"AwesomeProject/internal/env"
	"AwesomeProject/internal/store"
	"log"
)

const version = "0.0.1"

//	@title			GopherSocial API
//	@version		1.0
//	@description	API for Social network for Gophers
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8080
//	@BasePath	/v1

// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description
func main() {
	dbConfig := dbConfig{
		address:            env.GetString("DB_ADDRESS", "postgres://admin:adminpassword@localhost/socialnetwork?sslmode=disable"),
		maxOpenConnections: env.GetInt("DB_MAX_OPEN_CONNECTIONS", 30),
		maxIdleConnections: env.GetInt("DB_MAX_IDLE_CONNECTIONS", 30),
		maxIdleTime:        env.GetString("DB_MAX_IDLE_TIME", "15m"),
	}
	cfg := config{
		address: env.GetString("ADDRESS", ":8080"),
		db:      dbConfig,
		env:     "local",
	}
	_db, err := db.New(cfg.db.address, cfg.db.maxOpenConnections, cfg.db.maxIdleConnections, cfg.db.maxIdleTime)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer _db.Close()
	log.Println("Successfully connected to database")
	_store := store.NewStorage(_db)
	app := &application{cfg, &_store}
	mux := app.mount()
	log.Fatal(app.run(mux))
}
