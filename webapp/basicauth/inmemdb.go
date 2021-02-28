package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
)

// InitDB initializes the sqlite DB and creates user table
func InitDB() *sql.DB {
	os.Remove("sqlite-database.db") // Cleanup existing records

	log.Println("Creating sqlite-database.db...")
	file, err := os.Create("sqlite-database.db") // Create SQLite file
	if err != nil {
		log.Fatal(err.Error())
	}
	file.Close()
	log.Println("sqlite-database.db created")

	sqliteDatabase, _ := sql.Open("sqlite3", "./sqlite-database.db") // Open the created SQLite File
	createTable(sqliteDatabase)                                      // Create Database Tables

	// Insert Dummy Record
	InsertUser(sqliteDatabase, "vivek", "password")
	return sqliteDatabase
}

func createTable(db *sql.DB) {
	createSQL := `CREATE TABLE user(
		"username" TEXT,
		"password" TEXT,
		"token" TEXT
	  );` // SQL Statement for Create Table

	log.Println("Creating user table...")
	statement, err := db.Prepare(createSQL) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec() // Execute SQL Statements
	log.Println("user table created")
}

// InsertUser will be used to create a user record in the sqlite DB
func InsertUser(db *sql.DB, username string, password string) error {
	log.Println("Inserting user record ...")
	insertUserSQL := `INSERT INTO user(username, password, token) VALUES (?, ?, ?)`
	statement, err := db.Prepare(insertUserSQL) // Prepare statement.
	if err != nil {
		log.Printf("Error during insertion of user %v", err.Error())
		return err
	}
	_, err = statement.Exec(username, password, "")
	return err
}

// UpdateToken will be used to update token for user in sqlite DB
func UpdateToken(db *sql.DB, token string, username string) error {
	log.Println("Updating token for user ", username)
	updateTokenSQL := `update user set token = ? where username = ?`
	statement, err := db.Prepare(updateTokenSQL)
	if err != nil {
		log.Printf(err.Error())
		return err
	}
	_, err = statement.Exec(token, username)
	return err
}

// GetUserPass returns the password for given username
func GetUserPass(db *sql.DB, username string) (string, error) {
	var password string
	row := db.QueryRow("select password from user where username = ?", username)
	err := row.Scan(&password)
	if err != nil {
		return "", err
	}
	return password, nil
}

// GetUser returns the user for given username
func GetUser(db *sql.DB, username string) (*User, error) {
	user := User{Username: username, Token: "", Password: ""}
	//var token string
	//var password string
	row := db.QueryRow("select password, token from user where username = ?", username)
	err := row.Scan(&user.Password, &user.Token)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// DisplayUsers will display all the user records
func DisplayUsers(db *sql.DB) {
	row, err := db.Query("select * from user order by username")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	for row.Next() { // Iterate and fetch the records from result cursor
		var username string
		var password string
		row.Scan(&username, &password)
		log.Println("User: ", username, " ", password)
	}
}
