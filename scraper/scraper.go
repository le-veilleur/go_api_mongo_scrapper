package main

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gocolly/colly"
)

type Recipe struct {
	Name         string        `json:"name"`
	Page         string        `json:"page"`
	Image        string        `json:"image"`
	Ingredients  []Ingredient  `json:"ingredients"`
	Instructions []Instruction `json:"instructions"`
}

type Ingredient struct {
	Quantity string `json:"quantity"`
	Unit     string `json:"unit"`
}

type Instruction struct {
	Number      string `json:"number"`
	Description string `json:"description"`
}

func main() {
	startTime := time.Now() // Début du chronomètre global
	totalRequests := 0      // Compteur de requêtes

	log.Println("Démarrage du script de scraping...")

	// Créer une instance de collecteur
	c := colly.NewCollector(
		colly.Async(true), // Activer le mode asynchrone pour améliorer les performances
	)

	// Limiter le nombre de requêtes simultanées et ajouter un délai
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 10,
		Delay:       100 * time.Millisecond,
	})

	// Mesurer le temps pris par chaque requête
	c.OnRequest(func(r *colly.Request) {
		totalRequests++ // Incrémenter le compteur de requêtes
		log.Printf("Début de la requête vers %s\n", r.URL)
		r.Ctx.Put("start_time", time.Now())
	})

	c.OnResponse(func(r *colly.Response) {
		start := r.Ctx.GetAny("start_time").(time.Time)
		duration := time.Since(start)
		log.Printf("Réponse reçue de %s en %v\n", r.Request.URL, duration)
	})

	var recipes []Recipe

	// Sélectionner les liens de recette et visiter chaque page de recette
	c.OnHTML("div.mntl-taxonomysc-article-list-group .mntl-card", func(e *colly.HTMLElement) {
		page := e.Request.AbsoluteURL(e.Attr("href")) // Obtenir l'URL absolue
		title := e.ChildText("span.card__title-text")
		image := e.ChildAttr("img", "data-src")

		recipe := Recipe{Name: title, Page: page, Image: image}
		recipes = append(recipes, recipe)

		log.Printf("La recette '%s' a été collectée\n", recipe.Name)

		// Visiter la page de recette
		err := c.Visit(page)
		if err != nil {
			log.Println("Erreur lors de la visite de la page de recette: ", err)
		}
	})

	c.OnHTML("div.mntl-structured-ingredients", func(e *colly.HTMLElement) {
		ingredients := []Ingredient{}

		e.ForEach("li.mntl-structured-ingredients__list-item", func(_ int, ingr *colly.HTMLElement) {
			quantity := ingr.ChildText("span[data-ingredient-quantity=true]")
			unit := ingr.ChildText("span[data-ingredient-unit=true]")
			ingredients = append(ingredients, Ingredient{Quantity: quantity, Unit: unit})
		})
		if len(recipes) > 0 {
			recipes[len(recipes)-1].Ingredients = ingredients
		}
	})

	c.OnHTML("div.recipe__steps", func(e *colly.HTMLElement) {
		instructions := []Instruction{}
		e.ForEach("li", func(i int, inst *colly.HTMLElement) {
			number := strconv.Itoa(i + 1)
			description := inst.ChildText("p.mntl-sc-block")
			instructions = append(instructions, Instruction{Number: number, Description: description})
		})
		if len(recipes) > 0 {
			recipes[len(recipes)-1].Instructions = instructions
		}
	})

	// Enregistrer la recette dans le fichier JSON
	c.OnScraped(func(r *colly.Response) {
		content, err := json.MarshalIndent(recipes, "", "  ") // Ajouter un formatage lisible
		if err != nil {
			log.Println("Erreur lors de la sérialisation des recettes: ", err)
			return
		}
		fileName := "data.json"
		err = os.WriteFile(fileName, content, 0644)
		if err != nil {
			log.Println("Erreur lors de l'enregistrement des recettes: ", err)
			return
		}
		log.Printf("Toutes les recettes ont été enregistrées dans le fichier '%s'\n", fileName)
	})

	// Démarrer le scraping
	err := c.Visit("https://www.allrecipes.com/recipes/16369/soups-stews-and-chili/soup/")
	if err != nil {
		log.Println("Erreur lors de la visite du site principal: ", err)
		return
	}

	c.Wait() // Attendre la fin de toutes les requêtes asynchrones

	totalDuration := time.Since(startTime) // Fin du chronomètre global
	log.Printf("Script terminé en %v\n", totalDuration)
	log.Printf("Nombre total de requêtes effectuées: %d\n", totalRequests)
}
