# CI/CD Pipeline Documentation

## Vue d'ensemble

Ce projet utilise GitHub Actions pour automatiser l'intégration continue (CI) et le déploiement continu (CD). Le pipeline est composé de trois workflows principaux :

1. **Continuous Integration** (`ci.yml`) - Tests et vérifications
2. **Continuous Deployment** (`cd.yml`) - Déploiement automatique
3. **Release** (`release.yml`) - Création de releases

## Workflows

### 1. Continuous Integration (CI)

**Déclencheurs :**
- Push sur les branches `main` et `develop`
- Pull requests vers `main` et `develop`

**Jobs :**

#### Code Quality
- Vérification du formatage (`gofmt`)
- Analyse statique (`go vet`)
- Linting (`golangci-lint`)

#### Unit Tests
- Exécution des tests unitaires avec race detection
- Génération du rapport de couverture
- Upload vers Codecov

#### Integration Tests
- Tests avec MongoDB en service
- Tests du scraper avec timeout
- Validation de l'intégration complète

#### Build
- Compilation cross-platform (Linux, Windows, macOS)
- Architectures : amd64, arm64
- Upload des artefacts

#### Security
- Scan de sécurité avec Gosec
- Upload des résultats SARIF

#### Docker Build
- Construction de l'image Docker
- Tests de l'image
- Cache optimisé

### 2. Continuous Deployment (CD)

**Déclencheurs :**
- Succès du workflow CI sur `main`
- Push de tags `v*`

**Jobs :**

#### Deploy Staging
- Déploiement automatique en staging
- Push de l'image Docker avec tag `staging`
- Exécution sur push vers `main`

#### Deploy Production
- Déploiement en production
- Création de release GitHub
- Exécution sur tags `v*`

#### Rollback
- Rollback automatique en cas d'échec
- Restauration de la version précédente

### 3. Release

**Déclencheurs :**
- Push de tags `v*`

**Jobs :**

#### Create Release
- Génération automatique du changelog
- Création de la release GitHub
- Support des pre-releases (alpha, beta, rc)

#### Build Binaries
- Compilation pour toutes les plateformes
- Création d'archives (tar.gz, zip)
- Upload des assets vers la release

#### Build Docker
- Construction multi-architecture
- Push vers GitHub Container Registry
- Tags sémantiques

## Configuration

### Variables d'environnement

```yaml
env:
  GO_VERSION: '1.20'
  REGISTRY: ghcr.io
  IMAGE_NAME: ${{ github.repository }}
```

### Secrets requis

- `GITHUB_TOKEN` : Token automatique pour GitHub Actions
- Optionnel : `CODECOV_TOKEN` pour l'upload de couverture

### Cache

Le pipeline utilise plusieurs niveaux de cache :
- **Go modules** : Cache des dépendances Go
- **Docker layers** : Cache des couches Docker
- **Build cache** : Cache de compilation Go

## Utilisation

### Développement

1. **Créer une branche feature :**
   ```bash
   git checkout -b feature/nouvelle-fonctionnalite
   ```

2. **Développer et tester localement :**
   ```bash
   make test
   make build
   ```

3. **Créer une Pull Request :**
   - Le CI s'exécute automatiquement
   - Tous les checks doivent passer

### Déploiement

#### Staging
- Push vers `main` déclenche le déploiement staging
- Automatique après succès du CI

#### Production
1. **Créer un tag de version :**
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

2. **Le pipeline automatique :**
   - Crée la release GitHub
   - Compile les binaires
   - Push l'image Docker
   - Déploie en production

### Versioning

Le projet suit le versioning sémantique :
- `v1.0.0` : Release majeure
- `v1.1.0` : Release mineure
- `v1.1.1` : Patch
- `v1.0.0-alpha.1` : Pre-release

## Monitoring

### Status Badges

Ajoutez ces badges à votre README :

```markdown
![CI](https://github.com/maxime-louis14/api-golang/workflows/Continuous%20Integration/badge.svg)
![CD](https://github.com/maxime-louis14/api-golang/workflows/Continuous%20Deployment/badge.svg)
![Release](https://github.com/maxime-louis14/api-golang/workflows/Release/badge.svg)
```

### Notifications

Le pipeline inclut des notifications automatiques :
- Succès/échec des builds
- Status des déploiements
- Création de releases

## Dépendances

### Dependabot

Configuration automatique pour :
- **Go modules** : Mise à jour hebdomadaire
- **GitHub Actions** : Mise à jour hebdomadaire
- **Docker** : Mise à jour hebdomadaire

### Actions utilisées

- `actions/checkout@v4`
- `actions/setup-go@v4`
- `actions/cache@v3`
- `docker/build-push-action@v5`
- `golangci/golangci-lint-action@v3`
- `codecov/codecov-action@v3`

## Troubleshooting

### Échecs courants

#### Tests qui échouent
```bash
# Exécuter localement
make test-verbose
```

#### Problèmes de formatage
```bash
# Corriger le formatage
make fmt
```

#### Échecs de build
```bash
# Vérifier la compilation
make build
```

#### Problèmes Docker
```bash
# Tester localement
docker build -t test .
docker run --rm test
```

### Logs

- **GitHub Actions** : Onglet "Actions" du repository
- **Docker** : Logs des containers
- **Tests** : Rapports de couverture

## Sécurité

### Scans automatiques

- **Gosec** : Analyse de sécurité du code Go
- **Dependabot** : Alertes de sécurité des dépendances
- **SARIF** : Upload des résultats vers GitHub Security

### Bonnes pratiques

- Pas de secrets dans le code
- Images Docker minimales
- Scans de vulnérabilités
- Mise à jour régulière des dépendances

## Performance

### Optimisations

- **Cache agressif** : Go modules, Docker layers
- **Builds parallèles** : Matrix strategy
- **Images multi-stage** : Réduction de la taille
- **Compilation optimisée** : Flags `-ldflags="-s -w"`

### Métriques

- **Temps de build** : ~5-10 minutes
- **Taille d'image** : ~20-30 MB
- **Couverture de tests** : Objectif >80% 