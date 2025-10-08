package routes

import (
	"github.com/gin-gonic/gin"
	h "github.com/vova1001/Website-Ylia-fitness/internal/handlerJSON"
)

func RegisterRoutes(r *gin.Engine) {
	r.GET("", h.GetJson)
	r.POST("/newUser", h.PostNewUser)
}
