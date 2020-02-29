package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type GuestBookEntry struct {
	Author string `json:"author"`
	Text   string `json:"text"`
}

func getEntriesHandler(w http.ResponseWriter, r *http.Request) {
	/*
		Get the guest book entries from datastore
	*/
	entries, err := store.GetEntries()

	// Everything else is the same as before
	entriesBytes, err := json.Marshal(entries)

	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(entriesBytes)
}

func createEntryHandler(w http.ResponseWriter, r *http.Request) {
	entry := GuestBookEntry{}

	err := r.ParseForm()

	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	entry.Author = r.Form.Get("author")
	entry.Text = r.Form.Get("text")

	err = store.CreateEntry(&entry)
	if err != nil {
		fmt.Println(err)
	}

	http.Redirect(w, r, "/assets/", http.StatusFound)
}
