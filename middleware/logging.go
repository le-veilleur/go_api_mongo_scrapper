package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/maxime-louis14/api-golang/logger"
)

// generateRequestID génère un ID unique pour chaque requête
func generateRequestID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// LoggingMiddleware middleware de logging détaillé
func LoggingMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		requestID := generateRequestID()

		// Ajouter l'ID de requête au contexte
		c.Locals("requestID", requestID)

		// Log de début de requête
		logger.LogRequest(
			logger.INFO,
			"Début de requête",
			requestID,
			c.Method(),
			c.Path(),
			c.Get("User-Agent"),
			c.IP(),
			0, // Status code sera mis à jour après
			time.Since(start),
		)

		// Exécuter la requête
		err := c.Next()

		// Calculer la latence totale
		latency := time.Since(start)

		// Log de fin de requête
		logger.LogRequest(
			logger.INFO,
			"Fin de requête",
			requestID,
			c.Method(),
			c.Path(),
			c.Get("User-Agent"),
			c.IP(),
			c.Response().StatusCode(),
			latency,
		)

		return err
	}
}

// DatabaseLoggingMiddleware middleware pour logger les opérations de base de données
func DatabaseLoggingMiddleware(operation string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		requestID := c.Locals("requestID").(string)

		// Log de début d'opération DB
		logger.LogDatabase(
			logger.INFO,
			"Début d'opération de base de données",
			operation,
			"mongodb",
			time.Since(start),
			map[string]interface{}{
				"request_id": requestID,
				"path":       c.Path(),
			},
		)

		// Exécuter la requête
		err := c.Next()

		// Calculer la durée
		duration := time.Since(start)

		// Log de fin d'opération DB
		level := logger.INFO
		message := "Opération de base de données terminée"
		if err != nil {
			level = logger.ERROR
			message = "Erreur lors de l'opération de base de données"
		}

		logger.LogDatabase(
			level,
			message,
			operation,
			"mongodb",
			duration,
			map[string]interface{}{
				"request_id": requestID,
				"path":       c.Path(),
				"error":      err,
			},
		)

		return err
	}
}
