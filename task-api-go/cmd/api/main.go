package main

import (
    "log"
    "net/http"
    "time"

    "taskapi/internal/handlers"
    "taskapi/internal/middleware"
    "taskapi/internal/storage"
)

func main() {
    store := storage.NewMemoryStore()

    mux := http.NewServeMux()
    h := handlers.NewTaskHandler(store)

    // Routing
    mux.HandleFunc("/tasks", h.Tasks)

    // Middleware chain: logging -> auth -> mux
    var handler http.Handler = mux
    handler = middleware.Logging(handler, "task-api")
    handler = middleware.APIKey(handler, "secret12345")

    srv := &http.Server{
        Addr:              ":8080",
        Handler:           handler,
        ReadHeaderTimeout: 5 * time.Second,
    }

    log.Printf("listening on %s", srv.Addr)
    if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        log.Fatalf("server error: %v", err)
    }
}
