package main

import (
	"crypto/subtle"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	host   = "localhost"
	port   = "8080"
	secret = "secret123"
)

var db *sql.DB

// LoginHandler handles login request
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		http.ServeFile(w, r, "static/login.html")
	case "POST":
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		username := r.FormValue("username")
		password := r.FormValue("password")
		dbpassword, err := GetUserPass(db, username)
		if err != nil || subtle.ConstantTimeCompare([]byte(password), []byte(dbpassword)) != 1 {
			fmt.Printf("Incorrect username or password for %s %v", username, err)
			w.WriteHeader(401)
			w.Write([]byte("You are Unauthorized to access the application.\n"))
			return
		}
		token, err := GenerateToken(username)
		if err != nil {
			fmt.Fprintf(w, "Error generating token for username %s %v", username, err)
			return
		}
		UpdateToken(db, token, username)
		addCookie(w, token, username)
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func addUserHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		http.ServeFile(w, r, "static/form.html")
	case "POST":
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		username := r.FormValue("username")
		password := r.FormValue("password")
		err := InsertUser(db, username, password)
		if err != nil {
			fmt.Fprintf(w, "Error during insertion of user %v", err.Error())
			return
		}
		fmt.Fprintf(w, "User added!\n")
	default:
		fmt.Fprintf(w, "Only GET and POST methods are supported.")
	}
}

// AuthHandler handles the auth requests
func AuthHandler(w http.ResponseWriter, r *http.Request) {
	authCookie, _ := r.Cookie("auth")
	if authCookie != nil && VerifyToken(db, authCookie.Value) {
		http.Redirect(w, r, "http://www.google.com", 301)
		return
	}
	log.Println("No Auth cookie presented")
	// redirect to login
	http.ServeFile(w, r, "static/login.html")
}

func addCookie(w http.ResponseWriter, token string, user string) {
	log.Println("Adding auth cookie for user " + user)
	expire := time.Now().Add(time.Hour * 24)
	cookie := http.Cookie{
		Name:    "auth",
		Value:   token,
		Expires: expire,
		Path:    "/",
	}
	http.SetCookie(w, &cookie)
}

// GenerateToken generates JWT token for auth
func GenerateToken(sub string) (string, error) {
	log.Printf("Generating token for %s", sub)
	claims := &jwt.MapClaims{
		"sub": sub,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	sc, err := token.SignedString([]byte(secret))
	return sc, err
}

// VerifyToken validates token with subject
func VerifyToken(db *sql.DB, tknStr string) bool {
	tkn, err := jwt.ParseWithClaims(tknStr, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil || !tkn.Valid {
		log.Printf("Error verifying token")
		return false
	}
	claims := tkn.Claims.(jwt.MapClaims)
	user, err := GetUser(db, claims["sub"].(string))
	if err != nil || subtle.ConstantTimeCompare([]byte(tknStr), []byte(user.Token)) != 1 {
		log.Printf("Token mismatch for subject: %s", claims["sub"].(string))
		return false
	}
	return true
}

func main() {
	db = InitDB()
	defer db.Close()
	http.HandleFunc("/auth", AuthHandler)
	http.HandleFunc("/adduser", addUserHandler)
	http.HandleFunc("/login", LoginHandler)
	err := http.ListenAndServe(host+":"+port, nil)
	if err != nil {
		log.Fatal("error starting http server : ", err)
		return
	}
}
