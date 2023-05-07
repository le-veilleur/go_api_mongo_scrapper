package models

type Recette struct {
	Name         string        `json:"name" swagger:"description(Nom de la recette)"`
	Page         string        `json:"page" swagger:"description(URL de la page de la recette)"`
	Image        string        `json:"image" swagger:"description(URL de l'image de la recette)"`
	Ingredients  []Ingredient  `json:"ingredients" swagger:"description(Liste des ingrédients de la recette)"`
	Instructions []Instruction `json:"Instructions" swagger:"description(Liste des instructions de la recette)"`
}

type Ingredient struct {
	Quantity string `json:"quantity" swagger:"description(Quantité de l'ingrédient)"`
	Unit     string `json:"unit" swagger:"description(Unité de mesure de l'ingrédient)"`
}

type Instruction struct {
	Number      string `json:"number" swagger:"description(Numéro de l'instruction)"`
	Description string `json:"description" swagger:"description(Description de l'instruction)"`
}
