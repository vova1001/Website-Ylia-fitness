package main

import (
	"github.com/gin-gonic/gin"
	d "github.com/vova1001/Website-Ylia-fitness/internal/database"
	metrics "github.com/vova1001/Website-Ylia-fitness/internal/metrics"
	rout "github.com/vova1001/Website-Ylia-fitness/internal/routes"

	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	shopID := "1199000"
	apiKey := "live_mwZWuqw-qJGp7UYnoHzk5UGA-2dIEHviUQ4Vrc3rHIo"

	client := &http.Client{Timeout: 30 * time.Second}

	req, err := http.NewRequest("GET", "https://api.yookassa.ru/v3/payments", nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Базовая авторизация
	auth := base64.StdEncoding.EncodeToString([]byte(shopID + ":" + apiKey))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("User-Agent", "GoTestClient/1.0")

	fmt.Println("Sending request to Yookassa...")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	fmt.Println("Status:", resp.Status)
	fmt.Println("Body:", string(body))

	metrics.MetricsInit()
	d.DB_Conect()
	r := gin.Default()
	rout.RegisterRoutes(r)
	r.Run("0.0.0.0:8080")
}
