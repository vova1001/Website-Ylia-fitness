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
	var UserID int
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
	err = d.DB.QueryRow("SELECT id, password FROM users WHERE email =$1", User.Email).Scan(&UserID, &UserPass)
	if err != nil {
		return m.Token{}, fmt.Errorf("error checking password user: %w", err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(UserPass), []byte(User.Password))
	if err != nil {
		return m.Token{}, fmt.Errorf("error comparison password")
	}
	sekretKey := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":   UserID,
		"email":    User.Email,
		"timeLife": time.Now().Add(100 * time.Hour).Unix(),
	})
	signedToken, err := token.SignedString([]byte(sekretKey))
	if err != nil {
		return m.Token{}, fmt.Errorf("error signed token %w", err)
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
		return fmt.Errorf("err scan from pass_resets %w", err)
	}
	if tokenNP.Used || (time.Now().After(tokenNP.TimeLife)) {
		return fmt.Errorf("token invalid")
	}
	NewHashedPass, _ := bcrypt.GenerateFromPassword([]byte(tokenHash), bcrypt.DefaultCost)
	_, err = d.DB.Exec("UPDATE users SET password=$1 WHERE email=$2", NewHashedPass, tokenNP.EmailToToken)
	if err != nil {
		return fmt.Errorf("err update new password %w", err)
	}
	tokenNP.Used = true
	_, err = d.DB.Exec("UPDATE password_resets SET used=$1 WHERE token_hash=$2", tokenNP.Used, tokenNP.HashToken)
	if err != nil {
		return fmt.Errorf("err update table pas_resets after apdate new pass")
	}
	return nil
}

func ProductAddBasket(UserID, ProductID int, Email string) (string, error) {
	var Basket m.Basket
	Basket.UserID = UserID
	Basket.ProductID = ProductID
	Basket.Email = Email
	err := d.DB.QueryRow("SELECT product_name, product_price FROM products WHERE id=$1", Basket.ProductID).Scan(&Basket.ProductName, &Basket.ProductPrice)
	if err != nil {
		return "", fmt.Errorf("err scan from product %w", err)
	}

	_, err = d.DB.Exec(`
		INSERT INTO basket 
		(user_id, email, product_id, product_name, product_price )
		VALUES ($1, $2, $3, $4, $5)`,
		Basket.UserID,
		Basket.Email,
		Basket.ProductID,
		Basket.ProductName,
		Basket.ProductPrice,
	)
	if err != nil {
		return "", fmt.Errorf("err insert basket: %w", err)
	}
	return "Successfully", nil
}

func PurchesRequest(UserId int) (string, error) {
	var PR m.PurchaseRequest
	var items []m.PurchaseItem
	rows, _ := d.DB.Query("SELECT email, product_id, product_name, product_price FROM basket WHERE user_id=$1", UserId)
	for rows.Next() {
		var item m.PurchaseItem
		//Пусть перезаписывает почту, я так вижу :)
		err := rows.Scan(&PR.Email, &item.ProductID, &item.ProductName, &item.ProductPrice)
		if err != nil {
			return "", fmt.Errorf("err scan from basket %w", err)
		}
		PR.TotalAmount += item.ProductPrice
		items = append(items, item)
	}
	PR.CreateadAt = time.Now()
	PR.UserID = UserId
	PR.PaymentID = ""
	yc := o.NewYookassaClient(
		o.GetEnv("YOOKASSA_SHOP_ID", ""),
		o.GetEnv("YOOKASSA_API_KEY", ""),
	)

	var NamesItemsFromYK string
	if len(items) > 1 {
		NamesItemsFromYK = "Оплата товаров"
	} else {
		NamesItemsFromYK = "Оплата товара"
	}
	for _, item := range items {
		NamesItemsFromYK += fmt.Sprintf(",%s", item.ProductName)
	}

	resp, err := o.CreatePayment(yc, PR.TotalAmount, NamesItemsFromYK)
	if err != nil {
		return "", fmt.Errorf("err create payment %w", err)
	}
	PR.PaymentID = resp.ID

	var PurchaseRequestsID int
	err = d.DB.QueryRow(`
		INSERT INTO purchase_requests
		(user_id, email, total_amount, payment_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`,
		PR.UserID,
		PR.Email,
		PR.TotalAmount,
		PR.PaymentID,
		PR.CreateadAt,
	).Scan(&PurchaseRequestsID)
	if err != nil {
		return "", fmt.Errorf("err returned id purchase_request: %w", err)
	}

	for _, item := range items {
		_, err := d.DB.Exec(`
			INSERT INTO purchase_items
			(purchase_request_id, product_id, product_name, product_price)
			VALUES ($1, $2, $3, $4)
			`, PurchaseRequestsID, item.ProductID, item.ProductName, item.ProductPrice)
		if err != nil {
			return "", fmt.Errorf("err insert purchase_item: %w", err)
		}
	}

	return resp.Confirmation.ConfirmationURL, nil
}

