package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/bitly/go-simplejson"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var err error

type User struct {
	Id       int    `json:"id"`
	UserName string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Gender   string `json:"gender"`
}

func InitialMigration() {
	db, err := sql.Open("mysql", "nazmi:password@tcp(127.0.0.1:3306)/android")
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect database")
	}
	defer db.Close()
	// Migrate the schema
}

type Users []User

func Login(w http.ResponseWriter, r *http.Request) {

	db, err := sql.Open("mysql", "nazmi:password@tcp(127.0.0.1:3306)/android")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	vars := mux.Vars(r)
	username := vars["username"]
	password := vars["password"]

	json := simplejson.New()
	sqlStatement := `SELECT * FROM users WHERE username='` + username + `' AND password='` + password + `';`
	var user User
	row := db.QueryRow(sqlStatement)
	err2 := row.Scan(&user.Id, &username, &user.Email,
		&user.Password, &user.Gender)
	switch err2 {
	case sql.ErrNoRows:
		fmt.Println("Invalid username or password")
		json.Set("status", "failed")
		json.Set("message", "Invalid username or password")
		payload, err := json.MarshalJSON()
		if err != nil {
			log.Println(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(payload)
		return
	case nil:
		fmt.Println("Login successfull")
		json.Set("status", "success")
		json.Set("message", "Login successfull")
		payload, err := json.MarshalJSON()
		if err != nil {
			log.Println(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(payload)
		return
	default:
		fmt.Println("Error while Login  user")
		panic(err2)
	}
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "nazmi:password@tcp(127.0.0.1:3306)/android")

	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	vars := mux.Vars(r)
	username := vars["username"]
	email := vars["email"]
	password := vars["password"]
	gender := vars["gender"]
	// new query starts from here
	json := simplejson.New()

	sqlStatement := `SELECT * FROM users WHERE username='` + username + `' OR email='` + email + `';`
	var user User
	row := db.QueryRow(sqlStatement)
	err2 := row.Scan(&user.Id, &username, &user.Email,
		&user.Password, &user.Gender)
	switch err2 {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		var query string
		query = "INSERT INTO users (username,email, password, gender) VALUES('" + username + "','" + email + "','" + password + "','" + gender + "')"
		insert, err := db.Query(query)
		if err != nil {
			panic(err.Error())
		}
		defer insert.Close()
		fmt.Println("New User Successfully Created")
		json.Set("status", "success")
		json.Set("message", "New User Successfully Created")

		payload, err := json.MarshalJSON()
		if err != nil {
			log.Println(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(payload)
		return
	case nil:
		fmt.Println("User already registered")
		json.Set("status", "failed")
		json.Set("message", "User already registered")
		payload, err := json.MarshalJSON()
		if err != nil {
			log.Println(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(payload)
		return
	default:
		fmt.Println("Error while registering user")
		panic(err2)
	}
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/login/{username}/{password}", Login).Methods("POST")
	myRouter.HandleFunc("/register/{username}/{email}/{password}/{gender}", SignUp).Methods("POST")
	log.Fatal(http.ListenAndServe(":8081", myRouter))
}

func main() {
	fmt.Println("Starting server")
	InitialMigration()
	handleRequests()
}
