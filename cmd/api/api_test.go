package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRateLimiterMiddleware(t *testing.T) {
	cfg := config{
		ratelimiter: ratelimiterConfig{
			requestsPerTimeFrame: 20,
			timeFrame:            time.Second * 5,
			enabled:              true,
		},
		port: ":8080",
	}

	app := newTestApplication(t, cfg, false)
	ts := httptest.NewServer(app.mount(false))
	defer ts.Close()

	client := &http.Client{}
	mockIP := "192.168.1.1"
	marginOfError := 2

	for i := range cfg.ratelimiter.requestsPerTimeFrame + marginOfError {
		req, err := http.NewRequest("GET", ts.URL+"/v1/health", nil)
		if err != nil {
			t.Fatal("could not create request", err)
		}

		req.Header.Set("X-Forwarded-For", mockIP)

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("could not send request: %v", err)
		}

		if i < cfg.ratelimiter.requestsPerTimeFrame {
			if resp.StatusCode != http.StatusOK {
				t.Errorf("expected status OK, got %v", resp.Status)
			}
		} else {
			if resp.StatusCode != http.StatusTooManyRequests {
				t.Errorf("expected status TOO MANY REQUESTS, got %v", resp.Status)
			}
		}
	}
}
