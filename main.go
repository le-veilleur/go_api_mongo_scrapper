package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/maxime-louis14/api-golang/database"
	"github.com/maxime-louis14/api-golang/routes"
)

// Variables de versioning injectées lors du build
var (
	version   = "dev"
	gitCommit = "unknown"
	buildTime = "unknown"
)

// BuildInfo contient les informations de build
type BuildInfo struct {
	Version   string `json:"version"`
	GitCommit string `json:"git_commit"`
	BuildTime string `json:"build_time"`
	GoVersion string `json:"go_version"`
	OS        string `json:"os"`
	Arch      string `json:"arch"`
}

// HealthResponse structure pour le health check
type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Build     BuildInfo `json:"build"`
	Database  string    `json:"database"`
}

func main() {
	// Affichage des informations de version
	fmt.Printf("Go API MongoDB Scrapper\n")
	fmt.Printf("Version: %s\n", version)
	fmt.Printf("Git Commit: %s\n", gitCommit)
	fmt.Printf("Build Time: %s\n", buildTime)
	fmt.Printf("Go Version: %s\n", runtime.Version())
	fmt.Printf("OS/Arch: %s/%s\n\n", runtime.GOOS, runtime.GOARCH)

	// Initialisation des logs
	log.Println("Starting server...")

	// Initialisation de l'application Fiber avec configuration
	app := fiber.New(fiber.Config{
		AppName:      fmt.Sprintf("Go API MongoDB Scrapper v%s", version),
		ServerHeader: "Go API MongoDB Scrapper",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error":   true,
				"message": err.Error(),
				"version": version,
			})
		},
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} - ${latency}\n",
	}))
	app.Use(cors.New())

	log.Println("Fiber app initialized with middleware")

	// Connexion à MongoDB
	client := database.DBinstance()
	defer func() {
		log.Println("Closing MongoDB connection...")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := client.Disconnect(ctx); err != nil {
			log.Fatalf("Error disconnecting MongoDB client: %v", err)
		}
		log.Println("MongoDB connection closed")
	}()
	log.Println("Connected to MongoDB")

	// Route de health check
	app.Get("/health", func(c *fiber.Ctx) error {
		// Test de la connexion MongoDB
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		dbStatus := "connected"
		if err := client.Ping(ctx, nil); err != nil {
			dbStatus = "disconnected"
		}

		return c.JSON(HealthResponse{
			Status:    "ok",
			Timestamp: time.Now(),
			Build: BuildInfo{
				Version:   version,
				GitCommit: gitCommit,
				BuildTime: buildTime,
				GoVersion: runtime.Version(),
				OS:        runtime.GOOS,
				Arch:      runtime.GOARCH,
			},
			Database: dbStatus,
		})
	})

	// Route d'informations de version
	app.Get("/version", func(c *fiber.Ctx) error {
		return c.JSON(BuildInfo{
			Version:   version,
			GitCommit: gitCommit,
			BuildTime: buildTime,
			GoVersion: runtime.Version(),
			OS:        runtime.GOOS,
			Arch:      runtime.GOARCH,
		})
	})

	// Configuration des routes API
	routes.RecetteRoute(app)
	log.Println("Routes configured")

	// Démarrage du serveur
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on port %s", port)
	log.Printf("Health check available at: http://localhost:%s/health", port)
	log.Printf("Version info available at: http://localhost:%s/version", port)

	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
