package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"
)

// LogLevel définit les niveaux de log
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

// LogEntry représente une entrée de log structurée
type LogEntry struct {
	Timestamp  time.Time              `json:"timestamp"`
	Level      string                 `json:"level"`
	Message    string                 `json:"message"`
	Service    string                 `json:"service"`
	RequestID  string                 `json:"request_id,omitempty"`
	Method     string                 `json:"method,omitempty"`
	Path       string                 `json:"path,omitempty"`
	StatusCode int                    `json:"status_code,omitempty"`
	Latency    string                 `json:"latency,omitempty"`
	UserAgent  string                 `json:"user_agent,omitempty"`
	IP         string                 `json:"ip,omitempty"`
	Database   string                 `json:"database,omitempty"`
	Operation  string                 `json:"operation,omitempty"`
	Duration   int64                  `json:"duration_ns,omitempty"`
	Extra      map[string]interface{} `json:"extra,omitempty"`
}

// MetricsCollector collecte les métriques de l'application
type MetricsCollector struct {
	mu               sync.RWMutex
	TotalRequests    int64            `json:"total_requests"`
	TotalLatencyNs   int64            `json:"total_latency_ns"`
	RequestsByMethod map[string]int64 `json:"requests_by_method"`
	RequestsByPath   map[string]int64 `json:"requests_by_path"`
	StatusCodes      map[int]int64    `json:"status_codes"`
	DatabaseOps      map[string]int64 `json:"database_operations"`
	ErrorCount       int64            `json:"error_count"`
	StartTime        time.Time        `json:"start_time"`
	LastRequestTime  time.Time        `json:"last_request_time"`
	MemoryStats      runtime.MemStats `json:"memory_stats"`
}

var (
	collector *MetricsCollector
	once      sync.Once
)

// GetMetricsCollector retourne l'instance singleton du collecteur de métriques
func GetMetricsCollector() *MetricsCollector {
	once.Do(func() {
		collector = &MetricsCollector{
			RequestsByMethod: make(map[string]int64),
			RequestsByPath:   make(map[string]int64),
			StatusCodes:      make(map[int]int64),
			DatabaseOps:      make(map[string]int64),
			StartTime:        time.Now(),
		}
	})
	return collector
}

// LogRequest enregistre une requête HTTP
func LogRequest(level LogLevel, message, requestID, method, path, userAgent, ip string, statusCode int, latency time.Duration) {
	entry := LogEntry{
		Timestamp:  time.Now(),
		Level:      getLevelString(level),
		Message:    message,
		Service:    "go-api-mongo-scrapper",
		RequestID:  requestID,
		Method:     method,
		Path:       path,
		StatusCode: statusCode,
		Latency:    latency.String(),
		UserAgent:  userAgent,
		IP:         ip,
		Duration:   latency.Nanoseconds(),
	}

	// Mise à jour des métriques
	collector := GetMetricsCollector()
	collector.mu.Lock()
	collector.TotalRequests++
	collector.TotalLatencyNs += latency.Nanoseconds()
	collector.RequestsByMethod[method]++
	collector.RequestsByPath[path]++
	collector.StatusCodes[statusCode]++
	collector.LastRequestTime = time.Now()
	if statusCode >= 400 {
		collector.ErrorCount++
	}
	collector.mu.Unlock()

	// Log structuré
	logJSON(entry)
}

// LogDatabase enregistre une opération de base de données
func LogDatabase(level LogLevel, message, operation, database string, duration time.Duration, extra map[string]interface{}) {
	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     getLevelString(level),
		Message:   message,
		Service:   "go-api-mongo-scrapper",
		Database:  database,
		Operation: operation,
		Duration:  duration.Nanoseconds(),
		Extra:     extra,
	}

	// Mise à jour des métriques
	collector := GetMetricsCollector()
	collector.mu.Lock()
	collector.DatabaseOps[operation]++
	collector.mu.Unlock()

	logJSON(entry)
}

// LogInfo enregistre un message d'information général
func LogInfo(message string, extra map[string]interface{}) {
	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     getLevelString(INFO),
		Message:   message,
		Service:   "go-api-mongo-scrapper",
		Extra:     extra,
	}
	logJSON(entry)
}

