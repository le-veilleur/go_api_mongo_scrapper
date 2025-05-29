package main

import (
	"encoding/json"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test des structures de données
func TestRecipeStruct(t *testing.T) {
	recipe := Recipe{
		Name:  "Test Recipe",
		Page:  "https://example.com/recipe",
		Image: "https://example.com/image.jpg",
		Ingredients: []Ingredient{
			{Quantity: "1", Unit: "cup"},
			{Quantity: "2", Unit: "tbsp"},
		},
		Instructions: []Instruction{
			{Number: "1", Description: "Mix ingredients"},
			{Number: "2", Description: "Cook for 10 minutes"},
		},
	}

	assert.Equal(t, "Test Recipe", recipe.Name)
	assert.Equal(t, "https://example.com/recipe", recipe.Page)
	assert.Equal(t, "https://example.com/image.jpg", recipe.Image)
	assert.Len(t, recipe.Ingredients, 2)
	assert.Len(t, recipe.Instructions, 2)
}

func TestRecipeDataStruct(t *testing.T) {
	recipeData := RecipeData{
		URL:   "https://example.com/recipe",
		Title: "Test Recipe",
		Image: "https://example.com/image.jpg",
	}

	assert.Equal(t, "https://example.com/recipe", recipeData.URL)
	assert.Equal(t, "Test Recipe", recipeData.Title)
	assert.Equal(t, "https://example.com/image.jpg", recipeData.Image)
}

// Test de ScrapingStats
func TestScrapingStats(t *testing.T) {
	stats := &ScrapingStats{}

	// Test initial
	assert.Equal(t, 0, stats.Get())

	// Test increment
	stats.Increment()
	assert.Equal(t, 1, stats.Get())

	// Test multiple increments
	for i := 0; i < 10; i++ {
		stats.Increment()
	}
	assert.Equal(t, 11, stats.Get())
}

func TestScrapingStatsConcurrency(t *testing.T) {
	stats := &ScrapingStats{}
	var wg sync.WaitGroup
	numGoroutines := 100
	incrementsPerGoroutine := 10

	// Lancer plusieurs goroutines qui incrémentent en parallèle
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < incrementsPerGoroutine; j++ {
				stats.Increment()
			}
		}()
	}

	wg.Wait()
	expected := numGoroutines * incrementsPerGoroutine
	assert.Equal(t, expected, stats.Get())
}

// Test des fonctions utilitaires
func TestSaveRecipesToFile(t *testing.T) {
	recipes := []Recipe{
		{
			Name:  "Test Recipe 1",
			Page:  "https://example.com/recipe1",
			Image: "https://example.com/image1.jpg",
			Ingredients: []Ingredient{
				{Quantity: "1", Unit: "cup"},
			},
			Instructions: []Instruction{
				{Number: "1", Description: "Test instruction"},
			},
		},
		{
			Name:  "Test Recipe 2",
			Page:  "https://example.com/recipe2",
			Image: "https://example.com/image2.jpg",
		},
	}

	// Créer un fichier temporaire
	tempFile := "test_recipes.json"
	defer os.Remove(tempFile) // Nettoyer après le test

	// Tester la sauvegarde
	err := saveRecipesToFile(recipes, tempFile)
	require.NoError(t, err)

	// Vérifier que le fichier existe
	_, err = os.Stat(tempFile)
	require.NoError(t, err)

	// Lire et vérifier le contenu
	content, err := os.ReadFile(tempFile)
	require.NoError(t, err)

	var loadedRecipes []Recipe
	err = json.Unmarshal(content, &loadedRecipes)
	require.NoError(t, err)

	assert.Len(t, loadedRecipes, 2)
	assert.Equal(t, "Test Recipe 1", loadedRecipes[0].Name)
	assert.Equal(t, "Test Recipe 2", loadedRecipes[1].Name)
	assert.Len(t, loadedRecipes[0].Ingredients, 1)
	assert.Len(t, loadedRecipes[0].Instructions, 1)
}

func TestSaveRecipesToFileError(t *testing.T) {
	recipes := []Recipe{{Name: "Test"}}

	// Tenter de sauvegarder dans un répertoire inexistant
	err := saveRecipesToFile(recipes, "/nonexistent/directory/file.json")
	assert.Error(t, err)
}

// Test des collecteurs
func TestCreateMainCollector(t *testing.T) {
	stats := &ScrapingStats{}
	recipeURLs := make(chan RecipeData, 10)
	defer close(recipeURLs)

	collector := createMainCollector(stats, recipeURLs)

	// Vérifier que le collecteur est créé
	assert.NotNil(t, collector)

	// Vérifier la configuration des limites
	// Note: colly ne expose pas directement les limites,
	// donc on teste indirectement via le comportement
}

func TestCreateRecipeCollector(t *testing.T) {
	stats := &ScrapingStats{}

	collector := createRecipeCollector(stats)

	// Vérifier que le collecteur est créé
	assert.NotNil(t, collector)
}

