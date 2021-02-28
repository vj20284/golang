package main

import (
	_ "database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

func Test_GenerateToken_CreatesValidToken(t *testing.T) {
	sub := "test"
	tokenString, err := GenerateToken(sub)
	assert.Nil(t, err)

	assert.NotNil(t, tokenString)

	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		fmt.Println(err)
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		assert.Equal(t, sub, claims["sub"])
		assert.NotNil(t, claims["exp"])
	} else {
		assert.False(t, true)
	}
}

func Test_VerifyToken_ReturnsTrueForValidToken(t *testing.T) {
	username := "main_test1"
	password := "random"
	InsertUser(testdb, username, password)

	tokenString, _ := GenerateToken(username)
	UpdateToken(testdb, tokenString, username)

	assert.True(t, VerifyToken(testdb, tokenString))
}

func Test_VerifyToken_ReturnsFalseForInvalidToken(t *testing.T) {
	username := "main_test2"
	password := "random"
	InsertUser(testdb, username, password)

	tokenString, _ := GenerateToken(username)
	UpdateToken(testdb, tokenString, username)

	assert.False(t, VerifyToken(testdb, "abc"))
}

func Test_VerifyToken_ReturnsFalseForInvalidUser(t *testing.T) {
	username := "main_test3"
	tokenString, _ := GenerateToken(username)
	UpdateToken(testdb, tokenString, username)

	assert.False(t, VerifyToken(testdb, tokenString))
}

func Test_LoginHandler_ReturnsCookieWhenValidLogin(t *testing.T) {
	username := "main_test4"
	password := "random"
	db = testdb
	InsertUser(testdb, username, password)

	data := url.Values{"username": {username}, "password": {password}}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(LoginHandler)
	req, _ := http.NewRequest("POST", "/login", strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(data.Encode())))

	handler.ServeHTTP(rr, req)
	resp := rr.Result()
	authToken := ""
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "auth" {
			authToken = cookie.Value
		}
	}
	user, _ := GetUser(testdb, username)
	assert.Equal(t, user.Token, authToken)
}

func Test_LoginHandler_ReturnsForbiddenForIncorrectPassword(t *testing.T) {
	username := "main_test5"
	password := "random"
	db = testdb
	InsertUser(testdb, username, password)

	data := url.Values{"username": {username}, "password": {"blah"}}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(LoginHandler)
	req, _ := http.NewRequest("POST", "/login", strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(data.Encode())))

	handler.ServeHTTP(rr, req)
	resp := rr.Result()
	assert.Equal(t, 401, resp.StatusCode)
}

func Test_AuthHandler_ReturnsOkForValidCookie(t *testing.T) {
	username := "main_test6"
	password := "random"
	db = testdb
	InsertUser(testdb, username, password)

	data := url.Values{"username": {username}, "password": {password}}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(LoginHandler)
	req, _ := http.NewRequest("POST", "/login", strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(data.Encode())))

	handler.ServeHTTP(rr, req)
	resp := rr.Result()
	authToken := ""
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "auth" {
			authToken = cookie.Value
		}
	}
	user, _ := GetUser(testdb, username)
	assert.Equal(t, user.Token, authToken)

	handler2 := http.HandlerFunc(AuthHandler)
	authReq, _ := http.NewRequest("GET", "/auth", nil)
	cookie := http.Cookie{
		Name:  "auth",
		Value: authToken,
		Path:  "/",
	}
	authReq.AddCookie(&cookie)

	handler2.ServeHTTP(rr, authReq)
	document, _ := goquery.NewDocumentFromReader(rr.Body)
	link := ""
	document.Find("a").Each(func(_ int, s *goquery.Selection) {
		link, _ = s.Attr("href")
	})
	assert.Equal(t, "http://www.google.com", link)
}

func Test_AuthHandler_ReturnsLoginPageForInvalidCookie(t *testing.T) {
	username := "main_test7"
	password := "random"
	db = testdb
	InsertUser(testdb, username, password)

	data := url.Values{"username": {username}, "password": {password}}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(LoginHandler)
	req, _ := http.NewRequest("POST", "/login", strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(data.Encode())))

	handler.ServeHTTP(rr, req)
	resp := rr.Result()
	authToken := ""
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "auth" {
			authToken = cookie.Value
		}
	}
	user, _ := GetUser(testdb, username)
	assert.Equal(t, user.Token, authToken)

	handler2 := http.HandlerFunc(AuthHandler)
	authReq, _ := http.NewRequest("GET", "/auth", nil)
	cookie := http.Cookie{
		Name:  "auth",
		Value: "randomtoken",
		Path:  "/",
	}
	authReq.AddCookie(&cookie)

	handler2.ServeHTTP(rr, authReq)
	document, _ := goquery.NewDocumentFromReader(rr.Body)
	action := ""
	method := ""
	document.Find("form").Each(func(_ int, s *goquery.Selection) {
		action, _ = s.Attr("action")
		method, _ = s.Attr("method")
		fmt.Println("action ", action)
		fmt.Println("method ", method)
	})
	assert.Equal(t, "/login", action)
	assert.Equal(t, "POST", method)
}
