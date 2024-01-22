package api

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/vojdelenie/task0/internal/db"
)

func HandleIndex(w http.ResponseWriter, done chan bool) {
	templatePath := "frontend/index.html"
	temp, error := template.New("index.html").ParseFiles(templatePath)
	if error != nil {
		log.Print(error)
	}
	temp.Execute(w, nil)
	done <- true
}

func HandleGetOrder(w http.ResponseWriter, r *http.Request, doneGetOrder chan bool) {
	orderRequest := struct {
		ID string `json:"id"`
	}{}
	json.NewDecoder(r.Body).Decode(&orderRequest)
	response := string(handleResponse(orderRequest.ID))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	doneGetOrder <- true
}

func handleResponse(data string) []byte {
	order := db.Cache.Get(data)
	if order != nil {
		orderJSON, err := json.Marshal(order)
		if err != nil {
			log.Printf("Error marshalling order: %v", err)
			return []byte("An error occurred.")
		}
		log.Printf("Got order %v from CACHE", data)
		return orderJSON
	}

	id, err := uuid.Parse(data)
	if err != nil {
		log.Printf("Error parsing order id: %v", err)
		return []byte("Invalid ID")
	}
	order, err = db.Instance.GetOrder(id)
	if err != nil {
		log.Printf("Error getting order: %v", err)
		log.Printf("data: %v", data)
		return []byte("Order not found")
	}
	db.Cache.Set(data, order)
	order = db.Cache.Get(data)
	orderJSON, err := json.Marshal(order)
	if err != nil {
		log.Printf("Error marshalling order: %v", err)
		return []byte("An error occurred.")
	}
	return orderJSON
}
