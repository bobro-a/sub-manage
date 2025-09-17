// @title           Sub-Manage API
// @version         1.0
// @description     This is a simple subscription management API.
// @host            localhost:9090
// @BasePath        /
package main

import (
	"context"
	"log"
	"sub-manage/internal/app"
	"sub-manage/pkg/config"
)

func main() {
	log.Println("Start sub-manage service")

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Failed ti load config: %v", err)
	}
	log.Printf("Loaded config: DB=%s, Migrations=%s",
		cfg.Database.URL,
		cfg.Migrations.Path)

	ctx := context.Background()

	a, err := app.New(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to init app: %v", err)
	}
	log.Println("App initialized success! Starting service...")

	if err = a.Start(); err != nil {
		log.Fatalf("Failed to start app: %v", err)
	}
}
