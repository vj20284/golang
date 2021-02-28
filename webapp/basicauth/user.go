package main

import (
	_ "database/sql" // for the db tag
)

// User represents a user model
type User struct {
	Username string `db:"username"`
	Password string `db:"password"`
	Token    string `db:"token"`
}
