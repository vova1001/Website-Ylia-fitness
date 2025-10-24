package handler

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"

	d "github.com/vova1001/Website-Ylia-fitness/internal/database"
	m "github.com/vova1001/Website-Ylia-fitness/internal/model"
	o "github.com/vova1001/Website-Ylia-fitness/internal/otherFunc"
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

	if !o.EmailCheck(NewUser.Email) {
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

func AuthUser(User m.User) (m.Token, error) {
	var exist bool
	var UserPass string
	if !o.EmailCheck(User.Email) {
		return m.Token{}, fmt.Errorf("invalid email: %s", User.Email)
	}
	err := d.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", User.Email).Scan(&exist)
	if err != nil {
		return m.Token{}, fmt.Errorf("error checking existing user: %w", err)
	}
	if !exist {
		return m.Token{}, fmt.Errorf("there is no user with this email, please register")
	}
	err = d.DB.QueryRow("SELECT password FROM users WHERE email =$1", User.Email).Scan(&UserPass)
	if err != nil {
		return m.Token{}, fmt.Errorf("error checking password user: %w", err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(UserPass), []byte(User.Password))
	if err != nil {
		return m.Token{}, fmt.Errorf("error comparison password")
	}
	sekretKey := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":   User.ID,
		"email":    User.Email,
		"timeLife": time.Now().Add(100 * time.Hour).Unix(),
	})
	signedToken, err := token.SignedString([]byte(sekretKey))
	if err != nil {
		return m.Token{}, fmt.Errorf("error signed token")
	}
	var SignedToken m.Token
	SignedToken.JWT_Token = signedToken

	return SignedToken, nil

}

func FogotPass(email m.FogotPass) (string, error) {
	var exist bool
	var tokenNP m.TokenNewPass
	err := d.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", email.Email).Scan(&exist)
	if err != nil {
		return "", fmt.Errorf("error checking existing user: %w", err)
	}
	if !exist {
		return "", fmt.Errorf("there is no user with this email, please register")
	}
	createTable := `
	CREATE TABLE IF NOT EXISTS password_resets(
		id SERIAL PRIMARY KEY,
		email TEXT NOT NULL,
		token_hash TEXT NOT NULL,
		time_life TIMESTAMP NOT NULL,
		used BOOLEAN DEFAULT FALSE
	);
	`
	_, err = d.DB.Exec(createTable)
	if err != nil {
		return "", fmt.Errorf("table password_resets not created: %w", err)
	}
	token, err := o.GeneratorToken(32)
	if err != nil {
		return "", err
	}
	log.Printf("DEBUG reset token for %s = %s\n", email.Email, token)
	_, _ = d.DB.Exec("DELETE FROM password_resets WHERE email=$1 AND used=FALSE", email)
	hash := sha256.Sum256([]byte(token))
	tokenNP.HashToken = hex.EncodeToString(hash[:])
	tokenNP.TimeLife = time.Now().Add(time.Minute * 15)
	tokenNP.Used = false
	_, err = d.DB.Exec("INSERT INTO password_resets(email, token_hash, time_life) VALUES($1,$2,$3)", email.Email, tokenNP.HashToken, tokenNP.TimeLife)
	if err != nil {
		return "", fmt.Errorf("error adding token info: %w", err)
	}

	resetLink := fmt.Sprintf("https://yourfrontend.com/reset-password?token=%s", token)

	err = o.SendResetEmail(email.Email, resetLink)
	if err != nil {
		return "", fmt.Errorf("error sending email: %w", err)
	}
	return token, nil
}

func ResetPassword(NewPass m.NewPass) error {
	var tokenNP m.TokenNewPass
	hash := sha256.Sum256([]byte(NewPass.Token))
	tokenHash := hex.EncodeToString(hash[:])
	err := d.DB.QueryRow("SELECT email, time_life, used FROM password_resets WHERE token_hash = $1", tokenHash).Scan(&tokenNP.EmailToToken, &tokenNP.TimeLife, &tokenNP.Used)
	if err != nil {
		return fmt.Errorf("err scan from pass_resets")
	}
	if tokenNP.Used || (time.Now().After(tokenNP.TimeLife)) {
		return fmt.Errorf("token invalid")
	}
	NewHashedPass, _ := bcrypt.GenerateFromPassword([]byte(tokenHash), bcrypt.DefaultCost)
	_, err = d.DB.Exec("UPDATE users SET password=$1 WHERE email=$2", NewHashedPass, tokenNP.EmailToToken)
	if err != nil {
		return fmt.Errorf("err update new password")
	}
	tokenNP.Used = true
	_, err = d.DB.Exec("UPDATE password_resets SET used=$1 WHERE token_hash=$2", tokenNP.Used, tokenNP.HashToken)
	if err != nil {
		return fmt.Errorf("err update table pas_resets after apdate new pass")
	}
	return nil
}

func PurchesRequest(PR m.PurchaseRequest, UserID int, Email string) error {
	var Purchase m.Purchase
	Purchase.UserID = UserID
	Purchase.ProductID = PR.IdProduct
	Purchase.Email = Email
	Purchase.CreateadAt = time.Now()
	Purchase.PaymentID = ""
	err := d.DB.QueryRow("SELECT product_name, product_price FROM products WHERE id=$1", PR.IdProduct).Scan(&Purchase.ProductName, &Purchase.ProductPrice)
	if err != nil {
		return fmt.Errorf("err scan from product")
	}
	yk := o.NewYookassaClient(
		o.GetEnv("YOOKASSA_SHOP_ID", ""),
		o.GetEnv("YOOKASSA_API_KEY", ""),
	)
	return nil
}
