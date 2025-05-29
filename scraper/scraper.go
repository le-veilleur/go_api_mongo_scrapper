package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/gocolly/colly"
)

// Variables de versioning injectées lors du build
var (
	version   = "dev"
	gitCommit = "unknown"
	buildTime = "unknown"
)

// BuildInfo contient les informations de build
type BuildInfo struct {
	Version   string `json:"version"`
	GitCommit string `json:"git_commit"`
	BuildTime string `json:"build_time"`
	GoVersion string `json:"go_version"`
	OS        string `json:"os"`
	Arch      string `json:"arch"`
}

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

type RecipeData struct {
	URL   string
	Title string
	Image string
}

type ScrapingStats struct {
	TotalRequests int
	Mutex         sync.Mutex
}

func (s *ScrapingStats) Increment() {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.TotalRequests++
}

func (s *ScrapingStats) Get() int {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	return s.TotalRequests
}

// printVersionInfo affiche les informations de version
func printVersionInfo() {
	fmt.Printf("Go MongoDB Scrapper\n")
	fmt.Printf("Version: %s\n", version)
	fmt.Printf("Git Commit: %s\n", gitCommit)
	fmt.Printf("Build Time: %s\n", buildTime)
	fmt.Printf("Go Version: %s\n", runtime.Version())
	fmt.Printf("OS/Arch: %s/%s\n\n", runtime.GOOS, runtime.GOARCH)
}

// getBuildInfo retourne les informations de build
func getBuildInfo() BuildInfo {
	return BuildInfo{
		Version:   version,
		GitCommit: gitCommit,
		BuildTime: buildTime,
		GoVersion: runtime.Version(),
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
	}
}

// createMainCollector crée et configure le collecteur principal pour la page de liste
func createMainCollector(stats *ScrapingStats, recipeURLs chan<- RecipeData) *colly.Collector {
	collector := colly.NewCollector()
	collector.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 5,
		Delay:       50 * time.Millisecond,
	})

	collector.OnRequest(func(r *colly.Request) {
		stats.Increment()
		log.Printf("Début de la requête principale vers %s\n", r.URL)
	})

	collector.OnHTML("div.mntl-taxonomysc-article-list-group .mntl-card", func(e *colly.HTMLElement) {
		page := e.Request.AbsoluteURL(e.Attr("href"))
		title := e.ChildText("span.card__title-text")
		image := e.ChildAttr("img", "data-src")

		if page != "" && title != "" {
			recipeData := RecipeData{
				URL:   page,
				Title: title,
				Image: image,
			}

			select {
			case recipeURLs <- recipeData:
				log.Printf("URL de recette ajoutée à la queue: '%s'\n", title)
			default:
				log.Printf("Channel plein, recette ignorée: '%s'\n", title)
			}
		}
	})

	return collector
}

// createRecipeCollector crée un collecteur pour scraper une recette individuelle
func createRecipeCollector(stats *ScrapingStats) *colly.Collector {
	collector := colly.NewCollector()
	collector.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 1,
		Delay:       100 * time.Millisecond,
	})

	collector.OnRequest(func(r *colly.Request) {
		stats.Increment()
	})

	return collector
}

// scrapeRecipeDetails configure les handlers pour extraire les détails d'une recette
func scrapeRecipeDetails(collector *colly.Collector, recipe *Recipe, completedRecipes chan<- Recipe) {
	// Collecter les ingrédients
	collector.OnHTML("div.mntl-structured-ingredients", func(e *colly.HTMLElement) {
		var ingredients []Ingredient

		e.ForEach("li.mntl-structured-ingredients__list-item", func(_ int, ingr *colly.HTMLElement) {
			quantity := ingr.ChildText("span[data-ingredient-quantity=true]")
			unit := ingr.ChildText("span[data-ingredient-unit=true]")
			ingredients = append(ingredients, Ingredient{Quantity: quantity, Unit: unit})
		})

		recipe.Ingredients = ingredients
	})

	// Collecter les instructions
	collector.OnHTML("div.recipe__steps", func(e *colly.HTMLElement) {
		var instructions []Instruction
		e.ForEach("li", func(i int, inst *colly.HTMLElement) {
			number := strconv.Itoa(i + 1)
			description := inst.ChildText("p.mntl-sc-block")
			instructions = append(instructions, Instruction{Number: number, Description: description})
		})

		recipe.Instructions = instructions
	})

	// Quand le scraping de la recette est terminé
	collector.OnScraped(func(r *colly.Response) {
		completedRecipes <- *recipe
		log.Printf("Recette complétée: '%s'\n", recipe.Name)
	})
}

