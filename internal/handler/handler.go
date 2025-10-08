package handler

import (
	"fmt"
	"log"
	"net/mail"

	// "github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"

	d "github.com/vova1001/Website-Ylia-fitness/internal/database"
	m "github.com/vova1001/Website-Ylia-fitness/internal/model"
)

func RegisterNewUser(NewUser m.User) error {
	createTable := `
    CREATE TABLE IF NOT EXISTS users(
        id SERIAL PRIMARY KEY,
        password TEXT NOT NULL,
        email TEXT NOT NULL
    );`
	_, err := d.DB.Exec(createTable)
	if err != nil {
		return fmt.Errorf("table users not created: %w", err)
	}
	log.Println("Table Users created")

	if !EmailCheck(NewUser.Email) {
		return fmt.Errorf("invalid email: %s", NewUser.Email)
	}

	var exist bool
	err = d.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", NewUser.Email).Scan(&exist)
	if err != nil {
		return fmt.Errorf("error checking existing user: %w", err)
	}
	if exist {
		return fmt.Errorf("user with this email already exists")
	}

	hashPass, err := bcrypt.GenerateFromPassword([]byte(NewUser.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}

	_, err = d.DB.Exec("INSERT INTO users(password,email) VALUES($1,$2)", hashPass, NewUser.Email)
	if err != nil {
		return fmt.Errorf("error adding user: %w", err)
	}

	log.Println("User created successfully")
	return nil
}

func EmailCheck(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