func WebhookY(Webook m.YookassaWebhook) error {
	var PurchasePaid m.PurchasePaid
	if Webook.Event == "payment.succeeded" && Webook.Object.Status == "succeeded" && Webook.Object.Paid {
		PurchasePaid.PaymentID = Webook.Object.ID
	} else {
		return fmt.Errorf("payment failed")
	}

	err := d.DB.QueryRow(`
    SELECT user_id, email, payment_id, id
    FROM purchase_requests 
    WHERE payment_id=$1`, PurchasePaid.PaymentID).Scan(&PurchasePaid.UserID, &PurchasePaid.Email, &PurchasePaid.PaymentID, &PurchasePaid.ID)
	if err != nil {
		return fmt.Errorf("err scan from purchase_requests")
	}
	PurchasePaid.SubStart = time.Now()
	PurchasePaid.SubEnd = PurchasePaid.SubStart.Add(720 * time.Hour)

	rows, _ := d.DB.Query(`
		SELECT product_id, product_name, product_price 
		FROM purchase_items
		WHERE purchase_request_id=$1
	`, PurchasePaid.ID)

	defer rows.Close()

	var PurchasePaidItems []m.PurchaseItem
	for rows.Next() {
		var PurchasePaidItem m.PurchaseItem
		err := rows.Scan(&PurchasePaidItem.ProductID, &PurchasePaidItem.ProductName, &PurchasePaidItem.ProductPrice)
		if err != nil {
			return fmt.Errorf("Err scan PurchPaidItem:%w", err)
		}
		PurchasePaidItems = append(PurchasePaidItems, PurchasePaidItem)
	}

	for _, ItemPaid := range PurchasePaidItems {
		_, err := d.DB.Exec(`
			INSERT INTO successful_purchases
			(user_id, email, payment_id, sub_start, sub_end, product_name, product_price, product_id)
			VALUES($1,$2,$3,$4,$5,$6,$7,$8)
		`, PurchasePaid.UserID, PurchasePaid.Email, PurchasePaid.PaymentID, PurchasePaid.SubStart, PurchasePaid.SubEnd, ItemPaid.ProductName, ItemPaid.ProductPrice, ItemPaid.ProductID)
		if err != nil {
			return fmt.Errorf("err insert successful_purchases: %w", err)
		}
	}
	return nil

}

func GetBasket(userID int) ([]m.Basket, error) {
	var SliceBasket []m.Basket

	rows, err := d.DB.Query(`SELECT product_id, product_name, product_price FROM basket WHERE user_id=$1`, userID)
	if err != nil {
		return []m.Basket{}, fmt.Errorf("query basket error: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var ItemFromBasket m.Basket
		err := rows.Scan(&ItemFromBasket.ProductID, &ItemFromBasket.ProductName, &ItemFromBasket.ProductPrice)
		if err != nil {
			return []m.Basket{}, fmt.Errorf("err scan Basket:%w", err)
		}
		SliceBasket = append(SliceBasket, ItemFromBasket)
	}
	return SliceBasket, nil
}

func DeleteBasketItem(ProductID int) error {
	res, err := d.DB.Exec("DELETE FROM basket WHERE user_id=$1", ProductID)
	if err != nil {
		return fmt.Errorf("err delete from basket: %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("err getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no basket item found with user_id %d", ProductID)
	}

	return nil
}

// func GetCourse(userID int) ([]string, error) {
// 	var CourseUrl []string
// 	rows, err := d.DB.Query(`
//         SELECT p.url
//         FROM successful_purchases sp
//         JOIN products p ON sp.product_id = p.id
//         WHERE sp.user_id = $1`, userID)
// 	if err != nil {
// 		return nil, fmt.Errorf("err query rows: %w", err)
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var url string
// 		if err := rows.Scan(&url); err != nil {
// 			return nil, fmt.Errorf("err scan getcourse: %w", err)
// 		}
// 		CourseUrl = append(CourseUrl, url)
// 	}

// 	if err := rows.Err(); err != nil {
// 		return nil, fmt.Errorf("rows iteration err: %w", err)
// 	}

// 	return CourseUrl, nil
// }
