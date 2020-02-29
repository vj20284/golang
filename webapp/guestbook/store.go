package main

// The sql go library is needed to interact with the database
import (
	"database/sql"
)

// Our store will have two methods, to add a new Entry,
// and to get all existing entries
// Each method returns an error, in case something goes wrong
type Store interface {
	CreateEntry(entry *GuestBookEntry) error
	GetEntries() ([]*GuestBookEntry, error)
}

// The `dbStore` struct will implement the `Store` interface
// It also takes the sql DB connection object, which represents
// the database connection.
type dbStore struct {
	db *sql.DB
}

func (store *dbStore) CreateEntry(entry *GuestBookEntry) error {
	// 'Entry' is a simple struct which has "author" and "text" attributes
	// THe first underscore means that we don't care about what's returned from
	// this insert query. We just want to know if it was inserted correctly,
	// and the error will be populated if it wasn't
	_, err := store.db.Query("INSERT INTO entries(author, text) VALUES ($1,$2)", entry.Author, entry.Text)
	return err
}

func (store *dbStore) GetEntries() ([]*GuestBookEntry, error) {
	// Query the database for all entries, and return the result to the
	// `rows` object
	rows, err := store.db.Query("SELECT author, text from entries")
	// We return incase of an error, and defer the closing of the row structure
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Create the data structure that is returned from the function.
	// By default, this will be an empty array of entries
	entries := []*GuestBookEntry{}
	for rows.Next() {
		// For each row returned by the table, create a pointer to a entry,
		entry := &GuestBookEntry{}
		// Populate the `Author` and `Text` attributes of the entry,
		// and return incase of an error
		if err := rows.Scan(&entry.Author, &entry.Text); err != nil {
			return nil, err
		}
		// Finally, append the result to the returned array, and repeat for
		// the next row
		entries = append(entries, entry)
	}
	return entries, nil
}

// The store variable is a package level variable that will be available for
// use throughout our application code
var store Store

/*
We will need to call the InitStore method to initialize the store. This will
typically be done at the beginning of our application (in this case, when the server starts up)
This can also be used to set up the store as a mock, which we will be observing
later on
*/
func InitStore(s Store) {
	store = s
}
