package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/maxime-louis14/api-golang/database"
	"github.com/maxime-louis14/api-golang/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var recetteCollection *mongo.Collection = database.OpenCollection(database.Client, "recettes")

// getScraperDataPath retourne un chemin absolu vers data.json
func getScraperDataPath() (string, error) {
	// Chemin absolu pour Docker
	dataPath := "/go_api_mongo_scrapper/scraper/data.json"
	if _, err := os.Stat(dataPath); os.IsNotExist(err) {
		return "", errors.New("data.json file does not exist at " + dataPath)
	}
	return dataPath, nil
}

// PostRecette ajoute des recettes en batch depuis un fichier JSON
func PostRecette(c *fiber.Ctx) error {
	// Obtenir le chemin complet vers data.json
	dataPath, err := getScraperDataPath()
	if err != nil {
		log.Printf("Failed to get data path: %v", err)
		return c.Status(500).SendString("Erreur lors de la localisation du fichier data.json")
	}

	// Ouvrir le fichier data.json
	file, err := os.Open(dataPath)
	if err != nil {
		log.Printf("Failed to open file %s: %v", dataPath, err)
		return c.Status(500).SendString("Erreur lors de l'ouverture du fichier data.json")
	}
	defer file.Close()

	// Lire les données JSON
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("Failed to read file %s: %v", dataPath, err)
		return c.Status(500).SendString("Erreur lors de la lecture du fichier data.json")
	}

	// Décoder les données JSON
	var recettes []models.Recette
	if err := json.Unmarshal(data, &recettes); err != nil {
		log.Printf("Failed to unmarshal JSON: %v", err)
		return c.Status(500).SendString("Erreur lors du décodage des données JSON")
	}

	// Insérer les recettes dans MongoDB
	for _, recette := range recettes {
		_, err := recetteCollection.InsertOne(context.Background(), recette)
		if err != nil {
			log.Printf("Failed to insert recette: %v", err)
			return c.Status(500).SendString("Erreur lors de l'insertion des recettes")
		}
	}
	log.Println("Recettes ajoutées avec succès")

	return c.Status(201).SendString("Recettes ajoutées avec succès")
}

// GetAllRecettes retourne toutes les recettes
func GetAllRecettes(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Récupérer toutes les recettes
	cursor, err := recetteCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Printf("Failed to fetch recettes: %v", err)
		return c.Status(500).SendString("Erreur lors de la récupération des recettes")
	}
	defer cursor.Close(ctx)

	// Décoder les recettes
	var recettes []models.Recette
	if err := cursor.All(ctx, &recettes); err != nil {
		log.Printf("Failed to decode recettes: %v", err)
		return c.Status(500).SendString("Erreur lors du décodage des recettes")
	}

	return c.Status(200).JSON(recettes)
}

// GetRecetteByID retourne une recette spécifique en fonction de son ID
func GetRecetteByID(c *fiber.Ctx) error {
	id := c.Params("id")

	// Convertir l'ID en ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("Invalid recipe ID: %v", err)
		return c.Status(400).SendString("ID de recette invalide")
	}

	// Rechercher la recette
	filter := bson.M{"_id": objID}
	var recette models.Recette
	if err := recetteCollection.FindOne(context.Background(), filter).Decode(&recette); err != nil {
		log.Printf("Recipe not found: %v", err)
		return c.Status(404).SendString("Recette introuvable")
	}

	return c.Status(200).JSON(recette)
}

// GetRecetteByName retourne une recette en fonction de son nom
func GetRecetteByName(c *fiber.Ctx) error {
	nomRecette := strings.ReplaceAll(c.Params("name"), "%20", " ")

	// Rechercher la recette par nom
	filter := bson.M{"name": nomRecette}
	var recette models.Recette
	if err := recetteCollection.FindOne(context.Background(), filter).Decode(&recette); err != nil {
		log.Printf("Recipe not found: %v", err)
		return c.Status(404).SendString("Recette introuvable")
	}

	return c.Status(200).JSON(recette)
}

// GetRecettesByIngredient retourne toutes les recettes contenant un ingrédient spécifique
func GetRecettesByIngredient(c *fiber.Ctx) error {
	ingredient := c.Params("unit")

	// Rechercher les recettes par ingrédient
	filter := bson.M{"ingredients": bson.M{"$elemMatch": bson.M{"unit": ingredient}}}
	cursor, err := recetteCollection.Find(context.Background(), filter)
	if err != nil {
		log.Printf("Failed to fetch recettes by ingredient: %v", err)
		return c.Status(500).SendString("Erreur lors de la récupération des recettes")
	}
	defer cursor.Close(context.Background())

	// Décoder les recettes
	var recettes []models.Recette
	if err := cursor.All(context.Background(), &recettes); err != nil {
		log.Printf("Failed to decode recettes: %v", err)
		return c.Status(500).SendString("Erreur lors du décodage des recettes")
	}

	return c.Status(200).JSON(recettes)
}
