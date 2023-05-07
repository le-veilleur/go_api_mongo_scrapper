package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gocolly/colly"
)

type Recipe struct {
	Name         string        `json:"name"`
	Page         string        `json:"page"`
	Image        string        `json:"image"`
	Ingredients  []Ingredient  `json:"ingredients"`
	Instructions []Instruction `json:"Instructions"`
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
	// Créer une instance de collecteur
	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Scraping:", r.URL)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Status:", r.StatusCode)
	})

	var recipes []Recipe

	// Sélectionner les liens de recette et visiter chaque page de recette
	c.OnHTML("div.mntl-taxonomysc-article-list-group .mntl-card", func(e *colly.HTMLElement) {
		page := e.Attr("href")
		title := e.ChildText("span.card__title-text")
		image := e.ChildAttr("img", "data-src")

		recipe := Recipe{Name: title, Page: page, Image: image}

		recipes = append(recipes, recipe)

		fmt.Println("La recette", recipe.Name, "a été collectée")

		// Visiter la page de recette
		err := c.Visit(page)
		if err != nil {
			log.Println("Erreur lors de la visite de la page de recette: ", err)
			return
		}
		e.Request.Visit(page)
	})

	c.OnHTML("div.mntl-structured-ingredients", func(e *colly.HTMLElement) {
		ingredients := []Ingredient{}

		// Nameig := e.ChildText("h2.mntl-structured-ingredients__heading")
		// ingredients = append(ingredients, Ingredient{Nameig: Nameig})

		e.ForEach("li.mntl-structured-ingredients__list-item", func(_ int, ingr *colly.HTMLElement) {
			quantityElement := ingr.DOM.Find("span[data-ingredient-quantity=true]")
			unitElement := ingr.DOM.Find("span[data-ingredient-unit=true]")
			quantity := quantityElement.Text()
			unit := unitElement.Text()
			ingredients = append(ingredients, Ingredient{Quantity: quantity, Unit: unit})
		})
		recipes[len(recipes)-1].Ingredients = ingredients
	})

	c.OnHTML("div.recipe__steps", func(e *colly.HTMLElement) {
		instructions := []Instruction{}
		e.ForEach("li", func(i int, inst *colly.HTMLElement) {
			number := strconv.Itoa(i + 1)
			description := inst.ChildText("p.mntl-sc-block")
			instructions = append(instructions, Instruction{Number: number, Description: description})
		})
		recipes[len(recipes)-1].Instructions = instructions
	})

	// Enregistrer la recette dans le fichier JSON
	c.OnScraped(func(r *colly.Response) {
		content, err := json.Marshal(recipes)
		if err != nil {
			fmt.Println(err.Error())
		}
		os.WriteFile("data.json", content, 0644)
		fmt.Println("Toutes les recettes ont été enregistrées dans le fichier 'recettes.json'")
	})

	// Démarrer le scraping
	c.Visit("https://www.allrecipes.com/recipes/16369/soups-stews-and-chili/soup/")

}