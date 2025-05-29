# Go API MongoDB Scrapper

[![CI](https://github.com/maxime-louis14/go_api_mongo_scrapper/workflows/Continuous%20Integration/badge.svg)](https://github.com/maxime-louis14/go_api_mongo_scrapper/actions/workflows/ci.yml)
[![CD](https://github.com/maxime-louis14/go_api_mongo_scrapper/workflows/Continuous%20Deployment/badge.svg)](https://github.com/maxime-louis14/go_api_mongo_scrapper/actions/workflows/cd.yml)
[![Release](https://github.com/maxime-louis14/go_api_mongo_scrapper/workflows/Release/badge.svg)](https://github.com/maxime-louis14/go_api_mongo_scrapper/actions/workflows/release.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/maxime-louis14/go_api_mongo_scrapper)](https://goreportcard.com/report/github.com/maxime-louis14/go_api_mongo_scrapper)

Une API REST en Go avec MongoDB et un scraper de recettes performant utilisant des goroutines.

## Fonctionnalités

- **API REST** : Serveur web avec Fiber framework
- **Base de données** : MongoDB avec Docker
- **Scraper performant** : Scraping parallèle avec goroutines
- **Tests complets** : Tests unitaires avec couverture de code
- **CI/CD automatisé** : Pipeline GitHub Actions
- **Docker** : Containerisation complète
- **Cross-platform** : Binaires pour Linux, Windows, macOS

## Architecture

```
├── controllers/     # Contrôleurs API
├── database/       # Configuration MongoDB
├── models/         # Modèles de données
├── routes/         # Routes API
├── responses/      # Réponses API
├── scraper/        # Module de scraping
│   ├── scraper.go     # Code principal
│   ├── scraper_test.go # Tests unitaires
│   └── README_TESTS.md # Documentation tests
├── docs/           # Documentation
├── .github/        # Workflows CI/CD
└── Makefile        # Commandes de build
```

## Installation

### Prérequis

- Go 1.20+
- Docker & Docker Compose
- Make (optionnel)

### Installation rapide

```bash
# Cloner le repository
git clone https://github.com/maxime-louis14/go_api_mongo_scrapper.git
cd go_api_mongo_scrapper

# Installer les dépendances
go mod download

# Démarrer MongoDB avec Docker
docker-compose up -d

# Lancer l'API
go run main.go
```

### Avec Make

```bash
# Installation complète
make deps

# Lancer les tests
make test

# Compiler
make build-all

# Lancer l'API
make run-server
```

## Utilisation

### API REST

L'API est disponible sur `http://localhost:8080`

```bash
# Vérifier l'état de l'API
curl http://localhost:8080/health

# Endpoints disponibles
GET    /recipes          # Liste des recettes
POST   /recipes          # Créer une recette
GET    /recipes/:id      # Récupérer une recette
PUT    /recipes/:id      # Modifier une recette
DELETE /recipes/:id      # Supprimer une recette
```

### Scraper

```bash
# Lancer le scraper
cd scraper
go run scraper.go

# Ou avec Make
make run
```

Le scraper utilise des goroutines pour un scraping parallèle performant :
- 10 workers simultanés
- Protection contre les race conditions
- Statistiques en temps réel
- Gestion d'erreurs robuste

## Tests

### Exécution des tests

```bash
# Tests unitaires
make test

# Tests avec race detection
make test-verbose

# Rapport de couverture
make test-coverage

# Benchmarks
make benchmark
```

### Couverture de code

Le projet maintient une couverture de tests de **22.6%** avec :
- 12 tests unitaires
- 2 benchmarks
- Tests de concurrence
- Tests de validation

Voir [README_TESTS.md](scraper/README_TESTS.md) pour plus de détails.

## CI/CD Pipeline

Le projet utilise GitHub Actions pour l'automatisation :

### Continuous Integration (CI)

Déclenché sur push/PR vers `main` et `develop` :

- **Code Quality** : Formatage, linting, analyse statique
- **Tests** : Tests unitaires avec race detection
- **Build** : Compilation cross-platform
- **Security** : Scan de sécurité avec Gosec
- **Docker** : Build et test des images

### Continuous Deployment (CD)

- **Staging** : Déploiement automatique sur push vers `main`
- **Production** : Déploiement sur tags `v*`
- **Rollback** : Rollback automatique en cas d'échec

### Release

Création automatique de releases avec :
- Binaires multi-plateformes
- Images Docker multi-architecture
- Changelog automatique
- Assets GitHub

Voir [docs/CICD.md](docs/CICD.md) pour la documentation complète.

## Docker

### Images disponibles

```bash
# Dernière version
docker pull ghcr.io/maxime-louis14/go_api_mongo_scrapper:latest

# Version spécifique
docker pull ghcr.io/maxime-louis14/go_api_mongo_scrapper:v1.0.0
```

### Utilisation

```bash
# Avec Docker Compose (recommandé)
docker-compose up

# Ou manuellement
docker run -p 8080:8080 ghcr.io/maxime-louis14/go_api_mongo_scrapper:latest
```

## Développement

### Commandes Make

```bash
make help              # Afficher l'aide
make ci                # Pipeline CI local
make ci-full           # CI avec couverture et benchmarks
make docker-build      # Construire l'image Docker
make release VERSION=v1.0.0  # Créer une release
```

### Workflow de développement

1. **Créer une branche feature**
   ```bash
   git checkout -b feature/nouvelle-fonctionnalite
   ```

2. **Développer et tester**
   ```bash
   make test
   make ci
   ```

3. **Créer une Pull Request**
   - Le CI s'exécute automatiquement
   - Tous les checks doivent passer

4. **Merge vers main**
   - Déploiement automatique en staging

5. **Créer une release**
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

## Performance

### Scraper

- **Parallélisme** : 10 goroutines simultanées
- **Vitesse** : ~64 recettes en 3.2 secondes
- **Mémoire** : Optimisé avec channels et sync.Pool
- **Robustesse** : Gestion d'erreurs et timeouts

### API

- **Framework** : Fiber (Express-like pour Go)
- **Base de données** : MongoDB avec indexation
- **Middleware** : CORS, logging, compression
- **Performance** : ~10k req/s en conditions optimales

## Monitoring

### Métriques disponibles

- **Health check** : `/health` endpoint
- **Logs** : Structured logging avec niveaux
- **Métriques** : Temps de réponse, erreurs, throughput

### Alertes

- **CI/CD** : Notifications automatiques sur échecs
- **Dependabot** : Alertes de sécurité
- **GitHub Security** : Scan de vulnérabilités

## Contribution

1. Fork le projet
2. Créer une branche feature (`git checkout -b feature/AmazingFeature`)
3. Commit les changements (`git commit -m 'Add some AmazingFeature'`)
4. Push vers la branche (`git push origin feature/AmazingFeature`)
5. Ouvrir une Pull Request

### Standards de code

- **Formatage** : `gofmt` obligatoire
- **Linting** : `golangci-lint` sans erreurs
- **Tests** : Couverture minimale de 80%
- **Documentation** : Commentaires Go standard

## Licence

Ce projet est sous licence MIT. Voir le fichier [LICENSE](LICENSE) pour plus de détails.

## Support

- **Issues** : [GitHub Issues](https://github.com/maxime-louis14/go_api_mongo_scrapper/issues)
- **Discussions** : [GitHub Discussions](https://github.com/maxime-louis14/go_api_mongo_scrapper/discussions)
- **Documentation** : [docs/](docs/)

## Roadmap

- [ ] Authentification JWT
- [ ] Rate limiting
- [ ] Cache Redis
- [ ] Métriques Prometheus
- [ ] Dashboard Grafana
- [ ] Tests d'intégration E2E
- [ ] Déploiement Kubernetes

## Deployment

Pour activer le scraper fait.
```bash
  sudo docker compose up -d --build
```


# Bonjour, je suis Maxime ! Voici un projet demandé par mon école NWS 👋

## Consignes

Le restaurant Hótwings souhaite développer son activité avec de la vente en livraison. Le restaurant mise sur ce nouveau service de la façon suivante: une carte très étendue.

Pour plaire à tous les goûts, le restaurant vous a demandé de développer une API permettant de proposer beaucoup de plats et de recettes.

Vous allez devoir concevoir cette API, mais aussi devoir l'alimenter ! Le restaurant aimerait que vous récupériez les recettes depuis le site https://www.allrecipes.com/

Dans un soucis de benchmark, vous avez promis au client d'implémenter 2 bases de données différentes. Vous devez concevoir votre API avec une base de données SQL et NoSQL. Votre API doit pouvoir fonctionner avec n'importe quelle base de données, l'une sans l'autre.

Afin de vous assurer du fonctionnement de votre produit, vous veillerez à ce qu'un Swagger soit mis en place.

Votre scrapper sera capable de générer un fichier JSON contenant toutes les informations scrappées. Une route sur votre API vous permettra d'importer les nouvelles données dans la base de données choisie par l'utilisateur.

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

Vous pouvez retrouver l'API API_golang_Mysql et le scrapper_go

[![Golang scrapper_go](https://img.shields.io/badge/GitHub-100000?style=for-the-badge&logo=github&logoColor=white)](https://github.com/le-veilleur/scrapper_go)
[![API_golang_Mysql](https://img.shields.io/badge/GitHub-100000?style=for-the-badge&logo=github&logoColor=white)](https://github.com/le-veilleur/go_api__scrapper_mysql_docker)


