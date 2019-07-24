package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/AleF83/airbag/config"

	log "github.com/sirupsen/logrus"
	"github.com/AleF83/airbag/middleware"
	"github.com/justinas/alice"
	"github.com/rs/cors"
	"github.com/vulcand/oxy/forward"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.WarnLevel)


	cfg, err := config.Init()
	if err != nil {
		log.Fatalf("Failed to init configuration: %v", err)
	}
	log.Println("Configuration initialized successfully.")

	myApp := http.HandlerFunc(newProxy(cfg.BackendURL))
	middlewares := alice.New(
		cors.Default().Handler,
		middleware.NewAuthMiddleware(cfg.JWTProviders, cfg.UnauthenticatedRoutes),
	).Then(myApp)

	http.Handle("/metrics", promhttp.Handler())
	http.Handle("/", middlewares)

	addr := fmt.Sprintf(":%v", cfg.Port)
	http.ListenAndServe(addr, nil)
}

func newProxy(backendURL *url.URL) http.HandlerFunc {
	fwd, err := forward.New()
	if err != nil {
		log.Fatalf("Failed to init proxy: %v", err)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		r.URL = backendURL
		fwd.ServeHTTP(w, r)
	}
}
