package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mgiks/gopher-social/internal/auth"
	"github.com/mgiks/gopher-social/internal/store/cache"
	store "github.com/mgiks/gopher-social/internal/store/db"
	"go.uber.org/zap"
)

func newTestApplication(t *testing.T, useLogger bool) application {
	t.Helper()

	var logger *zap.SugaredLogger
	if useLogger {
		loggerConfig := zap.NewDevelopmentConfig()
		loggerConfig.DisableStacktrace = true
		logger = zap.Must(loggerConfig.Build()).Sugar()
	} else {
		logger = zap.NewNop().Sugar()
	}
	mockStore := store.NewMockStore()
	mockCacheStore := cache.NewMockStore()
	authenticator := auth.NewMockAuthenticator()

	return application{
		logger:        logger,
		store:         mockStore,
		cache:         mockCacheStore,
		authenticator: authenticator,
	}

}

func executeRequest(req *http.Request, mux http.Handler) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("expected response code %d, got %d", expected, actual)
	}
}
