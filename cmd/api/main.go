package main

import (
	"time"

	"github.com/mgiks/gopher-social/internal/auth"
	"github.com/mgiks/gopher-social/internal/db"
	"github.com/mgiks/gopher-social/internal/env"
	"github.com/mgiks/gopher-social/internal/mailer"
	"github.com/mgiks/gopher-social/internal/store/cache"
	store "github.com/mgiks/gopher-social/internal/store/db"
	"github.com/mgiks/gopher-social/internal/validator"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const apiVersion = "0.0.2"

//	@title			GopherSocial API
//	@description	API for GopherSocial, a social network for gophers
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath					/v1
//
// @securityDefinitions.apiKey	ApiKeyAuth
// @in							header
// @name						Authorization
func main() {
	cfg := config{
		port:        env.GetString("BACKEND_PORT", ":8080"),
		apiURL:      env.GetString("EXTERNAL_URL", "localhost:8080"),
		frontendURL: env.GetString("FRONTEND_URL", "http://localhost:5173"),
		db: dbConfig{
			url:          env.GetString("DB_URL", "postgres://admin:adminpassword@localhost:5434/gopher-social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		redis: redisConfig{
			addr:     env.GetString("REDIS_ADDR", "localhost:6379"),
			password: env.GetString("REDIS_PASSWORD", ""),
			db:       env.GetInt("REDIS_DB", 0),
			enabled:  env.GetBool("REDIS_ENABLED", false),
		},
		env: env.GetString("ENV", "development"),
		mail: mailConfig{
			exp:       time.Hour * 24 * 3, // 3 days
			fromEmail: env.GetString("FROM_EMAIL", ""),
			sendGrid: sendGridConfig{
				apiKey: env.GetString("SENDGRID_API_KEY", ""),
			},
			resend: resendConfig{
				apiKey: env.GetString("RESEND_API_KEY", ""),
			},
			mailTrap: mailTrapConfig{
				apiKey: env.GetString("MAILTRAP_API_KEY", ""),
			},
		},
		auth: authConfig{
			basic: basicConfig{
				username: env.GetString("BASIC_AUTH_USERNAME", "admin"),
				password: env.GetString("BASIC_AUTH_PASSWORD", "admin"),
			},
			token: tokenConfig{
				secret: env.GetString("TOKEN_AUTH_SECRET", "example"),
				exp:    time.Hour * 24 * 3, // days
				iss:    "gophersocial",
			},
		}}

	var loggerConfig zap.Config

	switch cfg.env {
	case "development":
		loggerConfig = zap.NewDevelopmentConfig()
		loggerConfig.DisableStacktrace = true
	case "production":
		loggerConfig = zap.NewProductionConfig()
	}

	logger := zap.Must(loggerConfig.Build()).Sugar()
	defer logger.Sync()

	db, err := db.New(
		cfg.db.url,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)
	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()
	logger.Info("database connection pool established")

	var rdb *redis.Client
	if cfg.redis.enabled {
		rdb = cache.NewRedisClient(cfg.redis.addr, cfg.redis.password, cfg.redis.db)
		logger.Info("redis cache connection established")
	}

	store := store.NewStore(db)
	cacheStore := cache.NewCacheStore(rdb)

	validator := validator.New()

	mailer, err := mailer.NewMailTrap(cfg.mail.mailTrap.apiKey, cfg.mail.fromEmail)
	if err != nil {
		logger.Fatal(err)
	}

	jwtAuthenticator := auth.NewJWTAuthenticator(
		cfg.auth.token.secret,
		cfg.auth.token.iss,
		cfg.auth.token.iss,
	)

	app := application{
		config:        cfg,
		store:         store,
		cache:         cacheStore,
		logger:        logger,
		validator:     validator,
		mailer:        mailer,
		authenticator: jwtAuthenticator,
	}

	mux := app.mount()

	app.logger.Fatal(app.run(mux))
}
