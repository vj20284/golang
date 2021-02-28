package main

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testdb *sql.DB

func init() {
	testdb = InitDB()
}

func Test_InsertUser_GetRetrievesUser(t *testing.T) {
	username := "test"
	password := "random"
	err := InsertUser(testdb, username, password)
	assert.Nil(t, err)

	user, err := GetUser(testdb, username)
	assert.Nil(t, err)
	assert.Equal(t, username, user.Username)
	assert.Equal(t, password, user.Password)
	assert.Equal(t, "", user.Token)
}

func Test_InsertUser_GetUserPassRetrievesPassword(t *testing.T) {
	username := "test2"
	password := "random2"
	err := InsertUser(testdb, username, password)
	assert.Nil(t, err)

	dbpassword, err := GetUserPass(testdb, username)
	assert.Nil(t, err)
	assert.Equal(t, password, dbpassword)
}

func Test_UpdateToken_VerifyTokenUpdated(t *testing.T) {
	username := "test3"
	password := "random3"
	err := InsertUser(testdb, username, password)
	assert.Nil(t, err)

	token := "mytoken"
	err = UpdateToken(testdb, token, username)
	assert.Nil(t, err)
	user, err := GetUser(testdb, username)
	assert.Nil(t, err)
	assert.Equal(t, username, user.Username)
	assert.Equal(t, password, user.Password)
	assert.Equal(t, token, user.Token)
}
