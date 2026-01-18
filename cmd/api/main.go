package main

import (
	"AwesomeProject/internal/db"
	"AwesomeProject/internal/env"
	"AwesomeProject/internal/store"
	"time"

	"go.uber.org/zap"
)

const version = "0.0.2"

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
		mail: mailConfig{
			exp: time.Hour * 24 * 3,
		},
		apiURL: env.GetString("EXTERNAL_URL", "localhost:8081"),
	}
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	_db, err := db.New(cfg.db.address, cfg.db.maxOpenConnections, cfg.db.maxIdleConnections, cfg.db.maxIdleTime)
	if err != nil {
		logger.Fatal(err)
		return
	}
	defer _db.Close()
	logger.Info("Successfully connected to database")
	_store := store.NewStorage(_db)
	app := &application{
		config: cfg,
		store:  &_store,
		logger: logger,
	}
	mux := app.mount()
	logger.Fatal(app.run(mux))
}
