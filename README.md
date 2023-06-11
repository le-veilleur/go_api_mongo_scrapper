
## Deployment

Pour activer le scraper fait. Ou (dans le dossier ou se trouve le main.go)
```bash
  docker comose up -d
```


# Bonjour, je suis Maxime ! Voici un projet demandé par mon école NWS 👋

## Consignes

Le restaurant Hótwings souhaite développer son activité avec de la vente en livraison. Le restaurant mise sur ce nouveau service de la façon suivante: une carte très étendue.

Pour plaire à tous les goûts, le restaurant vous a demandé de développer une API permettant de proposer beaucoup de plats et de recettes.

Vous allez devoir concevoir cette API, mais aussi devoir l’alimenter ! Le restaurant aimerait que vous récupériez les recettes depuis le site https://www.allrecipes.com/

Dans un soucis de benchmark, vous avez promis au client d’implémenter 2 bases de données différentes. Vous devez concevoir votre API avec une base de données SQL et NoSQL. Votre API doit pouvoir fonctionner avec n’importe quelle base de données, l’une sans l’autre.

Afin de vous assurer du fonctionnement de votre produit, vous veillerez à ce qu’un Swagger soit mis en place.

Votre scrapper sera capable de générer un fichier JSON contenant toutes les informations scrappées. Une route sur votre API vous permettra d’importer les nouvelles données dans la base de données choisie par l’utilisateur.

## Fonctionnalités attendues

### Fonctionnalités de Lecture

- Lister les recettes ⇒ get
- Lister une recette, ses ingrédients et ses étapes de préparation ⇒ get

### Fonctionnalité de Recherche

- Rechercher une recette par nom
- Rechercher une recette par ingrédient

### Importation de la base de données

- Importer la base de données depuis un fichier JSON dans la base de données choisie

### Outils & Stack

Voici la stack qui vous est **recommandée** pour le projet:

- MySQL / MariaDB pour le SQL
- MongoDB pour le NoSQL

## 🔗 Links

Vous pouvez retrouver l'API API_golang_Mysql et scrapper_go

[![Golang scrapper_go](https://img.shields.io/badge/GitHub-100000?style=for-the-badge&logo=github&logoColor=white)](https://github.com/maxime-louis14/scrapper_go)
[![API_golang_Mysql](https://img.shields.io/badge/GitHub-100000?style=for-the-badge&logo=github&logoColor=white)](https://github.com/maxime-louis14/go_api__scrapper_mysql_docker)


