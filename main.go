package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
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


func generateToken() string {
    bytes := make([]byte, 16)
    rand.Read(bytes)
    return hex.EncodeToString(bytes)
}

func sendVerificationEmail(email, token string) error {
    from := "your-email@example.com"
    password := "your-email-password"

    to := []string{email}
    smtpHost := "smtp.example.com"
    smtpPort := "587"

    message := []byte(fmt.Sprintf("Subject: Email Verification\n\nPlease verify your email using this link: http://localhost:8080/verify?token=%s", token))

    auth := smtp.PlainAuth("", from, password, smtpHost)
    return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
}


func hashPassword(password string) string{
	hashedPassword,err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
        return ""
    }
    return string(hashedPassword)

}

func registerUser(w http.ResponseWriter, r *http.Request){
	email := r.FormValue("email")
	password := r.FormValue("password")

	passwordHashValue := hashPassword(password)

	token := generateToken()

	user := User{
		Email:             email,
        PasswordHash:      passwordHashValue,
        VerificationToken: token,
	}

	if err := db.Create(&user).Error;  
	err!=nil {
		http.Error(w, "Could not create user", http.StatusInternalServerError)
        return
	}

	if err := sendVerificationEmail(user.Email, token); err != nil {
        http.Error(w, "Could not send verification email", http.StatusInternalServerError)
        return
    }

    fmt.Fprintln(w, "Registration successful! Please check your email to verify your account.")
}

func verifyUser(w http.ResponseWriter, r *http.Request) {
    token := r.URL.Query().Get("token")
    if token == "" {
        http.Error(w, "Invalid token", http.StatusBadRequest)
        return
    }

    var user User
    if err := db.Where("verification_token = ?", token).First(&user).Error; err != nil {
        http.Error(w, "Invalid token", http.StatusBadRequest)
        return
    }

    user.IsVerified = true
    user.VerificationToken = ""

    if err := db.Save(&user).Error; err != nil {
        http.Error(w, "Could not verify user", http.StatusInternalServerError)
        return
    }

    fmt.Fprintln(w, "Email verified successfully!")
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