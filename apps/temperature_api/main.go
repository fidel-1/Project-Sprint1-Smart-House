package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

// TemperatureResponse представляет собой JSON-ответ от API температуры.
type TemperatureResponse struct {
	Value       float64   `json:"value"`
	Unit        string    `json:"unit"`
	Timestamp   time.Time `json:"timestamp"`
	Location    string    `json:"location"`
	Status      string    `json:"status"`
	SensorID    string    `json:"sensor_id"`
	SensorType  string    `json:"sensor_type"`
	Description string    `json:"description"`
}

func main() {
	// Инициализация генератора случайных чисел
	rand.Seed(time.Now().UnixNano())

	// Регистрация обработчика для всех запросов, начинающихся с /temperature/
	// Это позволяет обрабатывать /temperature?location=... и /temperature/{id}
	http.HandleFunc("/temperature/", temperatureHandler)
	http.HandleFunc("/temperature", temperatureHandler)

	log.Println("Запуск temperature API на порту 8081...")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatalf("Не удалось запустить сервер: %v", err)
	}
}

// temperatureHandler обрабатывает запросы на получение температуры.
func temperatureHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	location := r.URL.Query().Get("location")
	sensorID := r.URL.Query().Get("sensorId")

	// Проверка наличия ID сенсора в пути (например, /temperature/1)
	pathID := strings.TrimPrefix(r.URL.Path, "/temperature/")
	if pathID != "" && pathID != "temperature" {
		sensorID = pathID
	}

	// Если локация не указана, определяем её по ID сенсора
	if location == "" {
		switch sensorID {
		case "1":
			location = "Living Room"
		case "2":
			location = "Bedroom"
		case "3":
			location = "Kitchen"
		default:
			location = "Unknown"
		}
	}

	// Если ID сенсора не указан, генерируем его на основе локации
	if sensorID == "" {
		switch location {
		case "Living Room":
			sensorID = "1"
		case "Bedroom":
			sensorID = "2"
		case "Kitchen":
			sensorID = "3"
		default:
			sensorID = "0"
		}
	}

	// Генерация случайной температуры от 15.0 до 30.0
	randomTemp := 15.0 + rand.Float64()*(30.0-15.0)

	response := TemperatureResponse{
		Value:       float64(int(randomTemp*10)) / 10, // Округление до одного знака после запятой
		Unit:        "°C",
		Timestamp:   time.Now(),
		Location:    location,
		Status:      "active",
		SensorID:    sensorID,
		SensorType:  "temperature",
		Description: fmt.Sprintf("Текущая температура в %s", location),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Не удалось закодировать ответ", http.StatusInternalServerError)
	}
}
