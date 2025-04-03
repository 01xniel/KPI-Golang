package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type Calculation1Result struct {
	Coefs                  map[string]float64 `json:"coefs"`
	CompositionDry         map[string]float64 `json:"compositionDry"`
	CompositionCombustible map[string]float64 `json:"compositionCombustible"`
	LowHeatingValues       map[string]float64 `json:"lowHeatingValues"`
}

type Calculation2Result struct {
	RawComposition map[string]float64 `json:"rawComposition"`
	RawLHV         float64            `json:"rawLHV"`
}

func getFuelComposition(components map[string]float64, fuelType string, coefficient float64) map[string]float64 {
	composition := make(map[string]float64)
	for key, value := range components {
		if key == "moisture" || (key == "ash" && fuelType == "combustible") {
			continue
		}
		composition[key] = value * coefficient
	}
	return composition
}

func getLHV(components map[string]float64) map[string]float64 {
	rawLHV := (339*components["carbon"] + 1030*components["hydrogen"] - 108.8*(components["oxygen"]-components["sulfur"]) - 25*components["moisture"]) / 1000
	dryLHV := (rawLHV + 0.025*components["moisture"]) * 100 / (100 - components["moisture"])
	combustibleLHV := (rawLHV + 0.025*components["moisture"]) * 100 / (100 - components["moisture"] - components["ash"])

	return map[string]float64{
		"raw":         rawLHV,
		"dry":         dryLHV,
		"combustible": combustibleLHV,
	}
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

	hydrogen, _ := strconv.ParseFloat(r.FormValue("hydrogen"), 64)
	carbon, _ := strconv.ParseFloat(r.FormValue("carbon"), 64)
	sulfur, _ := strconv.ParseFloat(r.FormValue("sulfur"), 64)
	nitrogen, _ := strconv.ParseFloat(r.FormValue("nitrogen"), 64)
	oxygen, _ := strconv.ParseFloat(r.FormValue("oxygen"), 64)
	moisture, _ := strconv.ParseFloat(r.FormValue("moisture"), 64)
	ash, _ := strconv.ParseFloat(r.FormValue("ash"), 64)

	coefs := map[string]float64{
		"dry":         100 / (100 - moisture),
		"combustible": 100 / (100 - moisture - ash),
	}

	components := map[string]float64{
		"hydrogen": hydrogen,
		"carbon":   carbon,
		"sulfur":   sulfur,
		"nitrogen": nitrogen,
		"oxygen":   oxygen,
		"moisture": moisture,
		"ash":      ash,
	}

	compositionDry := getFuelComposition(components, "dry", coefs["dry"])
	compositionCombustible := getFuelComposition(components, "combustible", coefs["combustible"])
	lowHeatingValues := getLHV(components)

	result := Calculation1Result{
		Coefs:                  coefs,
		CompositionDry:         compositionDry,
		CompositionCombustible: compositionCombustible,
		LowHeatingValues:       lowHeatingValues,
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(result)
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

	carbon, _ := strconv.ParseFloat(r.FormValue("carbon"), 64)
	hydrogen, _ := strconv.ParseFloat(r.FormValue("hydrogen"), 64)
	oxygen, _ := strconv.ParseFloat(r.FormValue("oxygen"), 64)
	sulfur, _ := strconv.ParseFloat(r.FormValue("sulfur"), 64)
	combustibleLHV, _ := strconv.ParseFloat(r.FormValue("combustibleLHV"), 64)
	rawMoisture, _ := strconv.ParseFloat(r.FormValue("rawMoisture"), 64)
	dryAsh, _ := strconv.ParseFloat(r.FormValue("dryAsh"), 64)
	combustibleVanadium, _ := strconv.ParseFloat(r.FormValue("combustibleVanadium"), 64)

	commonOperand := 100 - rawMoisture - dryAsh

	rawComposition := map[string]float64{
		"carbon":   carbon * commonOperand / 100,
		"hydrogen": hydrogen * commonOperand / 100,
		"oxygen":   oxygen * commonOperand / 100,
		"sulfur":   sulfur * commonOperand / 100,
		"moisture": rawMoisture,
		"ash":      dryAsh * (100 - rawMoisture) / 100,
		"vanadium": combustibleVanadium * (100 - rawMoisture) / 100,
	}

	rawLHV := combustibleLHV*commonOperand/100 - 0.025*rawMoisture

	result := Calculation2Result{
		RawComposition: rawComposition,
		RawLHV:         rawLHV,
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(result)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
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
