package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Обработчик для корневого пути
func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

// Обработчик для проверки доступности (Healthcheck)
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK")
}

func main() {
	log.Printf("INFO: Starting web application...")

	// Параметризация: порт из переменной окружения
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", helloHandler)
	http.HandleFunc("/health", healthHandler)
	http.Handle("/metrics", promhttp.Handler())

	log.Printf("INFO: Listening on :%s", port)

	// Логирование: Go-сервер автоматически логирует запросы на stdout/stderr,
	// которые Docker собирает и отправляет в централизованный лог.
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("FATAL: Server failed to start: %v", err)
	}
}
