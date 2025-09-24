# 📊 Système de Métriques et Logging

Ce document décrit le système de métriques et de logging avancé implémenté dans l'API Go MongoDB Scrapper.

## 🎯 Vue d'ensemble

Le système de métriques fournit une visibilité complète sur :
- **Performance** : Latence, débit, temps de réponse
- **Utilisation** : Requêtes par endpoint, codes de statut
- **Base de données** : Opérations MongoDB, temps d'exécution
- **Système** : Mémoire, goroutines, uptime
- **Métier** : Nombre de recettes, ingrédients, etc.

## 🏗️ Architecture

### Modules

- **`logger/logger.go`** : Module principal de logging et métriques
- **`middleware/logging.go`** : Middleware de logging des requêtes HTTP
- **`main.go`** : Intégration et exposition des métriques

### Types de logs

1. **Logs de requêtes HTTP** : Chaque requête est loggée avec des détails complets
2. **Logs de base de données** : Opérations MongoDB avec timing
3. **Logs d'erreurs** : Gestion centralisée des erreurs
4. **Logs de métriques** : Affichage périodique des statistiques

## 📈 Métriques collectées

### Métriques HTTP

```json
{
  "total_requests": 150,
  "avg_latency_ms": 45.67,
  "requests_by_method": {
    "GET": 120,
    "POST": 30
  },
  "requests_by_path": {
    "/recettes": 50,
    "/health": 20,
    "/metrics": 10
  },
  "status_codes": {
    "200": 140,
    "404": 5,
    "500": 5
  }
}
```

### Métriques de base de données

```json
{
  "database_operations": {
    "find_all": 25,
    "find_one": 30,
    "batch_insert": 5,
    "ping": 15
  }
}
```

### Métriques système

```json
{
  "memory_alloc_mb": 12.45,
  "memory_sys_mb": 25.67,
  "goroutines": 8,
  "uptime_seconds": 3600
}
```

## 🔧 Configuration

### Logging périodique

Les métriques sont affichées automatiquement toutes les 30 secondes :

```go
logger.StartMetricsLogger(30 * time.Second)
```

### Niveaux de log

- **DEBUG** : Informations détaillées de débogage
- **INFO** : Informations générales
- **WARN** : Avertissements
- **ERROR** : Erreurs

## 📊 Endpoints de métriques

### GET /metrics

Retourne toutes les métriques au format JSON :

```bash
curl http://localhost:8080/metrics
```

**Réponse :**
```json
{
  "timestamp": "2024-01-15T10:30:00Z",
  "uptime_seconds": 3600,
  "total_requests": 150,
  "avg_latency_ms": 45.67,
  "error_count": 5,
  "error_rate_percent": 3.33,
  "requests_by_method": {...},
  "requests_by_path": {...},
  "status_codes": {...},
  "database_operations": {...},
  "memory_alloc_mb": 12.45,
  "memory_sys_mb": 25.67,
  "goroutines": 8,
  "last_request": "2024-01-15T10:29:45Z"
}
```

## 🔍 Exemples de logs

### Log de requête HTTP

```json
{
  "timestamp": "2024-01-15T10:30:00Z",
  "level": "INFO",
  "message": "Début de requête",
  "service": "go-api-mongo-scrapper",
  "request_id": "a1b2c3d4",
  "method": "GET",
  "path": "/recettes",
  "status_code": 200,
  "latency": "45.67ms",
  "user_agent": "curl/7.68.0",
  "ip": "127.0.0.1",
  "duration_ns": 45670000
}
```

### Log d'opération de base de données

```json
{
  "timestamp": "2024-01-15T10:30:00Z",
  "level": "INFO",
  "message": "Récupération de toutes les recettes terminée",
  "service": "go-api-mongo-scrapper",
  "database": "mongodb",
  "operation": "find_all",
  "duration_ns": 12345678,
  "extra": {
    "request_id": "a1b2c3d4",
    "recettes_count": 25
  }
}
```

### Log de métriques périodiques

```json
{
  "timestamp": "2024-01-15T10:30:00Z",
  "level": "INFO",
  "message": "Métriques de l'application",
  "service": "go-api-mongo-scrapper",
  "extra": {
    "uptime_seconds": 3600,
    "total_requests": 150,
    "avg_latency_ms": "45.67",
    "error_count": 5,
    "memory_alloc_mb": "12.45",
    "goroutines": 8
  }
}
```

## 🚀 Utilisation

### Test des métriques

Utilisez le script de test fourni :

```bash
./scripts/test_metrics.sh
```

### Surveillance en temps réel

Pour surveiller les logs en temps réel :

```bash
# Logs de l'API
docker logs -f <container_id>

# Ou si vous lancez directement
go run main.go
```

### Intégration avec des outils de monitoring

Les logs JSON peuvent être facilement intégrés avec :
- **ELK Stack** (Elasticsearch, Logstash, Kibana)
- **Grafana + Loki**
- **Fluentd**
- **Prometheus** (avec un adaptateur)

## 🔧 Personnalisation

### Modifier la fréquence des métriques

```go
// Toutes les 60 secondes au lieu de 30
logger.StartMetricsLogger(60 * time.Second)
```

### Ajouter des métriques personnalisées

```go
// Dans votre code
logger.LogInfo("Opération personnalisée", map[string]interface{}{
    "custom_metric": "valeur",
    "count": 42,
})
```

### Modifier le format des logs

Les logs sont en JSON par défaut. Pour changer le format, modifiez la fonction `logJSON` dans `logger/logger.go`.

## 📋 Bonnes pratiques

1. **Utilisez des request_id** : Chaque requête a un ID unique pour le tracing
2. **Loggez les erreurs** : Toutes les erreurs sont automatiquement loggées
3. **Surveillez les métriques** : Vérifiez régulièrement les métriques de performance
4. **Configurez les alertes** : Mettez en place des alertes sur les métriques critiques
5. **Archivez les logs** : Conservez les logs pour l'analyse historique

## 🐛 Dépannage

### Problèmes courants

1. **Métriques vides** : Vérifiez que l'API reçoit des requêtes
2. **Logs manquants** : Vérifiez la configuration des middlewares
3. **Performance dégradée** : Surveillez les métriques de latence et mémoire

### Debug

Activez les logs de debug :

```go
logger.LogInfo("Debug info", map[string]interface{}{
    "debug": true,
    "details": "information de debug",
})
```
