package main

import (
	"encoding/json"
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"
)

type ResultData struct {
	Profit                 float64
	ElectricityNoImbalance float64
	Penalty                float64
	ElectricityImbalance   float64
	NetProfit              float64
}

type CapacityRange struct {
	Min float64
	Max float64
}

func normalPowerDistributionLaw(
	p float64,
	averageDailyCapacity float64,
	standardDeviation float64,
) float64 {
	normalizationFactor := 1 / (standardDeviation * math.Sqrt(2*math.Pi))
	exponentTerm := math.Pow(p-averageDailyCapacity, 2.0) / (2 * math.Pow(standardDeviation, 2.0))
	return normalizationFactor * math.Exp(-exponentTerm)
}

func gaussLegendreIntegration(
	averageDailyCapacity float64,
	standardDeviation float64,
	capacityRange CapacityRange,
) float64 {
	t := [5]float64{-0.90617985, -0.53846931, 0.0, 0.53846931, 0.90617985}
	coefA := [5]float64{0.23692688, 0.47862868, 0.56888889, 0.47862868, 0.23692688}

	integrationResult := 0.0

	for i := 0; i < 5; i++ {
		x := 0.5*(capacityRange.Max+capacityRange.Min) + 0.5*(capacityRange.Max-capacityRange.Min)*t[i]
		y := normalPowerDistributionLaw(x, averageDailyCapacity, standardDeviation)
		coefC := 0.5 * (capacityRange.Max - capacityRange.Min) * coefA[i]
		integrationResult += coefC * y
	}

	return integrationResult
}

var tmpl *template.Template

func init() {
	var err error
	tmpl, err = template.ParseFiles("templates/index.html")
	if err != nil {
		log.Fatal(err)
	}
}

func renderTemplate(w http.ResponseWriter, data ResultData) {
	err := tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleCalculator(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, ResultData{})
}

func handleCalculations(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	averageDailyCapacityStr := r.FormValue("average_daily_capacity")
	electricityCostStr := r.FormValue("electricity_cost")
	standardDeviationStr := r.FormValue("standard_deviation")

	averageDailyCapacity, err1 := strconv.ParseFloat(averageDailyCapacityStr, 64)
	electricityCost, err2 := strconv.ParseFloat(electricityCostStr, 64)
	standardDeviation, err3 := strconv.ParseFloat(standardDeviationStr, 64)

	if err1 != nil || err2 != nil || err3 != nil {
		http.Error(w, "Invalid input data", http.StatusBadRequest)
		return
	}

	forecastError := averageDailyCapacity * 0.05
	capacityRange := CapacityRange{
		Min: averageDailyCapacity - forecastError,
		Max: averageDailyCapacity + forecastError,
	}

	electricityNoImbalanceRelativeVal := gaussLegendreIntegration(averageDailyCapacity, standardDeviation, capacityRange)
	electricityImbalanceRelativeVal := 1 - electricityNoImbalanceRelativeVal

	electricityNoImbalance := averageDailyCapacity * 24 * electricityNoImbalanceRelativeVal
	electricityImbalance := averageDailyCapacity * 24 * electricityImbalanceRelativeVal

	profit := electricityNoImbalance * electricityCost
	penalty := electricityImbalance * electricityCost
	netProfit := profit - penalty

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(ResultData{
		Profit:                 profit,
		ElectricityNoImbalance: electricityNoImbalance,
		Penalty:                penalty,
		ElectricityImbalance:   electricityImbalance,
		NetProfit:              netProfit,
	})
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/", handleCalculator)
	http.HandleFunc("/evaluate", handleCalculations)
	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("static"))))

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
