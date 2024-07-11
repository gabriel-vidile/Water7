package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

var (
	irrigationData []IrrigationData
	dataLock       sync.Mutex
)

// Estrutura para os dados de irrigação
type IrrigationData struct {
	DesiredWaterAmount float64 `json:"desiredWaterAmount"`
	RainData           float64 `json:"rainData"`
	Temperature        float64 `json:"temperature"`
	SoilMoisture       float64 `json:"soilMoisture"`
	IrrigationAction   float64 `json:"irrigationAction"`
	Message            string  `json:"message"`
	SirenSound         string  `json:"sirenSound"`
}

// Estrutura para os dados de resposta da API pública
type WeatherResponse struct {
	Temp     float64 `json:"temp"`
	Date     string  `json:"date"`
	Rain     float64 `json:"rain"`
	CityName string  `json:"city_name"`
	Forecast []struct {
		Date string  `json:"date"`
		Max  float64 `json:"max"`
		Min  float64 `json:"min"`
		Rain float64 `json:"rain"`
	} `json:"forecast"`
}

func getRainDataHandler(w http.ResponseWriter, r *http.Request) {
	apiKey := "12eef7aa"
	url := fmt.Sprintf("https://api.hgbrasil.com/weather?woeid=440202&array_limit=2&fields=only_results,temp,city_name,forecast,max,min,date,rain&key=%s", apiKey)

	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "Erro ao buscar dados de clima", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Erro ao ler resposta da API", http.StatusInternalServerError)
		return
	}

	log.Printf("Resposta da API: %s", body)

	var weatherResponse WeatherResponse
	err = json.Unmarshal(body, &weatherResponse)
	if err != nil {
		http.Error(w, "Erro ao fazer parse da resposta da API", http.StatusInternalServerError)
		return
	}

	// Supondo que estamos interessados na previsão de hoje
	rainData := weatherResponse.Forecast[0].Rain
	log.Printf("Dados de Chuva: %f", rainData)

	dataLock.Lock()
	defer dataLock.Unlock()
	response := map[string]float64{"rainData": rainData}
	json.NewEncoder(w).Encode(response)
}

func saveIrrigationDataHandler(w http.ResponseWriter, r *http.Request) {
	var data IrrigationData
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	dataLock.Lock()
	irrigationData = append(irrigationData, data)
	dataLock.Unlock()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func main() {
	http.HandleFunc("/api/rain-data", getRainDataHandler)

	http.HandleFunc("/api/irrigation-data", saveIrrigationDataHandler)

	fmt.Println("Servidor iniciado na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