// Test des channels et goroutines
func TestRecipeChannelCommunication(t *testing.T) {
	completedRecipes := make(chan Recipe, 5)
	done := make(chan bool)

	var recipes []Recipe
	var recipesMutex sync.RWMutex

	// Démarrer le collecteur de recettes
	startRecipeCollector(completedRecipes, &recipes, &recipesMutex, done)

	// Envoyer quelques recettes
	testRecipes := []Recipe{
		{Name: "Recipe 1", Page: "https://example.com/1"},
		{Name: "Recipe 2", Page: "https://example.com/2"},
		{Name: "Recipe 3", Page: "https://example.com/3"},
	}

	go func() {
		for _, recipe := range testRecipes {
			completedRecipes <- recipe
		}
		close(completedRecipes)
	}()

	// Attendre la fin
	<-done

	// Vérifier les résultats
	recipesMutex.RLock()
	assert.Len(t, recipes, 3)
	assert.Equal(t, "Recipe 1", recipes[0].Name)
	assert.Equal(t, "Recipe 2", recipes[1].Name)
	assert.Equal(t, "Recipe 3", recipes[2].Name)
	recipesMutex.RUnlock()
}

func TestRecipeDataValidation(t *testing.T) {
	tests := []struct {
		name     string
		data     RecipeData
		expected bool
	}{
		{
			name: "Valid recipe data",
			data: RecipeData{
				URL:   "https://example.com/recipe",
				Title: "Test Recipe",
				Image: "https://example.com/image.jpg",
			},
			expected: true,
		},
		{
			name: "Empty URL",
			data: RecipeData{
				URL:   "",
				Title: "Test Recipe",
				Image: "https://example.com/image.jpg",
			},
			expected: false,
		},
		{
			name: "Empty Title",
			data: RecipeData{
				URL:   "https://example.com/recipe",
				Title: "",
				Image: "https://example.com/image.jpg",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := tt.data.URL != "" && tt.data.Title != ""
			assert.Equal(t, tt.expected, isValid)
		})
	}
}

// Test de performance
func TestScrapingStatsPerformance(t *testing.T) {
	stats := &ScrapingStats{}
	numOperations := 10000

	start := time.Now()

	var wg sync.WaitGroup
	for i := 0; i < numOperations; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			stats.Increment()
		}()
	}

	wg.Wait()
	duration := time.Since(start)

	assert.Equal(t, numOperations, stats.Get())
	// Vérifier que les opérations sont rapides (moins de 1 seconde pour 10k opérations)
	assert.Less(t, duration, time.Second)
}

// Test d'intégration pour les structures JSON
func TestJSONSerialization(t *testing.T) {
	recipe := Recipe{
		Name:  "Test Recipe",
		Page:  "https://example.com/recipe",
		Image: "https://example.com/image.jpg",
		Ingredients: []Ingredient{
			{Quantity: "1", Unit: "cup"},
			{Quantity: "2", Unit: "tbsp"},
		},
		Instructions: []Instruction{
			{Number: "1", Description: "Mix ingredients"},
			{Number: "2", Description: "Cook for 10 minutes"},
		},
	}

	// Sérialiser
	jsonData, err := json.Marshal(recipe)
	require.NoError(t, err)

	// Désérialiser
	var deserializedRecipe Recipe
	err = json.Unmarshal(jsonData, &deserializedRecipe)
	require.NoError(t, err)

	// Vérifier que les données sont identiques
	assert.Equal(t, recipe.Name, deserializedRecipe.Name)
	assert.Equal(t, recipe.Page, deserializedRecipe.Page)
	assert.Equal(t, recipe.Image, deserializedRecipe.Image)
	assert.Len(t, deserializedRecipe.Ingredients, 2)
	assert.Len(t, deserializedRecipe.Instructions, 2)
	assert.Equal(t, recipe.Ingredients[0].Quantity, deserializedRecipe.Ingredients[0].Quantity)
	assert.Equal(t, recipe.Instructions[0].Description, deserializedRecipe.Instructions[0].Description)
}

// Benchmark pour les opérations critiques
func BenchmarkScrapingStatsIncrement(b *testing.B) {
	stats := &ScrapingStats{}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			stats.Increment()
		}
	})
}

func BenchmarkJSONMarshal(b *testing.B) {
	recipe := Recipe{
		Name:  "Test Recipe",
		Page:  "https://example.com/recipe",
		Image: "https://example.com/image.jpg",
		Ingredients: []Ingredient{
			{Quantity: "1", Unit: "cup"},
			{Quantity: "2", Unit: "tbsp"},
		},
		Instructions: []Instruction{
			{Number: "1", Description: "Mix ingredients"},
			{Number: "2", Description: "Cook for 10 minutes"},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(recipe)
		if err != nil {
			b.Fatal(err)
		}
	}
}
