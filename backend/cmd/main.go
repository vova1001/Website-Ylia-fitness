package main

import (
	"github.com/gin-gonic/gin"
	d "github.com/vova1001/Website-Ylia-fitness/internal/database"
	rout "github.com/vova1001/Website-Ylia-fitness/internal/routes"
)

func main() {

	d.DB_Conect()
	r := gin.Default()
	rout.RegisterRoutes(r)
	r.Run("0.0.0.0:8080")
}
