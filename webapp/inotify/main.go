package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Hello World</h1>")
}

func check(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Health check</h1>")
}

func notification(w http.ResponseWriter, r *http.Request) {
	notif_bytes, _ := json.Marshal(notifications)
	w.Write(notif_bytes)
}

func main() {
	r := mux.NewRouter()
	go watch()
	notifications = append(notifications, "Test notification")
	r.HandleFunc("/hello", index).Methods("GET")
	r.HandleFunc("/check", check).Methods("GET")
	r.HandleFunc("/notifications", notification).Methods("GET")
	fmt.Println("Server starting...")
	http.ListenAndServe(":8090", r)
}
