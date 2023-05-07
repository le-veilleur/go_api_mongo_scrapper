package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/maxime-louis14/api-golang/configs"
	"github.com/maxime-louis14/api-golang/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

var recetteCollection *mongo.Collection = configs.GetCollection(configs.DB, "recettes")
var validate = validator.New()

func PostRecette(c *fiber.Ctx) error {
	// Ouvrir le fichier data.json
	file, err := os.Open("data.json")
	if err != nil {
		return err
	}
	defer file.Close()

	// Lire les données JSON dans un []byte
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	// Décodez les données JSON dans une variable slice de recettes
	var recettes []models.Recette
	err = json.Unmarshal(data, &recettes)
	if err != nil {
		return err
	}

	// Insérer chaque recette dans la collection "recettes" de la base de données "mydb"
	for _, recette := range recettes {
		_, err := recetteCollection.InsertOne(context.Background(), recette)
		if err != nil {
			return err
		}
	}
	fmt.Println(recettes)

	// Réponse HTTP avec un message de succès
	return c.SendString("Recettes ajoutées avec succès")
}

func GetAllRecettes(c *fiber.Ctx) error {
	// Récupérer tous les documents de la collection "recettes"
	cursor, err := recetteCollection.Find(context.Background(), bson.M{})
	if err != nil {
		return err
	}
	defer cursor.Close(context.Background())

	// Itérer sur tous les documents et les stocker dans une variable slice de recettes
	var recettes []models.Recette
	for cursor.Next(context.Background()) {
		var recette models.Recette
		err := cursor.Decode(&recette)
		if err != nil {
			return err
		}
		recettes = append(recettes, recette)
	}

	// Vérifier s'il y a eu une erreur lors de l'itération
	if err := cursor.Err(); err != nil {
		return err
	}

	// Réponse HTTP avec les recettes au format JSON
	return c.JSON(recettes)
}

func GetRecetteByID(c *fiber.Ctx) error {
	// Récupérer l'ID de la recette depuis les paramètres de l'URL
	id := c.Params("id")

	// Vérifier si l'ID est valide
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).SendString("ID de recette invalide")
	}

	// Créer un filtre pour chercher la recette avec l'ID correspondant
	filter := bson.M{"_id": objID}

	// Rechercher la recette dans la base de données
	var recette models.Recette
	err = recetteCollection.FindOne(context.Background(), filter).Decode(&recette)
	if err != nil {
		return c.Status(404).SendString("Recette introuvable")
	}

	// Retourner la recette en JSON
	return c.JSON(recette)
}

func GetRecetteByName(c *fiber.Ctx) error {
	// Récupérer le nom de la recette depuis les paramètres de l'URL
	nomRecette := strings.ReplaceAll(c.Params("name"), "%20", " ")

	// Créer un filtre pour chercher la recette avec le nom correspondant
	filter := bson.M{"name": nomRecette}

	// Rechercher la recette dans la base de données
	var recette models.Recette
	err := recetteCollection.FindOne(context.Background(), filter).Decode(&recette)
	if err != nil {
		return c.Status(404).SendString("Recette introuvable")
	}

	// Retourner la recette en JSON
	return c.JSON(recette)
}

func GetRecettesByIngredient(c *fiber.Ctx) error {
	// Récupérer l'ingrédient depuis les paramètres de l'URL
	ingredient := c.Params("unit")

	// Créer un filtre pour chercher toutes les recettes contenant l'ingrédient correspondant
	filter := bson.M{"ingredients": bson.M{"$elemMatch": bson.M{"unit": ingredient}}}

	// Rechercher les recettes dans la base de données
	var recettes []models.Recette
	cursor, err := recetteCollection.Find(context.Background(), filter)
	if err != nil {
		return c.Status(404).SendString("Aucune recette trouvée")
	}

	// Récupérer toutes les recettes dans un slice
	if err = cursor.All(context.Background(), &recettes); err != nil {
		return c.Status(404).SendString("Aucune recette trouvée")
	}

	// Retourner les recettes en JSON
	return c.JSON(recettes)
}
