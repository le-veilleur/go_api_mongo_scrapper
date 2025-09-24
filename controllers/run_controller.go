package controllers

import (
	"os"
	"os/exec"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/maxime-louis14/api-golang/logger"
)

// LaunchScraper lance le scraper via une route API
func LaunchScraper(c *fiber.Ctx) error {
	start := time.Now()
	requestID := c.Locals("requestID").(string)

	logger.LogInfo("Démarrage du scraper", map[string]interface{}{
		"request_id": requestID,
	})

	// Ajoute un délai de 4 secondes
	time.Sleep(4 * time.Second)

	// Exécute le scraper
	if err := RunScraper(); err != nil {
		logger.LogError("Erreur lors de l'exécution du scraper", err, map[string]interface{}{
			"request_id": requestID,
		})
		return c.Status(500).SendString("Erreur lors de l'exécution du scraper")
	}

	duration := time.Since(start)
	logger.LogInfo("Scraper exécuté avec succès", map[string]interface{}{
		"request_id": requestID,
		"duration":   duration.String(),
	})

	return c.Status(200).SendString("Scraper exécuté avec succès")
}

// RunScraper exécute le binaire du scraper
func RunScraper() error {
	start := time.Now()
	// Chemin vers le binaire du scraper
	scraperPath := "/go_api_mongo_scrapper/scraper/scraper"

	logger.LogInfo("Vérification de l'existence du binaire scraper", map[string]interface{}{
		"scraper_path": scraperPath,
	})

	// Vérifie que le fichier existe
	if _, err := os.Stat(scraperPath); os.IsNotExist(err) {
		logger.LogError("Binaire scraper introuvable", err, map[string]interface{}{
			"scraper_path": scraperPath,
		})
		return err
	}

	logger.LogInfo("Lancement du binaire scraper", map[string]interface{}{
		"scraper_path": scraperPath,
	})

	// Commande pour exécuter le scraper
	cmd := exec.Command(scraperPath)

	// Associe les sorties standard et erreur du scraper aux sorties du serveur
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Exécute la commande
	if err := cmd.Run(); err != nil {
		logger.LogError("Échec de l'exécution du scraper", err, map[string]interface{}{
			"scraper_path": scraperPath,
		})
		return err
	}

	duration := time.Since(start)
	logger.LogInfo("Scraper exécuté avec succès", map[string]interface{}{
		"scraper_path": scraperPath,
		"duration":     duration.String(),
	})
	return nil
}
