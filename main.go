package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type event struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type allEvents []event

var events = allEvents{
	{
		Description: "Event Description",
		ID:          1,
		Title:       "Event",
	},
}

func createEvent(w http.ResponseWriter, r *http.Request) {
	var newEvent event
	reqBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Fprintf(w, "Problems to read the requisition body")
	}

	json.Unmarshal(reqBody, &newEvent)
	events = append(events, newEvent)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newEvent)
}

func getOneEvent(w http.ResponseWriter, r *http.Request) {
	eventId := mux.Vars(r)["id"]

	eventIdInt, err := strconv.Atoi(eventId)

	if err != nil {
		fmt.Fprintf(w, "Cannot get eventId: %s", eventId)
	}

	for _, event := range events {
		if event.ID == eventIdInt {
			json.NewEncoder(w).Encode(event)
		}
	}
}

func getAllEvents(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(events)
}

func updateEvent(w http.ResponseWriter, r *http.Request) {
	eventId := mux.Vars(r)["id"]
	eventIdInt, err := strconv.Atoi(eventId)

	if err != nil {
		fmt.Fprintf(w, "Cannot get event: %s", eventId)
	}

	var updateEvent event
	reqBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Fprintf(w, "Cannot get body data")
	}

	json.Unmarshal(reqBody, &updateEvent)

	for eventIndex, event := range events {
		if event.ID == eventIdInt {
			eventToUpdate := events[eventIndex]
			eventToUpdate.Description = updateEvent.Description
			eventToUpdate.Title = updateEvent.Title

			events[eventIndex] = eventToUpdate

			json.NewEncoder(w).Encode(eventToUpdate)
		}
	}
}

func deleteEvent(w http.ResponseWriter, r *http.Request) {
	eventId := mux.Vars(r)["id"]
	eventIdInt, err := strconv.Atoi(eventId)

	if err != nil {
		fmt.Fprintf(w, "Cannot get event: %s", eventId)
	}

	for eventIndex, event := range events {
		if event.ID == eventIdInt {
			eventToDelete := event
			events = append(events[:eventIndex], events[eventIndex+1:]...)

			json.NewEncoder(w).Encode(eventToDelete)
		}
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome here!")
}

func main() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", home)
	router.HandleFunc("/event", createEvent).Methods("POST")
	router.HandleFunc("/events", getAllEvents).Methods("GET")
	router.HandleFunc("/events/{id}", getOneEvent).Methods("GET")
	router.HandleFunc("/events/{id}", updateEvent).Methods("PUT")
	router.HandleFunc("/events/{id}", deleteEvent).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", router))
}
