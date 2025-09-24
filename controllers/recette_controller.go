package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/maxime-louis14/api-golang/database"
	"github.com/maxime-louis14/api-golang/logger"
	"github.com/maxime-louis14/api-golang/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var recetteCollection *mongo.Collection = database.OpenCollection(database.Client, "recettes")

// getScraperDataPath retourne un chemin absolu vers data.json
func getScraperDataPath() (string, error) {
	// Essayer d'abord le chemin local en développement
	localPath := "/home/maka/GitHub/go_api_mongo_scrapper/scraper/data.json"
	if _, err := os.Stat(localPath); err == nil {
		return localPath, nil
	}

	// Essayer le chemin du volume monté
	volumePath := "/go_api_mongo_scrapper/scraper/data.json"
	if _, err := os.Stat(volumePath); err == nil {
		return volumePath, nil
	}

	// Chemin absolu pour Docker
	dataPath := "/go_api_mongo_scrapper/scraper/data.json"
	if _, err := os.Stat(dataPath); err == nil {
		return dataPath, nil
	}

	return "", errors.New("data.json file does not exist at " + localPath + ", " + volumePath + ", or " + dataPath)
}

// PostRecette ajoute des recettes en batch depuis un fichier JSON
func PostRecette(c *fiber.Ctx) error {
	start := time.Now()
	requestID := c.Locals("requestID").(string)

	logger.LogInfo("Début de l'importation des recettes", map[string]interface{}{
		"request_id": requestID,
	})

	// Obtenir le chemin complet vers data.json
	dataPath, err := getScraperDataPath()
	if err != nil {
		logger.LogError("Échec de localisation du fichier data.json", err, map[string]interface{}{
			"request_id": requestID,
		})
		return c.Status(500).SendString("Erreur lors de la localisation du fichier data.json")
	}

	// Debug: afficher le chemin trouvé
	logger.LogInfo("Chemin du fichier data.json trouvé", map[string]interface{}{
		"request_id": requestID,
		"file_path":  dataPath,
	})

	// Ouvrir le fichier data.json
	file, err := os.Open(dataPath)
	if err != nil {
		logger.LogError("Échec d'ouverture du fichier data.json", err, map[string]interface{}{
			"request_id": requestID,
			"file_path":  dataPath,
		})
		return c.Status(500).SendString("Erreur lors de l'ouverture du fichier data.json")
	}
	defer file.Close()

	// Lire les données JSON
	data, err := ioutil.ReadAll(file)
	if err != nil {
		logger.LogError("Échec de lecture du fichier data.json", err, map[string]interface{}{
			"request_id": requestID,
			"file_path":  dataPath,
		})
		return c.Status(500).SendString("Erreur lors de la lecture du fichier data.json")
	}

	// Décoder les données JSON
	var recettes []models.Recette
	if err := json.Unmarshal(data, &recettes); err != nil {
		logger.LogError("Échec du décodage JSON", err, map[string]interface{}{
			"request_id": requestID,
		})
		return c.Status(500).SendString("Erreur lors du décodage des données JSON")
	}

	// Insérer les recettes dans MongoDB
	insertedCount := 0
	for _, recette := range recettes {
		_, err := recetteCollection.InsertOne(context.Background(), recette)
		if err != nil {
			logger.LogError("Échec d'insertion d'une recette", err, map[string]interface{}{
				"request_id": requestID,
				"recette":    recette.Name,
			})
			return c.Status(500).SendString("Erreur lors de l'insertion des recettes")
		}
		insertedCount++
	}

	duration := time.Since(start)
	logger.LogDatabase(logger.INFO, "Importation des recettes terminée", "batch_insert", "mongodb", duration, map[string]interface{}{
		"request_id":     requestID,
		"recettes_count": insertedCount,
	})

	return c.Status(201).SendString("Recettes ajoutées avec succès")
}

// GetAllRecettes retourne toutes les recettes
func GetAllRecettes(c *fiber.Ctx) error {
	start := time.Now()
	requestID := c.Locals("requestID").(string)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logger.LogDatabase(logger.INFO, "Début de récupération de toutes les recettes", "find_all", "mongodb", time.Since(start), map[string]interface{}{
		"request_id": requestID,
	})

	// Récupérer toutes les recettes
	cursor, err := recetteCollection.Find(ctx, bson.M{})
	if err != nil {
		logger.LogError("Échec de récupération des recettes", err, map[string]interface{}{
			"request_id": requestID,
		})
		return c.Status(500).SendString("Erreur lors de la récupération des recettes")
	}
	defer cursor.Close(ctx)

	// Décoder les recettes
	var recettes []models.Recette
	if err := cursor.All(ctx, &recettes); err != nil {
		logger.LogError("Échec du décodage des recettes", err, map[string]interface{}{
			"request_id": requestID,
		})
		return c.Status(500).SendString("Erreur lors du décodage des recettes")
	}

	duration := time.Since(start)
	logger.LogDatabase(logger.INFO, "Récupération de toutes les recettes terminée", "find_all", "mongodb", duration, map[string]interface{}{
		"request_id":     requestID,
		"recettes_count": len(recettes),
	})

	return c.Status(200).JSON(recettes)
}

