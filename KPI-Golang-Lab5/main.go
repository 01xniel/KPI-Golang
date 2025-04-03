package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type ReliabilityIndicatorData struct {
	FailureRate           float64
	RecoveryTime          float64
	RepairFrequency       float64
	CurrentRepairDuration *int
}

func intPtr(i int) *int {
	return &i
}

type Calculator2Params struct {
	LossesEmergency        string `json:"lossesEmergency"`
	LossesScheduled        string `json:"lossesScheduled"`
	Pm                     string `json:"pm"`
	Tm                     string `json:"tm"`
	FailureRate            string `json:"failureRate"`
	AverageRecoveryTime    string `json:"averageRecoveryTime"`
	AveragePlannedDowntime string `json:"averagePlannedDowntime"`
}

type Calculation2Result struct {
	ExpectedOutagesScheduled int `json:"expectedOutagesScheduled"`
	ExpectedOutagesEmergency int `json:"expectedOutagesEmergency"`
	ExpectedLosses           int `json:"expectedLosses"`
}

func handleCalculations1(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	reliabilityIndicators := map[string]ReliabilityIndicatorData{
		"pl110Q":           {0.007, 10.0, 0.167, intPtr(35)},
		"pl35Q":            {0.02, 8.0, 0.167, intPtr(35)},
		"pl10Q":            {0.02, 10.0, 0.167, intPtr(35)},
		"kl10TrenchQ":      {0.03, 44.0, 1.0, intPtr(9)},
		"kl10CableQ":       {0.005, 17.5, 1.0, intPtr(9)},
		"t110Q":            {0.015, 100.0, 1.0, intPtr(43)},
		"t35Q":             {0.02, 80.0, 1.0, intPtr(28)},
		"t10CableNetworkQ": {0.005, 60.0, 0.5, intPtr(10)},
		"t10AirQ":          {0.05, 60.0, 0.5, intPtr(10)},
		"v110Q":            {0.01, 30.0, 0.1, intPtr(30)},
		"v10LowOilQ":       {0.02, 15.0, 0.33, intPtr(15)},
		"v10VacuumQ":       {0.01, 15.0, 0.33, intPtr(15)},
		"busBar10Q":        {0.03, 2.0, 0.167, intPtr(5)},
		"av038Q":           {0.05, 4.0, 0.33, intPtr(10)},
		"ed610Q":           {0.1, 160.0, 0.5, nil},
		"ed038Q":           {0.1, 50.0, 0.5, nil},
	}

	numParameters := make(map[string]int)
	for key := range reliabilityIndicators {
		value := r.FormValue(key)
		intValue, err := strconv.Atoi(value)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid value for %s", key), http.StatusBadRequest)
			return
		}
		numParameters[key] = intValue
	}

	var failureRateSCS float64
	var averageRecoveryTime float64

	for key, count := range numParameters {
		if count > 0 {
			indicator := reliabilityIndicators[key]
			element := indicator.FailureRate * float64(count)
			failureRateSCS += element
			averageRecoveryTime += element * indicator.RecoveryTime
		}
	}

	averageRecoveryTime /= failureRateSCS
	coefEmergencyDowntimeSCS := (failureRateSCS * averageRecoveryTime) / 8760
	coefScheduledDowntimeSCS := 1.2 * 43 / 8760
	failureRateTCS := 2 * failureRateSCS * (coefEmergencyDowntimeSCS + coefScheduledDowntimeSCS)
	failureRateWithSectionalizerTCS := failureRateTCS + 0.02

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(map[string]float64{
		"failureRateSCS":                  failureRateSCS,
		"averageRecoveryTime":             averageRecoveryTime,
		"coefEmergencyDowntimeSCS":        coefEmergencyDowntimeSCS,
		"coefScheduledDowntimeSCS":        coefScheduledDowntimeSCS,
		"failureRateTCS":                  failureRateTCS,
		"failureRateWithSectionalizerTCS": failureRateWithSectionalizerTCS,
	})
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func handleCalculations2(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	params := Calculator2Params{
		LossesEmergency:        r.FormValue("lossesEmergency"),
		LossesScheduled:        r.FormValue("lossesScheduled"),
		Pm:                     r.FormValue("pm"),
		Tm:                     r.FormValue("tm"),
		FailureRate:            r.FormValue("failureRate"),
		AverageRecoveryTime:    r.FormValue("averageRecoveryTime"),
		AveragePlannedDowntime: r.FormValue("averagePlannedDowntime"),
	}

	pm, _ := strconv.ParseFloat(params.Pm, 64)
	tm, _ := strconv.ParseFloat(params.Tm, 64)
	failureRate, _ := strconv.ParseFloat(params.FailureRate, 64)
	avgRecoveryTime, _ := strconv.ParseFloat(params.AverageRecoveryTime, 64)
	avgPlannedDowntime, _ := strconv.ParseFloat(params.AveragePlannedDowntime, 64)
	lossesEmergency, _ := strconv.ParseFloat(params.LossesEmergency, 64)
	lossesScheduled, _ := strconv.ParseFloat(params.LossesScheduled, 64)

	commonOperand := pm * 1000 * tm

	expectedOutagesScheduled := int(avgPlannedDowntime * commonOperand)
	expectedOutagesEmergency := int(failureRate * avgRecoveryTime * commonOperand)
	expectedLosses := int(lossesEmergency*float64(expectedOutagesEmergency) + lossesScheduled*float64(expectedOutagesScheduled))

	result := Calculation2Result{
		ExpectedOutagesScheduled: expectedOutagesScheduled,
		ExpectedOutagesEmergency: expectedOutagesEmergency,
		ExpectedLosses:           expectedLosses,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})

	http.HandleFunc("/calc1", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/calc1.html")
	})
	http.HandleFunc("/evaluate1", handleCalculations1)

	http.HandleFunc("/calc2", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/calc2.html")
	})
	http.HandleFunc("/evaluate2", handleCalculations2)

	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("static"))))

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
