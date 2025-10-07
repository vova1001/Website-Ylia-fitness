package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func main() {
	// DSN для локального запуска
	dsn := "host=localhost port=5432 user=myuser password=mypassword dbname=mydb sslmode=disable"

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("Не удалось подключиться к БД:", err)
	}

	fmt.Println("✅ Подключение к БД успешно!")

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(200, "API работает!")
	})

	r.Run("0.0.0.0:8080")
}