// GetRecetteByID retourne une recette spécifique en fonction de son ID
func GetRecetteByID(c *fiber.Ctx) error {
	start := time.Now()
	requestID := c.Locals("requestID").(string)
	id := c.Params("id")

	logger.LogInfo("Recherche de recette par ID", map[string]interface{}{
		"request_id": requestID,
		"recipe_id":  id,
	})

	// Convertir l'ID en ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		logger.LogError("ID de recette invalide", err, map[string]interface{}{
			"request_id": requestID,
			"recipe_id":  id,
		})
		return c.Status(400).SendString("ID de recette invalide")
	}

	// Rechercher la recette
	filter := bson.M{"_id": objID}
	var recette models.Recette
	if err := recetteCollection.FindOne(context.Background(), filter).Decode(&recette); err != nil {
		logger.LogError("Recette introuvable", err, map[string]interface{}{
			"request_id": requestID,
			"recipe_id":  id,
		})
		return c.Status(404).SendString("Recette introuvable")
	}

	duration := time.Since(start)
	logger.LogDatabase(logger.INFO, "Recette trouvée par ID", "find_one", "mongodb", duration, map[string]interface{}{
		"request_id":  requestID,
		"recipe_id":   id,
		"recipe_name": recette.Name,
	})

	return c.Status(200).JSON(recette)
}

// GetRecetteByName retourne une recette en fonction de son nom
func GetRecetteByName(c *fiber.Ctx) error {
	start := time.Now()
	requestID := c.Locals("requestID").(string)
	nomRecette := strings.ReplaceAll(c.Params("name"), "%20", " ")

	logger.LogInfo("Recherche de recette par nom", map[string]interface{}{
		"request_id":  requestID,
		"recipe_name": nomRecette,
	})

	// Rechercher la recette par nom
	filter := bson.M{"name": nomRecette}
	var recette models.Recette
	if err := recetteCollection.FindOne(context.Background(), filter).Decode(&recette); err != nil {
		logger.LogError("Recette introuvable par nom", err, map[string]interface{}{
			"request_id":  requestID,
			"recipe_name": nomRecette,
		})
		return c.Status(404).SendString("Recette introuvable")
	}

	duration := time.Since(start)
	logger.LogDatabase(logger.INFO, "Recette trouvée par nom", "find_one", "mongodb", duration, map[string]interface{}{
		"request_id":  requestID,
		"recipe_name": nomRecette,
	})

	return c.Status(200).JSON(recette)
}

// GetRecettesByIngredient retourne toutes les recettes contenant un ingrédient spécifique
func GetRecettesByIngredient(c *fiber.Ctx) error {
	start := time.Now()
	requestID := c.Locals("requestID").(string)
	ingredient := c.Params("unit")

	logger.LogInfo("Recherche de recettes par ingrédient", map[string]interface{}{
		"request_id": requestID,
		"ingredient": ingredient,
	})

	// Rechercher les recettes par ingrédient
	filter := bson.M{"ingredients": bson.M{"$elemMatch": bson.M{"unit": ingredient}}}
	cursor, err := recetteCollection.Find(context.Background(), filter)
	if err != nil {
		logger.LogError("Échec de récupération des recettes par ingrédient", err, map[string]interface{}{
			"request_id": requestID,
			"ingredient": ingredient,
		})
		return c.Status(500).SendString("Erreur lors de la récupération des recettes")
	}
	defer cursor.Close(context.Background())

	// Décoder les recettes
	var recettes []models.Recette
	if err := cursor.All(context.Background(), &recettes); err != nil {
		logger.LogError("Échec du décodage des recettes par ingrédient", err, map[string]interface{}{
			"request_id": requestID,
			"ingredient": ingredient,
		})
		return c.Status(500).SendString("Erreur lors du décodage des recettes")
	}

	duration := time.Since(start)
	logger.LogDatabase(logger.INFO, "Recettes trouvées par ingrédient", "find_many", "mongodb", duration, map[string]interface{}{
		"request_id":     requestID,
		"ingredient":     ingredient,
		"recettes_count": len(recettes),
	})

	return c.Status(200).JSON(recettes)
}
