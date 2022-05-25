package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"time"
)

type key int

const (
	requestIDKey key = 0
	listenAddr string = ":8000"
)

var (
	healthy    int32
)

var store = make(map[string]Event)

func main() {
	logger := log.New(os.Stdout, "http: ", log.LstdFlags)
	logger.Println("Server is starting...")

	router := http.NewServeMux()
	router.Handle("/log-event", logEvent())

	nextRequestID := func() string {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}

	server := &http.Server{
		Addr:         listenAddr,
		Handler:      tracing(nextRequestID)(logging(logger)(router)),
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		logger.Println("Server is shutting down...")
		atomic.StoreInt32(&healthy, 0)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			logger.Fatalf("Could not gracefully shutdown the server: %v\n", err)
		}
		close(done)
	}()

	logger.Println("Server is ready to handle requests at", listenAddr)
	atomic.StoreInt32(&healthy, 1)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("Could not listen on %s: %v\n", listenAddr, err)
	}

	<-done
	logger.Println("Server stopped")
}

func logEvent() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		setupHeaders(w)

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method != "POST" {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		bodyBytes, error := io.ReadAll(r.Body)
		if error != nil {
			http.Error(w, error.Error(), http.StatusBadRequest)
			return
		}

		jsonBody := make(map[string]interface{})
		error = json.Unmarshal(bodyBytes, &jsonBody)
		if error != nil {
			http.Error(w, error.Error(), http.StatusBadRequest)
			return
		}

		sessionId, eventType, error := ValidateRequestBody(jsonBody)

		if error != nil {
			http.Error(w, error.Error(), http.StatusBadRequest)
			return
		}

		event := store[sessionId]
		if IsComplete(event) {
			fmt.Println("Event already completed")
			fmt.Println(event)
			w.WriteHeader(http.StatusOK)
			return
		}
		event, error = ParseEvent(eventType, event, bodyBytes)
		if error != nil {
			http.Error(w, error.Error(), http.StatusBadRequest)
			return
		}
		store[sessionId] = event
		if IsComplete(event) {
			fmt.Println("Event complete!")
		}
		fmt.Println(event)
		w.WriteHeader(http.StatusOK)
		return
	})
}
