package main

import (
	"embed"
	"log"
	"net/http"

	"rdips-dashboard/internal/config"
	"rdips-dashboard/internal/controllers"
	"rdips-dashboard/internal/services"
	"rdips-dashboard/internal/views"
)

//go:embed templates/*.html static/*
var assets embed.FS

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	renderer, err := views.New(assets)
	if err != nil {
		log.Fatalf("load templates: %v", err)
	}

	api := services.NewAPIService(cfg)
	dashboard := controllers.NewDashboardController(api, renderer)

	mux := http.NewServeMux()
	mux.Handle("GET /static/", http.FileServerFS(assets))
	dashboard.RegisterRoutes(mux)

	log.Printf("rdips-dashboard listening on :%s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, mux))
}
