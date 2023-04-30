package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/maxime-louis14/api-golang/configs"
	"github.com/maxime-louis14/api-golang/models"
	"github.com/maxime-louis14/api-golang/responses"
	"go.mongodb.org/mongo-driver/mongo"
)

var recetteCollection *mongo.Collection = configs.GetCollection(configs.DB, "recette")

func PostRecettes(c *fiber.Ctx) error {

	fmt.Println("Début de la fonction PostRecettes")

	// Ouvrir le fichier JSON
	file, err := os.OpenFile("scraper/recettes.json", os.O_RDONLY, 0644)
	if err != nil {
		fmt.Println("Erreur lors de l'ouverture du fichier JSON :", err)
		return c.Status(http.StatusInternalServerError).JSON(responses.RecetteResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	defer file.Close()

	// Ajouter un log pour afficher le contenu du fichier
	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("Erreur lors de la lecture du fichier JSON :", err)
		return c.Status(http.StatusInternalServerError).JSON(responses.RecetteResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	fmt.Println("Contenu du fichier JSON :", string(fileContents))

	var recettes []models.Recette
	err = json.NewDecoder(file).Decode(&recettes)
	if err != nil {
		fmt.Println("Erreur lors de la décodage du JSON :", err)
		return c.Status(http.StatusBadRequest).JSON(responses.RecetteResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": "invalid JSON format"}})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Convertir les recettes en un tableau d'interface{}
	var recettesInterface []interface{}
	for _, recette := range recettes {
		recettesInterface = append(recettesInterface, recette)
	}

	// Insérer les recettes dans la base de données
	result, err := recetteCollection.InsertMany(ctx, recettesInterface)
	if err != nil {
		fmt.Println("Erreur lors de l'insertion des recettes :", err)
		return c.Status(http.StatusInternalServerError).JSON(responses.RecetteResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	// Renvoyer le nombre de recettes insérées
	fmt.Println("Nombre de recettes insérées :", len(result.InsertedIDs))
	return c.Status(http.StatusOK).JSON(responses.RecetteResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"count": len(result.InsertedIDs)}})
}
