package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"
)

// Параметризация: Используем переменные окружения для целевого хоста и порта
const (
	MonitorLogFile = "monitor.log"
)

func init() {
	// Настройка логирования для записи в файл
	file, err := os.OpenFile(MonitorLogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("FATAL: Failed to open monitor log file: %v", err)
	}
	log.SetOutput(file)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {
	// Параметризация: Получаем хост/порт приложения из переменных окружения
	appHost := os.Getenv("APP_TARGET_HOST")
	appPort := os.Getenv("APP_TARGET_PORT")
	if appHost == "" || appPort == "" {
		log.Fatal("FATAL: APP_TARGET_HOST and APP_TARGET_PORT must be set.")
	}

	targetURL := fmt.Sprintf("http://%s:%s/health", appHost, appPort)
	log.Printf("INFO: Starting health check loop for %s", targetURL)

	// Скрипт мониторинга запускается каждые N секунд
	intervalStr := os.Getenv("CHECK_INTERVAL")
	if intervalStr == "" {
		intervalStr = "5s"
	}
	interval, err := time.ParseDuration(intervalStr)
	if err != nil {
		log.Fatalf("FATAL: Invalid CHECK_INTERVAL: %v", err)
	}

	for {
		// Проверка доступности приложения
		resp, err := http.Get(targetURL)

		if err != nil || resp.StatusCode != http.StatusOK {
			// Логирование результатов проверки (сбой)
			log.Printf("ERROR: Health check failed for %s. Attempting restart. Error: %v, Status: %d",
				targetURL, err, func() int {
					if resp != nil {
						return resp.StatusCode
					}
					return 0
				}())

			log.Println("INFO: Executing restart script for web-app...")

			// Вызов Bash-скрипта, который выполнит docker compose restart
			cmd := exec.Command("/restart_app.sh")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			if err := cmd.Run(); err != nil {
				log.Printf("FATAL: Restart script failed: %v", err)
			} else {
				log.Println("SUCCESS: Restart script executed. Waiting for service to come up.")
			}

		} else {
			// Логирование результатов проверки (успех)
			log.Printf("INFO: Health check succeeded. Status: %d", resp.StatusCode)
		}

		if resp != nil {
			resp.Body.Close()
		}

		time.Sleep(interval)
	}
}
