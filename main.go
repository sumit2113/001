package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	_ "github.com/sumit2113/001/docs"
	"github.com/gorilla/mux"
	"github.com/swaggo/http-swagger"
)

// Alarm represents an alarm object
type Alarm struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	StartTime time.Time `json:"startTime"`
	Severity  string    `json:"severity"`
}

var alarms []Alarm

// @title Alarm API
// @version 1.0
// @description This is a sample Alarm API
// @host localhost:8080
// @BasePath /
func main() {
	alarms = append(alarms, Alarm{ID: rand.Intn(9999), Message: "Write the test cases for the sample code", Title: "Testing", StartTime: time.Now(), Severity: "Major"})
	alarms = append(alarms, Alarm{ID: rand.Intn(9999), Title: "Testing", Message: "Testing successfully completed", StartTime: time.Now(), Severity: "Minor"})
	alarms = append(alarms, Alarm{ID: rand.Intn(9999), Title: "Execution", Message: "Execute the above code", StartTime: time.Now(), Severity: "Major"})

	handleRequests()
}

// @Summary Generate an alarm
// @Description Generate a new alarm
// @Accept json
// @Produce json
// @Param alarm body Alarm true "Alarm object to be generated"
// @Success 201 {object} Alarm
// @Router /generateAlarms [post]
func createAlarm(w http.ResponseWriter, r *http.Request) {
	var newAlarm Alarm
	_ = json.NewDecoder(r.Body).Decode(&newAlarm)

	newAlarm.ID = rand.Intn(9999)
	newAlarm.StartTime = time.Now()

	alarms = append(alarms, newAlarm)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newAlarm)
	fmt.Fprintf(w, "Alarm with ID %d and title %s is generated\n", newAlarm.ID, newAlarm.Title)
}

// @Summary Get all alarms
// @Description Get all alarms
// @Produce json
// @Success 200 {array} Alarm
// @Router /alarms [get]
func getAllAlarms(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(alarms)
}

// @Summary Clear an alarm
// @Description Clear an alarm by ID
// @Produce json
// @Param id path int true "Alarm ID to be deleted"
// @Success 200 {string} string
// @Router /clearAlarms/{id} [delete]
func deleteAlarm(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	for index, alarm := range alarms {
		if alarm.ID == id {
			alarms = append(alarms[:index], alarms[index+1:]...)
			fmt.Fprintf(w, "Alarm with title %s and ID %d has been deleted\n", alarm.Title, alarm.ID)
			break
		}
		w.WriteHeader(http.StatusOK)
	}
	w.WriteHeader(http.StatusNotFound)
}

// @Summary Update an alarm
// @Description Update an existing alarm by ID
// @Accept json
// @Produce json
// @Param id path int true "Alarm ID to be updated"
// @Param alarm body Alarm true "Updated alarm object"
// @Success 200 {object} Alarm
// @Router /updateAlarms/{id} [put]
func updateAlarm(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	var updateAlarm Alarm
	_ = json.NewDecoder(r.Body).Decode(&updateAlarm)

	for index, alarm := range alarms {
		if alarm.ID == id {
			alarm.Title = updateAlarm.Title
			alarm.Message = updateAlarm.Message
			alarm.Severity = updateAlarm.Severity

			alarms[index] = alarm
			json.NewEncoder(w).Encode(alarm)
			fmt.Fprintf(w, "Alarm with title %s is updated\n", alarm.Title)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "Alarm with ID %d is not found", id)
}

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Alarm server.
// @BasePath /swagger
func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/generateAlarms", createAlarm).Methods("POST")
	router.HandleFunc("/alarms", getAllAlarms).Methods("GET")
	router.HandleFunc("/clearAlarms/{id}", deleteAlarm).Methods("DELETE")
	router.HandleFunc("/updateAlarms/{id}", updateAlarm).Methods("PUT")

	// Swagger UI endpoint
	//router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)
	//router.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger/", http.FileServer(http.Dir("./docs"))))

	/*router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/docs/swagger.json"), //The url pointing to API definition
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet)*/

	log.Println("Starting API Gateway on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
