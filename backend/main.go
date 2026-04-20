package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/Seeyko/casamia-api/handlers"
	custommiddleware "github.com/Seeyko/casamia-api/middleware"
	"github.com/Seeyko/casamia-api/services"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = "./uploads"
	}

	frontendURL := os.Getenv("FRONTEND_URL")
	allowedOrigins := []string{"http://localhost:8000", "http://localhost:4200", "http://localhost:3000"}
	if frontendURL != "" {
		allowedOrigins = append(allowedOrigins, frontendURL)
	}

	adminSecretPath := os.Getenv("ADMIN_SECRET_PATH")
	if adminSecretPath == "" {
		adminSecretPath = "_admin"
	}

	var adminAllowedIPs []string
	if ips := os.Getenv("ADMIN_ALLOWED_IPS"); ips != "" {
		for _, ip := range strings.Split(ips, ",") {
			ip = strings.TrimSpace(ip)
			if ip != "" {
				adminAllowedIPs = append(adminAllowedIPs, ip)
			}
		}
	}

	// Initialize database
	db, err := services.NewDatabase()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Seed data on first run
	services.SeedIfEmpty(db, uploadDir)

	// Initialize JWT service
	jwtService := services.NewJWTService()

	// Initialize handlers
	imageHandler := handlers.NewImageHandler(uploadDir)
	locationHandler := handlers.NewLocationHandler(db)
	newsHandler := handlers.NewNewsHandler(db)
	menuHandler := handlers.NewMenuHandler(db)
	adminHandler := handlers.NewAdminHandler(db, imageHandler)
	authHandler := handlers.NewAuthHandler(db, jwtService)

	// Initialize admin auth middleware (JWT-based)
	adminAuth := custommiddleware.NewAdminAuthMiddleware(jwtService, adminAllowedIPs)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Authorization"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Dynamic config.js for frontend
	r.Get("/js/config.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		w.Header().Set("Cache-Control", "no-cache")
		fmt.Fprintf(w, `window.APP_CONFIG = { API_URL: "", ADMIN_SECRET_PATH: "%s" };`, adminSecretPath)
	})

	r.Route("/api", func(r chi.Router) {
		r.Get("/health", handlers.HealthCheck)

		// Public endpoints
		r.Get("/locations", locationHandler.List)
		r.Get("/status", locationHandler.Status)
		r.Get("/news", newsHandler.ListPublic)
		r.Get("/menu", menuHandler.GetMenu)
		r.Get("/images/{filename}", imageHandler.ServeImage)

		// Auth endpoints (public — login + reset)
		r.Post("/auth/login", authHandler.Login)
		r.Post("/auth/reset-request", authHandler.RequestReset)
		r.Post("/auth/reset-password", authHandler.ResetPassword)

		// Admin endpoints (hidden behind secret path + JWT auth)
		r.Route("/"+adminSecretPath, func(r chi.Router) {
			r.Use(adminAuth.Authenticate)

			// Auth (authenticated)
			r.Post("/auth/change-password", authHandler.ChangePassword)

			// News
			r.Get("/news", adminHandler.ListNews)
			r.Post("/news", adminHandler.CreateNews)
			r.Put("/news/{id}", adminHandler.UpdateNews)
			r.Delete("/news/{id}", adminHandler.DeleteNews)

			// Menu categories
			r.Get("/menu/categories", adminHandler.ListCategories)
			r.Post("/menu/categories", adminHandler.CreateCategory)
			r.Put("/menu/categories/{id}", adminHandler.UpdateCategory)
			r.Delete("/menu/categories/{id}", adminHandler.DeleteCategory)

			// Menu items
			r.Get("/menu/items", adminHandler.ListItems)
			r.Post("/menu/items", adminHandler.CreateItem)
			r.Post("/menu/items/reorder", adminHandler.ReorderItems)
			r.Put("/menu/items/{id}", adminHandler.UpdateItem)
			r.Delete("/menu/items/{id}", adminHandler.DeleteItem)

			// Locations
			r.Put("/locations/{id}", adminHandler.UpdateLocation)
		})
	})

	// Serve frontend static files if directory exists
	staticDir := os.Getenv("STATIC_DIR")
	if staticDir == "" {
		staticDir = "./static"
	}

	if info, err := os.Stat(staticDir); err == nil && info.IsDir() {
		log.Printf("Serving frontend from: %s", staticDir)
		fs := http.FileServer(http.Dir(staticDir))

		// SPA routes — serve specific HTML files
		spaRoutes := map[string]string{
			"/menu":     "menu.html",
			"/histoire": "histoire.html",
		}
		// Admin route
		adminRoute := "/" + adminSecretPath
		spaRoutes[adminRoute] = "admin.html"

		for route, file := range spaRoutes {
			filePath := staticDir + "/" + file
			r.Get(route, func(w http.ResponseWriter, r *http.Request) {
				http.ServeFile(w, r, filePath)
			})
		}

		// Catch-all: try static file, fallback to index.html
		r.NotFound(func(w http.ResponseWriter, req *http.Request) {
			// Try serving the static file directly
			path := staticDir + req.URL.Path
			if fInfo, fErr := os.Stat(path); fErr == nil && !fInfo.IsDir() {
				fs.ServeHTTP(w, req)
				return
			}
			// Fallback to index.html for SPA
			http.ServeFile(w, req, staticDir+"/index.html")
		})
	} else {
		r.NotFound(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, `{"error":"not_found","path":"%s"}`, r.URL.Path)
		})
	}

	log.Printf("===========================================")
	log.Printf("CasaMia API starting on port %s", port)
	log.Printf("Upload directory: %s", uploadDir)
	log.Printf("CORS origins: %v", allowedOrigins)
	log.Printf("Admin API: /api/%s/*", adminSecretPath)
	log.Printf("Auth: POST /api/auth/login")
	log.Printf("===========================================")

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
