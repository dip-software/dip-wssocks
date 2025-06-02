package health

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

func StartServer(port int) {
	// Define HTTP endpoints
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"UP"}`))
	})

	// Create a simple HTTP server with reasonable timeouts
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start the HTTP server in a goroutine
	go func() {
		slog.Info("Starting monitoring server",
			"port", port,
			"endpoints", []string{"/health"})

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Monitoring server error", "error", err)
		}
	}()
}
