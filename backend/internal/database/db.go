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
	createTableBasket()
	createTablePurchaseRequest()
	createTableSuccessfulPurchases()
}

func createTableProduct() {
	createTable := `
	CREATE TABLE IF NOT EXISTS products (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		price DECIMAL(10,2) NOT NULL,
		currency TEXT DEFAULT 'RUB',
		url TEXT NOT NULL
	);
	`
	_, err := DB.Exec(createTable)
	if err != nil {
		log.Fatal("Error created table product")
	}
	fmt.Println("Table products created successefully")
}

func createTablePurchaseRequest() {
	createTable := `
	CREATE TABLE IF NOT EXISTS purchase_request(
		id SERIAL PRIMARY KEY,
        user_id INTEGER NOT NULL,
        email TEXT NOT NULL,
        product_id INTEGER NOT NULL,
        product_name TEXT NOT NULL,
        product_price DECIMAL(10,2) NOT NULL,
        payment_id TEXT,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err := DB.Exec(createTable)
	if err != nil {
		log.Fatal("Error created table purchase_request")
	}
	fmt.Println("Table purchase_request created successefully")
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
			purchased_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		`
	_, err := DB.Exec(createTable)
	if err != nil {
		log.Fatal("Error created table successful_purchases")
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
			product_price DECIMAL(10,2) NOT NULL,
		);
	`
	_, err := DB.Exec(createTable)
	if err != nil {
		log.Fatal("Error created table basket")
	}
	fmt.Println("Table basket created successefully")
}
