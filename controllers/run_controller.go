package controllers

import (
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/gofiber/fiber/v2"
)

// LaunchScraper lance le scraper via une route API
func LaunchScraper(c *fiber.Ctx) error {
	// Ajoute un délai de 4 secondes
	time.Sleep(4 * time.Second)

	// Exécute le scraper
	if err := RunScraper(); err != nil {
		log.Printf("Error while running scraper: %v", err)
		return c.Status(500).SendString("Erreur lors de l'exécution du scraper")
	}

	return c.Status(200).SendString("Scraper exécuté avec succès")
}

// RunScraper exécute le binaire du scraper
func RunScraper() error {
	// Chemin vers le binaire du scraper
	scraperPath := "/go_api_mongo_scrapper/scraper/scraper"

	// Vérifie que le fichier existe
	if _, err := os.Stat(scraperPath); os.IsNotExist(err) {
		return err
	}

	// Commande pour exécuter le scraper
	cmd := exec.Command(scraperPath)

	// Associe les sorties standard et erreur du scraper aux sorties du serveur
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Exécute la commande
	if err := cmd.Run(); err != nil {
		log.Printf("Failed to run scraper: %v", err)
		return err
	}

	log.Println("Scraper executed successfully")
	return nil
}
