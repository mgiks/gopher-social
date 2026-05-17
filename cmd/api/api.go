package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"go.uber.org/zap"

	"github.com/mgiks/gopher-social/docs" // This is required to generate swagger docs
	"github.com/mgiks/gopher-social/internal/auth"
	"github.com/mgiks/gopher-social/internal/mailer"
	"github.com/mgiks/gopher-social/internal/ratelimiter"
	"github.com/mgiks/gopher-social/internal/store/cache"
	store "github.com/mgiks/gopher-social/internal/store/db"
	"github.com/mgiks/gopher-social/internal/validator"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type application struct {
	config        config
	store         store.Store
	cache         cache.Store
	logger        *zap.SugaredLogger
	validator     validator.Validator
	mailer        mailer.SenderCreator
	authenticator auth.Authenticator
	rateLimiter   ratelimiter.Limiter
}

type config struct {
	port        string
	db          dbConfig
	env         string
	apiURL      string
	mail        mailConfig
	frontendURL string
	auth        authConfig
	redis       redisConfig
	ratelimiter ratelimiterConfig
}

type ratelimiterConfig struct {
	requestsPerTimeFrame int
	timeFrame            time.Duration
	enabled              bool
}

type redisConfig struct {
	addr     string
	password string
	db       int
	enabled  bool
}

type dbConfig struct {
	url          string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

type mailConfig struct {
	exp       time.Duration
	fromEmail string
	sendGrid  sendGridConfig
	resend    resendConfig
	mailTrap  mailTrapConfig
}

type sendGridConfig struct {
	apiKey string
}

type resendConfig struct {
	apiKey string
}

type mailTrapConfig struct {
	apiKey string
}

type authConfig struct {
	basic basicConfig
	token tokenConfig
}

type basicConfig struct {
	username string
	password string
}

type tokenConfig struct {
	secret string
	exp    time.Duration
	iss    string
}

func (app application) mount(useLogger bool) http.Handler {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)

	// Needed to turn off logging during tests
	if useLogger {
		r.Use(middleware.Logger)
	}

	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(app.RateLimiterMiddleware)

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)

		docsURL := fmt.Sprintf("%s/swagger/doc.json", app.config.port)
		r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(docsURL)))

		r.Route("/posts", func(r chi.Router) {
			r.Use(app.AuthTokenMiddleware)

			r.Post("/", app.createPostHandler)

			r.Route("/{postID}", func(r chi.Router) {
				r.Use(app.postContextMiddleware)

				r.Get("/", app.getPostHandler)
				r.Patch("/", app.checkPostOwnership("moderator", app.updatePostHandler))
				r.Delete("/", app.checkPostOwnership("admin", app.deletePostHandler))
			})
		})

		r.Route("/users", func(r chi.Router) {
			r.Put("/activate/{token}", app.activateUserHandler)

			r.Route("/{userID}", func(r chi.Router) {
				r.Use(app.AuthTokenMiddleware)

				r.Get("/", app.getUserHandler)
				r.Put("/follow", app.followUserHandler)
				r.Put("/unfollow", app.unfollowUserHandler)
			})

			r.Group(func(r chi.Router) {
				r.Use(app.AuthTokenMiddleware)
				r.Get("/feed", app.getUserFeedHandler)
			})
		})

		r.Route("/auth", func(r chi.Router) {
			r.Post("/user", app.registerUserHandler)
			r.Post("/token", app.createTokenHandler)
		})
	})

	return r
}

func (app application) run(mux http.Handler) error {
	docs.SwaggerInfo.Version = apiVersion
	docs.SwaggerInfo.Host = app.config.apiURL
	docs.SwaggerInfo.BasePath = "/v1"

	srv := &http.Server{
		Addr:         app.config.port,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	shutdown := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGTERM, os.Interrupt)
		s := <-quit

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		app.logger.Infow("signal caught", "signal", s.String())
		shutdown <- srv.Shutdown(ctx)
	}()

	app.logger.Infow("server has started", "port", app.config.port, "env", app.config.env)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	start := time.Now()
	err = <-shutdown
	if err != nil {
		return err
	}

	app.logger.Infow("server has stopped", "port", app.config.port, "env", app.config.env, "shutdown in", time.Since(start))
	return nil
}
