package api

import (
	"log"
	"net/http"
)

func Router() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/orders/", http.StatusPermanentRedirect)
	})
	mux.HandleFunc("/orders/", handleOrder)
	err := http.ListenAndServe(":3000", mux)
	log.Fatal(err)
}

func handleOrder(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		doneGet := make(chan bool, 1)
		go HandleIndex(w, doneGet)
		<-doneGet

	case http.MethodPost:
		donePost := make(chan bool, 1)
		go HandleGetOrder(w, r, donePost)
		<-donePost

	default:
		w.Header().Set("Allow", "GET, POST")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
