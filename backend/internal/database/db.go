package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	o "github.com/vova1001/Website-Ylia-fitness/internal/otherFunc"
)

var DB *sql.DB

func DB_Conect() {
	dbHost := o.GetEnv("DB_HOST", "postgres")
	dbPort := o.GetEnv("DB_PORT", "5432")
	dbUser := o.GetEnv("DB_USER", "myuser")
	dbPassword := o.GetEnv("DB_PASSWORD", "mypassword")
	dbName := o.GetEnv("DB_NAME", "mydb")
	dbSSLMode := o.GetEnv("DB_SSLMODE", "disable")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode)

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Error conect from DB", err)
	}

	fmt.Println("DB connected")

	createTableProduct()
	createTableVideo()
	createTableBasket()
	createTablePurchaseItems()
	createTablePurchaseRequest()
	createTableSuccessfulPurchases()
}

// сами курсы (4 шт)
func createTableProduct() {
	createTable := `
	CREATE TABLE IF NOT EXISTS products (
		id SERIAL PRIMARY KEY,
		product_name TEXT NOT NULL,
		product_price DECIMAL(10,2) NOT NULL,
		currency TEXT DEFAULT 'RUB'
	);
	`
	_, err := DB.Exec(createTable)
	if err != nil {
		log.Fatal("Error created table product", err)
	}
	fmt.Println("Table products created successefully")
}

// 12 видео под каждый из курсов (4 курса, 48 видео)
func createTableVideo() {
	createTable := `
	CREATE TABLE IF NOT EXISTS video (
		id SERIAL PRIMARY KEY,
		product_id INTEGER NOT NULL,
		url TEXT NOT NULL,
		video_name TEXT NOT NULL
	);
	`
	_, err := DB.Exec(createTable)
	if err != nil {
		log.Fatal("Error created table video", err)
	}
	fmt.Println("Table video created successefully")
}

func createTablePurchaseRequest() {
	createTable := `
	CREATE TABLE IF NOT EXISTS purchase_request(
		id SERIAL PRIMARY KEY,
        user_id INTEGER NOT NULL,
        email TEXT NOT NULL,
		total_amount DECIMAL(10,2) NOT NULL,
        payment_id TEXT,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err := DB.Exec(createTable)
	if err != nil {
		log.Fatal("Error created table purchase_request", err)
	}
	fmt.Println("Table purchase_request created successefully")
}

func createTablePurchaseItems() {
	createTable := `
		CREATE TABLE IF NOT EXISTS purchase_item(
			id SERIAL PRIMARY KEY,
			purchase_request_id INTEGER NOT NULL REFERENCES purchase_request(id) ON DELETE CASCADE,
			product_id INTEGER NOT NULL,
			product_name TEXT NOT NULL,
			product_price DECIMAL(10,2) NOT NULL
		);
	`
	_, err := DB.Exec(createTable)
	if err != nil {
		log.Fatal("Error created table purchase_items", err)
	}
	fmt.Println("Table purchase_items created successefully")
}

func createTableSuccessfulPurchases() {
	createTable := `
		CREATE TABLE IF NOT EXISTS successful_purchases(
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL,
			email TEXT NOT NULL,
			product_id INTEGER NOT NULL,
			product_name TEXT NOT NULL,
			product_price DECIMAL(10,2) NOT NULL,
			payment_id TEXT,
			sub_start TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			sub_end TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		`
	_, err := DB.Exec(createTable)
	if err != nil {
		log.Fatal("Error created table successful_purchases", err)
	}
	fmt.Println("Table successful_purchases created successefully")
}

func createTableBasket() {
	createTable := `
		CREATE TABLE IF NOT EXISTS basket(
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL,
			email TEXT NOT NULL,
			product_id INTEGER NOT NULL,
			product_name TEXT NOT NULL,
			product_price DECIMAL(10,2) NOT NULL
		);
	`
	_, err := DB.Exec(createTable)
	if err != nil {
		log.Fatal("Error created table basket", err)
	}
	fmt.Println("Table basket created successefully")
}
