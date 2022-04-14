package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"time"
)

type User struct {
	ID   int    `json:"id"`
	FirstName string `json:"fist_name"`
	LastName string `json:"last_name"`
	EmailAddress string `json:"email_address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy string `json:"created_by"`
	UpdatedBy string `json:"updated_by"`
}

type JsonResponse struct {
	Type    string `json:"type"`
	Data    []User `json:"data"`
	Message string `json:"message"`
}

func setupDB() *sql.DB {
	db, err := sql.Open("mysql", "appMotelDev:12345@tcp(localhost:3306)/motel?parseTime=true")

	checkErr(err)

	return db
}

func printMessage(message string) {
	fmt.Println("")
	fmt.Println(message)
	fmt.Println("")
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

// Get all movies

// GetUsers response and request handlers
func GetUsers(w http.ResponseWriter, r *http.Request) {
	db := setupDB()

	printMessage("Getting movies...")

	// Get all movies from movies table that don't have movieID = "1"
	rows, err := db.Query("SELECT * FROM users")

	// check errors
	checkErr(err)

	// var response []JsonResponse
	var users []User

	// Foreach movie
	for rows.Next() {
		var id int
		var firstName string
		var lastName string
		var emailAddress string
		var createdAt time.Time
		var updatedAt time.Time
		var createdBy string
		var updatedBy string

		err = rows.Scan(&id, &createdAt, &createdBy, &emailAddress, &firstName, &lastName,   &updatedAt,  &updatedBy)

		// check errors
		checkErr(err)

		users = append(users, User{
			ID: id,
			FirstName: firstName,
			LastName: lastName,
			EmailAddress: emailAddress,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
			CreatedBy: createdBy,
			UpdatedBy: updatedBy,
		})
	}

	var response = JsonResponse{Type: "success", Data: users}

	json.NewEncoder(w).Encode(response)
}

// Create a movie

// CreateUser response and request handlers
func CreateUser(w http.ResponseWriter, r *http.Request) {
	userId := r.FormValue("userId")
	userFirstName := r.FormValue("userFirstName")

	var response = JsonResponse{}

	if userId == "" || userFirstName == "" {
		response = JsonResponse{Type: "error", Message: "You are missing userId or userFirstName parameter."}
	} else {
		db := setupDB()

		printMessage("Inserting movie into DB")

		fmt.Println("Inserting new user with ID: " + userId + " and name: " + userFirstName)

		var lastInsertID int
		err := db.QueryRow("INSERT INTO user(id, created_at,created_by,email_address,first_name,last_name,updated_at,updated_by) VALUES($1, $2) returning id;", userId, userFirstName).Scan(&lastInsertID)

		// check errors
		checkErr(err)

		response = JsonResponse{Type: "success", Message: "The user has been inserted successfully!"}
	}

	json.NewEncoder(w).Encode(response)
}

// DeleteUser response and request handlers
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	userId := params["userId"]

	var response = JsonResponse{}

	if userId == "" {
		response = JsonResponse{Type: "error", Message: "You are missing movieID parameter."}
	} else {
		db := setupDB()

		printMessage("Deleting user from DB")

		_, err := db.Exec("DELETE FROM users where id = $1", userId)

		// check errors
		checkErr(err)

		response = JsonResponse{Type: "success", Message: "The user has been deleted successfully!"}
	}

	json.NewEncoder(w).Encode(response)
}

// Delete all users

// DeleteUsers response and request handlers
func DeleteUsers(w http.ResponseWriter, r *http.Request) {
	db := setupDB()

	printMessage("Deleting all users...")

	_, err := db.Exec("DELETE FROM users")

	// check errors
	checkErr(err)

	printMessage("All users have been deleted successfully!")

	var response = JsonResponse{Type: "success", Message: "All users have been deleted successfully!"}

	json.NewEncoder(w).Encode(response)
}


func main() {
	// Init the mux router
	router := mux.NewRouter()

	// Route handles & endpoints

	// Get all users
	router.HandleFunc("/users/", GetUsers).Methods("GET")

	// Create a users
	router.HandleFunc("/users/", CreateUser).Methods("POST")

	// Delete a specific users by the userid
	router.HandleFunc("/users/{userid}", DeleteUsers).Methods("DELETE")

	// Delete all users
	router.HandleFunc("/users/", DeleteUsers).Methods("DELETE")

	// serve the app
	fmt.Println("Server at 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}