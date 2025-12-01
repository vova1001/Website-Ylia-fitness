package main

import (
	"github.com/gin-gonic/gin"
	d "github.com/vova1001/Website-Ylia-fitness/internal/database"
	metrics "github.com/vova1001/Website-Ylia-fitness/internal/metrics"
	rout "github.com/vova1001/Website-Ylia-fitness/internal/routes"

	"log"
	"net"
	"time"
)

func main() {
	ip := "109.235.165.99:443"
	log.Println("Testing TCP connection to", ip)
	conn, err := net.DialTimeout("tcp", ip, 10*time.Second)
	if err != nil {
		log.Println("TCP connection failed:", err)
	} else {
		log.Println("TCP connection succeeded")
		conn.Close()
	}

	metrics.MetricsInit()
	d.DB_Conect()
	r := gin.Default()
	rout.RegisterRoutes(r)
	r.Run("0.0.0.0:8080")
}