// LogError enregistre une erreur
func LogError(message string, err error, extra map[string]interface{}) {
	if extra == nil {
		extra = make(map[string]interface{})
	}
	if err != nil {
		extra["error"] = err.Error()
	}

	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     getLevelString(ERROR),
		Message:   message,
		Service:   "go-api-mongo-scrapper",
		Extra:     extra,
	}

	// Mise à jour des métriques
	collector := GetMetricsCollector()
	collector.mu.Lock()
	collector.ErrorCount++
	collector.mu.Unlock()

	logJSON(entry)
}

// LogMetrics affiche les métriques actuelles
func LogMetrics() {
	collector := GetMetricsCollector()
	collector.mu.RLock()
	defer collector.mu.RUnlock()

	// Mise à jour des stats mémoire
	runtime.ReadMemStats(&collector.MemoryStats)

	// Calcul des moyennes
	avgLatencyMs := float64(0)
	if collector.TotalRequests > 0 {
		avgLatencyMs = float64(collector.TotalLatencyNs) / float64(collector.TotalRequests) / 1e6
	}

	uptime := time.Since(collector.StartTime)

	metrics := map[string]interface{}{
		"timestamp":           time.Now(),
		"uptime_seconds":      uptime.Seconds(),
		"total_requests":      collector.TotalRequests,
		"avg_latency_ms":      fmt.Sprintf("%.2f", avgLatencyMs),
		"error_count":         collector.ErrorCount,
		"error_rate_percent":  fmt.Sprintf("%.2f", float64(collector.ErrorCount)/float64(collector.TotalRequests)*100),
		"requests_by_method":  collector.RequestsByMethod,
		"requests_by_path":    collector.RequestsByPath,
		"status_codes":        collector.StatusCodes,
		"database_operations": collector.DatabaseOps,
		"memory_alloc_mb":     fmt.Sprintf("%.2f", float64(collector.MemoryStats.Alloc)/1024/1024),
		"memory_sys_mb":       fmt.Sprintf("%.2f", float64(collector.MemoryStats.Sys)/1024/1024),
		"goroutines":          runtime.NumGoroutine(),
		"last_request":        collector.LastRequestTime,
	}

	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     getLevelString(INFO),
		Message:   "Métriques de l'application",
		Service:   "go-api-mongo-scrapper",
		Extra:     metrics,
	}

	logJSON(entry)
}

// GetMetricsJSON retourne les métriques au format JSON
func GetMetricsJSON() ([]byte, error) {
	collector := GetMetricsCollector()
	collector.mu.RLock()
	defer collector.mu.RUnlock()

	// Mise à jour des stats mémoire
	runtime.ReadMemStats(&collector.MemoryStats)

	// Calcul des moyennes
	avgLatencyMs := float64(0)
	if collector.TotalRequests > 0 {
		avgLatencyMs = float64(collector.TotalLatencyNs) / float64(collector.TotalRequests) / 1e6
	}

	uptime := time.Since(collector.StartTime)

	metrics := map[string]interface{}{
		"timestamp":           time.Now(),
		"uptime_seconds":      uptime.Seconds(),
		"total_requests":      collector.TotalRequests,
		"avg_latency_ms":      avgLatencyMs,
		"error_count":         collector.ErrorCount,
		"error_rate_percent":  float64(collector.ErrorCount) / float64(collector.TotalRequests) * 100,
		"requests_by_method":  collector.RequestsByMethod,
		"requests_by_path":    collector.RequestsByPath,
		"status_codes":        collector.StatusCodes,
		"database_operations": collector.DatabaseOps,
		"memory_alloc_mb":     float64(collector.MemoryStats.Alloc) / 1024 / 1024,
		"memory_sys_mb":       float64(collector.MemoryStats.Sys) / 1024 / 1024,
		"goroutines":          runtime.NumGoroutine(),
		"last_request":        collector.LastRequestTime,
	}

	return json.MarshalIndent(metrics, "", "  ")
}

// logJSON affiche un log au format JSON
func logJSON(entry LogEntry) {
	jsonData, err := json.Marshal(entry)
	if err != nil {
		log.Printf("Erreur lors de la sérialisation du log: %v", err)
		return
	}
	log.Printf("%s", string(jsonData))
}

// getLevelString retourne la représentation string du niveau de log
func getLevelString(level LogLevel) string {
	switch level {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	default:
		return "INFO"
	}
}

// StartMetricsLogger démarre un logger périodique des métriques
func StartMetricsLogger(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			LogMetrics()
		}
	}()
}
