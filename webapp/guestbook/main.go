package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// The new router function creates the router and
// returns it to us. We can now use this function
// to instantiate and test the router outside of the main function
func newRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/hello", handler).Methods("GET")

	// Declare the static file directory and point it to the directory we just made
	staticFileDirectory := http.Dir("./assets/")
	// Declare the handler, that routes requests to their respective filename.
	// The fileserver is wrapped in the `stripPrefix` method, because we want to
	// remove the "/assets/" prefix when looking for files.
	// For example, if we type "/assets/index.html" in our browser, the file server
	// will look for only "index.html" inside the directory declared above.
	// If we did not strip the prefix, the file server would look for "./assets/assets/index.html", and yield an error
	staticFileHandler := http.StripPrefix("/assets/", http.FileServer(staticFileDirectory))
	// The "PathPrefix" method acts as a matcher, and matches all routes starting
	// with "/assets/", instead of the absolute route itself
	r.PathPrefix("/assets/").Handler(staticFileHandler).Methods("GET")

	r.HandleFunc("/gbe", getEntriesHandler).Methods("GET")
	r.HandleFunc("/gbe", createEntryHandler).Methods("POST")
	return r
}

func getDbConnection() *sql.DB {
	for i := 1; i <= 5; i++ {
		connString := "host=postgres port=5432 user=postgres password=example dbname=temp sslmode=disable"
		_db, _err := sql.Open("postgres", connString)
		_err = _db.Ping()
		if _err == nil {
			return _db
		}
		fmt.Println("Waiting for connection")
		time.Sleep(3 * time.Second)
	}
	return nil
}

func main() {
	fmt.Println("Starting server...")

	db := getDbConnection()

	m, err := migrate.New(
		"file://migrations",
		"postgres://postgres:example@postgres:5432/temp?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	m.Up()

	InitStore(&dbStore{db: db})

	// The router is now formed by calling the `newRouter` constructor function
	// that we defined above. The rest of the code stays the same
	r := newRouter()
	fmt.Println("Serving on port 8090")
	http.ListenAndServe(":8090", r)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}
