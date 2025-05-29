# Configuration

Ce document décrit toutes les variables de configuration disponibles pour l'application.

## Variables d'environnement

### Application

| Variable | Description | Valeur par défaut | Requis |
|----------|-------------|-------------------|---------|
| `PORT` | Port d'écoute du serveur | `8080` | Non |
| `ENV` | Environnement d'exécution | `development` | Non |

### Base de données

| Variable | Description | Valeur par défaut | Requis |
|----------|-------------|-------------------|---------|
| `MONGODB_URI` | URI de connexion MongoDB | `mongodb://localhost:27017/recipes` | Oui |
| `MONGODB_DATABASE` | Nom de la base de données | `recipes` | Non |
| `MONGODB_COLLECTION` | Nom de la collection | `recipes` | Non |

### Scraper

| Variable | Description | Valeur par défaut | Requis |
|----------|-------------|-------------------|---------|
| `SCRAPER_MAX_WORKERS` | Nombre de workers parallèles | `10` | Non |
| `SCRAPER_TIMEOUT` | Timeout des requêtes | `30s` | Non |
| `SCRAPER_BASE_URL` | URL de base pour le scraping | `https://www.allrecipes.com` | Non |

### Logs

| Variable | Description | Valeur par défaut | Requis |
|----------|-------------|-------------------|---------|
| `LOG_LEVEL` | Niveau de log (debug, info, warn, error) | `info` | Non |
| `LOG_FORMAT` | Format des logs (json, text) | `json` | Non |

### Docker

| Variable | Description | Valeur par défaut | Requis |
|----------|-------------|-------------------|---------|
| `DOCKER_REGISTRY` | Registry Docker | `ghcr.io` | Non |
| `DOCKER_IMAGE_NAME` | Nom de l'image Docker | `go-api-mongo-scrapper` | Non |

### CI/CD

| Variable | Description | Valeur par défaut | Requis |
|----------|-------------|-------------------|---------|
| `CI_ENVIRONMENT` | Environnement CI | `local` | Non |
| `CODECOV_TOKEN` | Token Codecov pour la couverture | - | Non |

### Sécurité

| Variable | Description | Valeur par défaut | Requis |
|----------|-------------|-------------------|---------|
| `JWT_SECRET` | Secret pour les tokens JWT | - | Oui (production) |
| `API_KEY` | Clé API pour l'authentification | - | Non |

### Monitoring

| Variable | Description | Valeur par défaut | Requis |
|----------|-------------|-------------------|---------|
| `HEALTH_CHECK_INTERVAL` | Intervalle des health checks | `30s` | Non |
| `METRICS_ENABLED` | Activer les métriques | `true` | Non |

## Fichiers de configuration

### .env (local)

Créez un fichier `.env` à la racine du projet :

```bash
# Configuration de l'application
PORT=8080
ENV=development

# Configuration MongoDB
MONGODB_URI=mongodb://localhost:27017/recipes
MONGODB_DATABASE=recipes
MONGODB_COLLECTION=recipes

# Configuration du scraper
SCRAPER_MAX_WORKERS=10
SCRAPER_TIMEOUT=30s
SCRAPER_BASE_URL=https://www.allrecipes.com

# Configuration des logs
LOG_LEVEL=info
LOG_FORMAT=json
```

### docker-compose.yml

Le fichier `docker-compose.yml` contient la configuration pour l'environnement Docker :

```yaml
version: '3.8'
services:
  mongodb:
    image: mongo:latest
    environment:
      MONGO_INITDB_ROOT_USERNAME: admin
      MONGO_INITDB_ROOT_PASSWORD: password
    ports:
      - "27017:27017"
    volumes:
      - mongodb_data:/data/db

  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - MONGODB_URI=mongodb://admin:password@mongodb:27017/recipes?authSource=admin
      - ENV=production
    depends_on:
      - mongodb

volumes:
  mongodb_data:
```

## Configuration par environnement

### Développement

```bash
ENV=development
LOG_LEVEL=debug
MONGODB_URI=mongodb://localhost:27017/recipes
SCRAPER_MAX_WORKERS=5
```

### Test

```bash
ENV=test
LOG_LEVEL=warn
MONGODB_URI=mongodb://localhost:27017/recipes_test
SCRAPER_MAX_WORKERS=2
```

### Staging

```bash
ENV=staging
LOG_LEVEL=info
MONGODB_URI=mongodb://staging-db:27017/recipes
SCRAPER_MAX_WORKERS=10
METRICS_ENABLED=true
```

### Production

```bash
ENV=production
LOG_LEVEL=warn
MONGODB_URI=mongodb://prod-db:27017/recipes
SCRAPER_MAX_WORKERS=20
METRICS_ENABLED=true
JWT_SECRET=your_secure_jwt_secret
```

## Validation de la configuration

L'application valide automatiquement la configuration au démarrage :

- **Variables requises** : Vérification de la présence
- **Formats** : Validation des URLs, durées, etc.
- **Valeurs** : Vérification des plages de valeurs
- **Connexions** : Test de connectivité aux services externes

## Secrets et sécurité

### Variables sensibles

Les variables suivantes contiennent des informations sensibles :

- `MONGODB_URI` (si elle contient des credentials)
- `JWT_SECRET`
- `API_KEY`
- `CODECOV_TOKEN`

### Bonnes pratiques

1. **Ne jamais commiter** les fichiers `.env` avec des secrets
2. **Utiliser des secrets managers** en production (AWS Secrets Manager, Azure Key Vault, etc.)
3. **Rotation régulière** des secrets
4. **Principe du moindre privilège** pour les accès

### GitHub Secrets

Pour le CI/CD, configurez ces secrets dans GitHub :

```bash
# Repository Settings > Secrets and variables > Actions

CODECOV_TOKEN=your_codecov_token
DOCKER_REGISTRY_TOKEN=your_registry_token
MONGODB_URI_TEST=mongodb://test:test@localhost:27017/test
```

## Troubleshooting

### Erreurs courantes

#### MongoDB connection failed

```bash
# Vérifier que MongoDB est démarré
docker-compose up mongodb

# Vérifier l'URI de connexion
echo $MONGODB_URI
```

#### Port already in use

```bash
# Changer le port
export PORT=8081

# Ou tuer le processus existant
lsof -ti:8080 | xargs kill -9
```

#### Invalid configuration

```bash
# Vérifier les variables d'environnement
env | grep -E "(MONGODB|PORT|ENV)"

# Valider la configuration
make validate-config
```

### Logs de débogage

Activez les logs de débogage pour diagnostiquer les problèmes :

```bash
export LOG_LEVEL=debug
make run-server
```

## Migration de configuration

### Depuis une version antérieure

Si vous migrez depuis une version antérieure, vérifiez :

1. **Nouvelles variables** : Ajoutez les nouvelles variables requises
2. **Variables obsolètes** : Supprimez les variables non utilisées
3. **Formats changés** : Vérifiez les nouveaux formats
4. **Valeurs par défaut** : Vérifiez les nouvelles valeurs par défaut

### Script de migration

```bash
#!/bin/bash
# migrate-config.sh

echo "Migration de la configuration..."

# Sauvegarder l'ancienne configuration
cp .env .env.backup

# Ajouter les nouvelles variables
echo "METRICS_ENABLED=true" >> .env
echo "HEALTH_CHECK_INTERVAL=30s" >> .env

echo "Migration terminée. Vérifiez le fichier .env"
``` 