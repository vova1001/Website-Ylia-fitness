package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vova1001/Website-Ylia-fitness/internal/handlerJSON"
)

func RegisterRoutes(r *gin.Engine) {
	r.GET("", handlerJSON.GetJson)
}
