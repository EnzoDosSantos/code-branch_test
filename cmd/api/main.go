package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/EnzoDosSantos/code-branch_test/internal/server"
)

func main() {
    mux := server.SetupRouter()
    
    srv := &http.Server{
        Addr:    ":8080",
        Handler: mux,
    }

    done := make(chan os.Signal, 1)
    signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

    go func() {
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("server error: %v", err)
        }
    }()

    log.Println("Server started on :8080")

    <-done
    log.Println("Server is shutting down...")

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
        log.Fatalf("Server shutdown failed: %v", err)
    }
    log.Println("Server exited properly")
}