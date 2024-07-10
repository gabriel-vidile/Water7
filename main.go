package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

var (
	rainData       float64
	irrigationData []IrrigationData
	dataLock       sync.Mutex
)

// Estrutura para os dados de irrigação
type IrrigationData struct {
	DesiredWaterAmount float64 `json:"desiredWaterAmount"`
	RainData           float64 `json:"rainData"`
	SoilMoisture       float64 `json:"soilMoisture"`
	IrrigationAction   float64 `json:"irrigationAction"`
	Message            string  `json:"message"`
	SirenSound         string  `json:"sirenSound"`
}

// Função para buscar dados de chuva
func getRainDataHandler(w http.ResponseWriter, r *http.Request) {
	dataLock.Lock()
	defer dataLock.Unlock()
	response := map[string]float64{"rainData": rainData}
	json.NewEncoder(w).Encode(response)
}

// Função para guardar dados de irrigação
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
