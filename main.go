package main

import (
	"time"
    "github.com/jinzhu/gorm"
    _ "github.com/go-sql-driver/mysql"
    "log"
	"github.com/gorilla/mux"
    "crypto/rand"
    "encoding/hex"
    "fmt"
    "log"
    "net/http"
    "net/smtp"
)

// github.com/jinzhu/gorm: GORM is an ORM (Object Relational Mapper) library for Go, which simplifies interactions with the database.
// _ "github.com/go-sql-driver/mysql": The underscore before the import path means that the package is imported solely for its side effects. In this case, the side effect is registering the MySQL driver with GORM, even though the package is not explicitly referenced in the code.

type User struct {
    ID            uint      `gorm:"primary_key"`
    Email         string    `gorm:"unique;not null"`
    PasswordHash  string    `gorm:"not null"`
    IsVerified    bool      `gorm:"default:false"`
    VerificationToken string `gorm:"not null"`
    CreatedAt     time.Time
    UpdatedAt     time.Time
}


var db *gorm.DB
var err error

// db of type *gorm.DB which will hold the database connection.

// This function is responsible for initializing the database connection.


func initDB(){
	db,err := gorm.Open("mysql", "user:password@/dbname?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	db.AutoMigrate(&User{})
}


func registerUser(w http.ResponseWriter, r *http.Request){

}

func verifyUser(w http.ResponseWriter, r *http.Request){
	
}


func main (){
	initDB()
	defer db.Close()

	r := mux.NewRouter()

	r.HandleFunc("/register",registerUser).Methods("POST")
	r.HandleFunc("/verify", verifyUser).Methods("GET")

	http.Handle("/", r)
http.ListenAndServe(":8080", nil)
}

// defer db.Close(): This ensures that the database connection is closed when the main function exits.