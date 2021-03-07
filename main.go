package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"gopkg.in/square/go-jose.v2/json"
)

func main() {
	addr := "127.0.0.1:5000"

	logger, _ := zap.NewProduction()
	defer logger.Sync()
	logger.Info("boot",
		zap.String("version", "0.0.1"),
		zap.String("addr", addr))

	r := mux.NewRouter()

	// log middleware
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Info("request",
				zap.String("uri", r.RequestURI))

			next.ServeHTTP(w, r)
		})
	})

	r.HandleFunc("/", RootHandler)
	r.HandleFunc("/status", StatusHandler)

	srv := &http.Server{
		Handler:      r,
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logger.Error("exit", zap.Error(err))
		}
	}()

	// graceful shutdown with deadline
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello world")
}

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"ok": true,
	})
}
