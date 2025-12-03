package main

import (
	"github.com/gin-gonic/gin"
	d "github.com/vova1001/Website-Ylia-fitness/internal/database"
	metrics "github.com/vova1001/Website-Ylia-fitness/internal/metrics"
	rout "github.com/vova1001/Website-Ylia-fitness/internal/routes"
)

func main() {
	metrics.MetricsInit()
	d.DB_Conect()
	r := gin.Default()
	rout.RegisterRoutes(r)
	r.Run("0.0.0.0:8080")
}
