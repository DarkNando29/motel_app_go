package main

import (
"database/sql"
"fmt"
_ "github.com/go-sql-driver/mysql"
"log"
"time"
)

type User1 struct {
	ID   int    `json:"id"`
	FirstName string `json:"fist_name"`
	LastName string `json:"last_name"`
	EmailAddress string `json:"email_address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy string `json:"created_by"`
	UpdatedBy string `json:"updated_by"`
}

func main() {
	fmt.Println("Go MySQL Tutorial")

	// Open up our database connection.
	// I've set up a database on my local machine using phpmyadmin.
	// The database is called testDb
	db, err := sql.Open("mysql", "appMotelDev:12345@tcp(localhost:3306)/motel")

	// if there is an error opening the connection, handle it
	if err != nil {
		panic(err.Error())
	}

	// defer the close till after the main function has finished
	// executing
	defer db.Close()

	results, err := db.Query("SELECT id, first_name FROM users")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	for results.Next() {
		var user User
		// for each row, scan the result into our tag composite object
		err = results.Scan(&user.ID, &user.FirstName)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		// and then print out the tag's Name attribute
		log.Printf(user.FirstName)
	}
}