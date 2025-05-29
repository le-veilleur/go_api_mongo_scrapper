# Documentation Docker

Ce document décrit l'utilisation de Docker avec le projet Go API MongoDB Scrapper, incluant le versioning automatique et les configurations optimisées.

## Vue d'ensemble

Le projet utilise une approche Docker multi-services avec :
- **API Server** : Serveur web Fiber avec MongoDB
- **Scraper** : Service de scraping de recettes
- **MongoDB** : Base de données NoSQL
- **MongoDB Express** : Interface web pour MongoDB (optionnel)

## Architecture Docker

### Images disponibles

1. **go-api-mongo-scrapper** : API Server principal
2. **go-scraper** : Service de scraping
3. **mongo:7.0** : Base de données MongoDB
4. **mongo-express** : Interface web MongoDB

### Versioning automatique

Toutes les images sont construites avec des informations de versioning :
- **Version** : Tag Git ou version spécifiée
- **Git Commit** : Hash du commit actuel
- **Build Time** : Timestamp de construction
- **Go Version** : Version de Go utilisée
- **OS/Arch** : Système d'exploitation et architecture

## Utilisation

### Construction des images

#### Build complet avec script automatisé
```bash
# Build avec version automatique (depuis Git)
./scripts/build.sh

# Build avec version spécifique
./scripts/build.sh v1.0.0

# Build avec push vers registry
PUSH_TO_REGISTRY=true ./scripts/build.sh v1.0.0
```

#### Build avec Make
```bash
# Build toutes les images
make docker-build VERSION=v1.0.0

# Build API uniquement
make docker-build-api VERSION=v1.0.0

# Build Scraper uniquement
make docker-build-scraper VERSION=v1.0.0
```

#### Build manuel
```bash
# API Server
docker build \
  --build-arg VERSION=v1.0.0 \
  --build-arg GIT_COMMIT=$(git rev-parse --short HEAD) \
  --build-arg BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
  -t go-api-mongo-scrapper:v1.0.0 \
  -f dockerfile .

# Scraper
docker build \
  --build-arg VERSION=v1.0.0 \
  --build-arg GIT_COMMIT=$(git rev-parse --short HEAD) \
  --build-arg BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
  -t go-scraper:v1.0.0 \
  -f scraper/dockerfile .
```

### Exécution des services

#### Avec Docker Compose (recommandé)

```bash
# Démarrer l'API et MongoDB
make docker-run
# ou
docker-compose up api-server mongodb

# Exécuter le scraper
make docker-run-scraper
# ou
docker-compose --profile scraper up scraper

# Démarrer MongoDB Express
make docker-run-tools
# ou
docker-compose --profile tools up mongo-express

# Arrêter tous les services
make docker-stop
# ou
docker-compose down
```

#### Avec Docker directement

```bash
# Réseau Docker
docker network create app-network

# MongoDB
docker run -d \
  --name go-api-mongodb \
  --network app-network \
  -p 27017:27017 \
  -e MONGO_INITDB_ROOT_USERNAME=admin \
  -e MONGO_INITDB_ROOT_PASSWORD=password123 \
  mongo:7.0

# API Server
docker run -d \
  --name go-api-server \
  --network app-network \
  -p 8080:8080 \
  -e MONGODB_URI=mongodb://admin:password123@go-api-mongodb:27017/recipes?authSource=admin \
  go-api-mongo-scrapper:latest

# Scraper (exécution unique)
docker run --rm \
  --name go-scraper \
  --network app-network \
  -v $(pwd)/data:/app/data \
  go-scraper:latest
```

## Configuration

### Variables d'environnement

#### API Server
```bash
PORT=8080                    # Port d'écoute
ENV=production              # Environnement
MONGODB_URI=mongodb://...   # URI MongoDB
LOG_LEVEL=info             # Niveau de log
```

#### Scraper
```bash
SCRAPER_MAX_WORKERS=10              # Nombre de workers
SCRAPER_TIMEOUT=30s                 # Timeout des requêtes
SCRAPER_BASE_URL=https://...        # URL de base
LOG_LEVEL=info                      # Niveau de log
```

#### MongoDB
```bash
MONGO_INITDB_ROOT_USERNAME=admin    # Utilisateur admin
MONGO_INITDB_ROOT_PASSWORD=password # Mot de passe admin
MONGO_INITDB_DATABASE=recipes       # Base par défaut
```

### Volumes

```yaml
volumes:
  mongodb_data:      # Données MongoDB persistantes
  mongodb_config:    # Configuration MongoDB
  scraper_data:      # Données du scraper
  api_logs:          # Logs de l'API
```

### Réseaux

```yaml
networks:
  app-network:       # Réseau interne pour communication inter-services
    subnet: 172.20.0.0/16
```

## Optimisations

### Images multi-stage

Les Dockerfiles utilisent une approche multi-stage :
1. **Builder stage** : Compilation avec Go complet
2. **Runtime stage** : Image Alpine minimale

### Sécurité

- **Utilisateur non-root** : Exécution avec utilisateur dédié
- **Images minimales** : Alpine Linux pour réduire la surface d'attaque
- **Certificats CA** : Inclus pour les connexions HTTPS
- **Health checks** : Surveillance automatique de l'état

