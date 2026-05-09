package main

import (
	"time"

	"github.com/mgiks/gopher-social/internal/db"
	"github.com/mgiks/gopher-social/internal/env"
	"github.com/mgiks/gopher-social/internal/store"
	"github.com/mgiks/gopher-social/internal/validator"
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
		port:   env.GetString("PORT", ":8080"),
		apiURL: env.GetString("EXTERNAL_URL", "localhost:8080"),
		db: dbConfig{
			url:          env.GetString("DB_URL", "postgres://admin:adminpassword@localhost:5434/gopher-social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env: env.GetString("ENV", "development"),
		mail: mailConfig{
			exp: time.Hour * 24 * 3, // 3 days
		},
	}

	loggerConfig := zap.NewDevelopmentConfig()
	loggerConfig.DisableStacktrace = true
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

	store := store.NewStore(db)

	validator := validator.New()

	app := application{
		config:    cfg,
		store:     store,
		logger:    logger,
		validator: validator,
	}

	mux := app.mount()

	app.logger.Fatal(app.run(mux))
}
