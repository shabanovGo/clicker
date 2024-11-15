package main

import (
    "log"

    "clicker/internal/app"
    "clicker/internal/config"
)

func main() {
    cfg, err := config.New()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    app := app.New(cfg)
    if err := app.Run(); err != nil {
        log.Fatalf("Application error: %v", err)
    }
}