### Performance

- **Cache Docker** : Optimisation des layers
- **Binaires statiques** : CGO_ENABLED=0 pour portabilité
- **Compression** : Flags -ldflags="-s -w" pour réduire la taille

## Monitoring et Debugging

### Health Checks

```bash
# Vérifier l'état de l'API
curl http://localhost:8080/health

# Vérifier les informations de version
curl http://localhost:8080/version

# Health check Docker
docker inspect go-api-server | grep -A 10 Health
```

### Logs

```bash
# Logs de tous les services
make docker-logs

# Logs d'un service spécifique
docker-compose logs -f api-server
docker-compose logs -f mongodb
docker-compose logs scraper

# Logs en temps réel
docker logs -f go-api-server
```

### Debugging

```bash
# Entrer dans un container
docker exec -it go-api-server sh
docker exec -it go-api-mongodb mongosh

# Inspecter les images
docker inspect go-api-mongo-scrapper:latest
docker history go-api-mongo-scrapper:latest

# Vérifier les ressources
docker stats
docker system df
```

## Profils Docker Compose

### Profil par défaut
Services démarrés automatiquement :
- `api-server`
- `mongodb`

### Profil `scraper`
```bash
docker-compose --profile scraper up
```
Services additionnels :
- `scraper`

### Profil `tools`
```bash
docker-compose --profile tools up
```
Services additionnels :
- `mongo-express` (http://localhost:8081)

### Tous les profils
```bash
docker-compose --profile scraper --profile tools up
```

## Déploiement

### Environnements

#### Développement
```bash
# Variables d'environnement
export VERSION=dev
export ENV=development
export LOG_LEVEL=debug

# Démarrage
docker-compose up
```

#### Staging
```bash
# Variables d'environnement
export VERSION=$(git describe --tags)
export ENV=staging
export LOG_LEVEL=info

# Démarrage
docker-compose -f docker-compose.yml -f docker-compose.staging.yml up
```

#### Production
```bash
# Variables d'environnement
export VERSION=v1.0.0
export ENV=production
export LOG_LEVEL=warn

# Démarrage
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d
```

### Registry

#### Push vers GitHub Container Registry
```bash
# Login
echo $GITHUB_TOKEN | docker login ghcr.io -u USERNAME --password-stdin

# Tag et push
docker tag go-api-mongo-scrapper:v1.0.0 ghcr.io/maxime-louis14/go-api-mongo-scrapper:v1.0.0
docker push ghcr.io/maxime-louis14/go-api-mongo-scrapper:v1.0.0

# Ou avec le script
PUSH_TO_REGISTRY=true ./scripts/build.sh v1.0.0
```

#### Pull depuis le registry
```bash
# Pull des images
docker pull ghcr.io/maxime-louis14/go-api-mongo-scrapper:latest
docker pull ghcr.io/maxime-louis14/go-scraper:latest

# Utilisation
docker run -p 8080:8080 ghcr.io/maxime-louis14/go-api-mongo-scrapper:latest
```

## Troubleshooting

### Problèmes courants

#### Port déjà utilisé
```bash
# Vérifier les ports utilisés
netstat -tulpn | grep :8080
lsof -i :8080

# Arrêter les services
docker-compose down
```

#### Problèmes de connexion MongoDB
```bash
# Vérifier les logs MongoDB
docker logs go-api-mongodb

# Tester la connexion
docker exec go-api-mongodb mongosh --eval "db.adminCommand('ping')"

# Vérifier les variables d'environnement
docker exec go-api-server env | grep MONGODB
```

#### Images corrompues
```bash
# Nettoyer les images
docker system prune -a

# Rebuild complet
make clean
make docker-build
```

#### Problèmes de permissions
```bash
# Vérifier les permissions des volumes
docker exec go-api-server ls -la /app

# Reconstruire avec les bonnes permissions
docker-compose down
docker-compose build --no-cache
docker-compose up
```

### Commandes utiles

```bash
# Informations système Docker
docker info
docker version

# Espace disque utilisé
docker system df

# Nettoyer l'espace
docker system prune -a --volumes

# Inspecter un container
docker inspect go-api-server

# Copier des fichiers
docker cp go-api-server:/app/logs ./logs

# Exporter/Importer des images
docker save go-api-mongo-scrapper:latest | gzip > api-server.tar.gz
gunzip -c api-server.tar.gz | docker load
```

## Sécurité

### Bonnes pratiques

1. **Utilisateurs non-root** : Tous les containers s'exécutent avec des utilisateurs dédiés
2. **Secrets** : Utilisation de variables d'environnement pour les credentials
3. **Réseaux isolés** : Communication inter-services via réseau Docker privé
4. **Images minimales** : Alpine Linux pour réduire les vulnérabilités
5. **Scans de sécurité** : Intégration avec les outils CI/CD

### Scan de vulnérabilités

```bash
# Avec Docker Scout (si disponible)
docker scout cves go-api-mongo-scrapper:latest

# Avec Trivy
trivy image go-api-mongo-scrapper:latest

# Avec Clair (si configuré)
clair-scanner go-api-mongo-scrapper:latest
``` 