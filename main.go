package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/maxime-louis14/api-golang/database"
	"github.com/maxime-louis14/api-golang/routes"
)

func main() {
	// Initialisation des logs
	log.Println("Starting server...")

	// Initialisation de l'application Fiber
	app := fiber.New()
	log.Println("Fiber app initialized")

	// Connexion à MongoDB
	client := database.DBinstance()
	defer func() {
		log.Println("Closing MongoDB connection...")
		if err := client.Disconnect(nil); err != nil {
			log.Fatalf("Error disconnecting MongoDB client: %v", err)
		}
		log.Println("MongoDB connection closed")
	}()
	log.Println("Connected to MongoDB")

	// Configuration des routes
	routes.RecetteRoute(app)
	log.Println("Routes configured")

	// Démarrage du serveur
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on port %s", port)

	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
