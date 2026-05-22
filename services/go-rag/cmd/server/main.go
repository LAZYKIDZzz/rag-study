package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/lazykidzzz/rag-study/services/go-rag/internal/rag"
)

func main() {
	addr := ":" + envOrDefault("PORT", "8080")

	service := rag.NewService()
	handler := rag.NewHTTPHandler(service)

	server := &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("go-rag listening on %s", addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server failed: %v", err)
	}
}

func envOrDefault(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