// processRecipe traite une recette individuelle dans une goroutine
func processRecipe(recipeData RecipeData, stats *ScrapingStats, completedRecipes chan<- Recipe, wg *sync.WaitGroup) {
	defer wg.Done()

	log.Printf("Traitement de la recette: %s\n", recipeData.Title)

	// Créer un collecteur dédié pour cette recette
	recipeCollector := createRecipeCollector(stats)

	recipe := Recipe{
		Name:  recipeData.Title,
		Page:  recipeData.URL,
		Image: recipeData.Image,
	}

	// Configurer le scraping des détails
	scrapeRecipeDetails(recipeCollector, &recipe, completedRecipes)

	// Visiter la page de la recette
	err := recipeCollector.Visit(recipeData.URL)
	if err != nil {
		log.Printf("Erreur lors de la visite de la page de recette '%s': %v\n", recipeData.Title, err)
	}
}

// startRecipeProcessor démarre la goroutine qui traite les URLs de recettes
func startRecipeProcessor(recipeURLs <-chan RecipeData, completedRecipes chan<- Recipe, stats *ScrapingStats, wg *sync.WaitGroup) {
	go func() {
		const maxWorkers = 10
		semaphore := make(chan struct{}, maxWorkers)

		for recipeData := range recipeURLs {
			wg.Add(1)

			// Acquérir un slot dans le semaphore
			semaphore <- struct{}{}

			go func(rd RecipeData) {
				defer func() { <-semaphore }() // Libérer le slot
				processRecipe(rd, stats, completedRecipes, wg)
			}(recipeData)
		}

		// Attendre que toutes les goroutines se terminent
		wg.Wait()
		close(completedRecipes)
	}()
}

// startRecipeCollector démarre la goroutine qui collecte les recettes terminées
func startRecipeCollector(completedRecipes <-chan Recipe, recipes *[]Recipe, recipesMutex *sync.RWMutex, done chan<- bool) {
	go func() {
		for recipe := range completedRecipes {
			recipesMutex.Lock()
			*recipes = append(*recipes, recipe)
			recipesMutex.Unlock()
		}
		done <- true
	}()
}

// saveRecipesToFile sauvegarde les recettes dans un fichier JSON
func saveRecipesToFile(recipes []Recipe, filename string) error {
	content, err := json.MarshalIndent(recipes, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, content, 0644)
}

// printStats affiche les statistiques finales
func printStats(startTime time.Time, totalRequests int, recipesCount int, filename string) {
	totalDuration := time.Since(startTime)
	log.Printf("Script terminé en %v\n", totalDuration)
	log.Printf("Nombre total de requêtes effectuées: %d\n", totalRequests)
	log.Printf("Nombre de recettes collectées: %d\n", recipesCount)
	log.Printf("Toutes les recettes ont été enregistrées dans le fichier '%s'\n", filename)
}

func main() {
	// Affichage des informations de version
	printVersionInfo()

	startTime := time.Now()
	stats := &ScrapingStats{}

	log.Printf("Démarrage du script de scraping avec goroutines (version %s)...\n", version)
	log.Printf("Build info: %+v\n", getBuildInfo())

	// Channels pour la communication entre goroutines
	recipeURLs := make(chan RecipeData, 100)
	completedRecipes := make(chan Recipe, 100)
	done := make(chan bool)

	// Slice thread-safe pour stocker les recettes
	var recipes []Recipe
	var recipesMutex sync.RWMutex

	// WaitGroup pour synchroniser les goroutines
	var wg sync.WaitGroup

	// Créer le collecteur principal
	mainCollector := createMainCollector(stats, recipeURLs)

	// Démarrer les goroutines de traitement
	startRecipeCollector(completedRecipes, &recipes, &recipesMutex, done)
	startRecipeProcessor(recipeURLs, completedRecipes, stats, &wg)

	// Démarrer le scraping de la page principale
	log.Println("Début du scraping de la page principale...")
	err := mainCollector.Visit("https://www.allrecipes.com/recipes/16369/soups-stews-and-chili/soup/")
	if err != nil {
		log.Printf("Erreur lors de la visite du site principal: %v\n", err)
		return
	}

	// Fermer le channel des URLs après avoir terminé la collecte
	close(recipeURLs)

	// Attendre que toutes les recettes soient collectées
	<-done

	// Sauvegarder les résultats
	filename := "data.json"
	recipesMutex.RLock()
	err = saveRecipesToFile(recipes, filename)
	recipesCount := len(recipes)
	recipesMutex.RUnlock()

	if err != nil {
		log.Printf("Erreur lors de l'enregistrement des recettes: %v\n", err)
		return
	}

	// Afficher les statistiques finales
	printStats(startTime, stats.Get(), recipesCount, filename)

	// Afficher les informations de build dans les stats finales
	log.Printf("Scraping terminé avec la version %s (commit: %s)\n", version, gitCommit)
}
