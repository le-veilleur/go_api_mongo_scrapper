package main

import (
	"encoding/json"

	"fmt"

	"os"

	"github.com/gocolly/colly"
)

type data struct {
	Recettes []recettes `json:"recettes"`
}

type recettes struct {
	Name         string `json:"name"`
	Descriptions string `json:"descriptions"`
	Ingredients  string `json:"ingredients"`
	Photos       string `json:"photos"`
	Instructions string `json:"instructions"`
	Page         string `json:"line"`
}

var allRecettes []recettes
var image string
var link string

func main() {

	allData := data{
		Recettes: []recettes{},
	}

	c := colly.NewCollector(
		colly.AllowedDomains("allrecipes.com", "www.allrecipes.com"),
	)

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Scraping:", r.URL)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Status:", r.StatusCode)
	})

	// OnHTML enregistre une fonction. La fonction sera exécutée sur chaque HTML élément correspondant au paramètre
	c.OnHTML("div.mntl-taxonomysc-article-list-group .mntl-card", func(h *colly.HTMLElement) {
		link = h.Attr("href")
		image = h.ChildAttr("img", "data-src")
		c.Visit(link)
	})

	c.OnHTML("article.mntl-article", func(h *colly.HTMLElement) {
		recettes := recettes{
			Name:         h.ChildText("h1.type--lion"),
			Descriptions: h.ChildText("p.article-subheading"),
			Page:         link,
			Photos:       image,
			Ingredients:  h.ChildText("div.mntl-structured-ingredients"),
			Instructions: h.ChildText("div.recipe__steps"),
		}
		fmt.Println(recettes)
		allData.Recettes = append(allData.Recettes, recettes)
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "nError:", err)
	})

	c.Visit("https://www.allrecipes.com/recipes/17562/dinner/")

	content, err := json.Marshal(allData)
	if err != nil {
		fmt.Println(err.Error())
	}

	os.WriteFile("recettes.json", content, 0644)
	fmt.Println("Total recettes: ", len(allRecettes))
}
