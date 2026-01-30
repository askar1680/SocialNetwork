package main

import (
	"AwesomeProject/internal/auth"
	"AwesomeProject/internal/db"
	"AwesomeProject/internal/env"
	"AwesomeProject/internal/mailer"
	"AwesomeProject/internal/rateLimiter"
	"AwesomeProject/internal/store"
	"AwesomeProject/internal/store/cache"
	"time"

	"github.com/redis/go-redis/v9"

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
		env:     env.GetString("ENV", "development"),
		mail: mailConfig{
			fromEmail: env.GetString("FROM_EMAIL", ""),
			sendGrid: sendGridConfig{
				apiKey: env.GetString("SENDGRID_API_KEY", ""),
			},
			mailTrap: mailTrapConfig{
				apiKey: env.GetString("MAIL_TRAP_API_KEY", ""),
			},
			exp: time.Hour * 24 * 3,
		},
		apiURL:      env.GetString("EXTERNAL_URL", "localhost:8081"),
		frontendURL: env.GetString("FRONTEND_URL", "http://localhost:4000"),
		auth: authConfig{
			basic: basicConfig{
				username: env.GetString("AUTH_USERNAME", ""),
				password: env.GetString("AUTH_PASSWORD", ""),
			},
			token: tokenConfig{
				secret: env.GetString("AUTH_SECRET", ""),
				exp:    time.Hour * 24 * 3,
				iss:    "social network",
			},
		},
		redis: redisConfig{
			addr:    env.GetString("REDIS_ADDR", "localhost:6379"),
			pw:      env.GetString("REDIS_PASSWORD", ""),
			db:      0,
			enabled: env.GetBool("REDIS_ENABLED", true), // for now TRUE
		},
		rateLimiter: rateLimiter.Config{
			// TODO: change it
			RequestsPerTimeFrame: 20,
			TimeFrame:            time.Second * 5,
			Enabled:              true,
		},
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

	// Cache
	var rdb *redis.Client
	if cfg.redis.enabled {
		rdb = cache.NewRedisClient("127.0.0.1:6380", cfg.redis.pw, cfg.redis.db)
		logger.Info("Successfully connected to redis cache")
	}
	cacheStorage := cache.NewRedisStorage(rdb)

	_store := store.NewStorage(_db)

	// mailer := mailer.NewSendGridMailer(cfg.mail.sendGrid.apiKey, cfg.mail.fromEmail)
	mailer := mailer.NewMailtrapMailer(cfg.mail.mailTrap.apiKey, cfg.mail.fromEmail)
	if err != nil {
		logger.Fatal(err)
	}

	jwtAuthenticator := auth.NewJWTAuthenticator(cfg.auth.token.secret, cfg.auth.token.iss, cfg.auth.token.iss)
	_rateLimiter := rateLimiter.NewFixedWindowRateLimiter(cfg.rateLimiter.RequestsPerTimeFrame, cfg.rateLimiter.TimeFrame)
	app := &application{
		config:       cfg,
		store:        &_store,
		cacheStorage: &cacheStorage,
		logger:       logger,
		mailer:       mailer,
		auth:         jwtAuthenticator,
		rateLimiter:  _rateLimiter,
	}
	mux := app.mount()
	logger.Fatal(app.run(mux))
}
