package routes

import (
	"github.com/gin-gonic/gin"
	h "github.com/vova1001/Website-Ylia-fitness/internal/handlerJSON"
)

func RegisterRoutes(r *gin.Engine) {
	auth := r.Group("/auth")
	auth.Use(h.JWT_Middleware())
	{
		auth.GET("/hi", h.GetAuthJson)
		auth.POST("/purchase", h.PurchaseJSON)
		auth.GET("/get-course", h.GetCourseJSON)
	}
	r.POST("/registerUser", h.PostNewUserJson)
	r.POST("/authUser", h.PostAuthJson)
	r.POST("/fogotPassword", h.FogotPassJSON)
	r.POST("/resetPassword", h.ResetPasswordJSON)
	r.GET("/", h.GetHethJSON)
	r.POST("/webhook", h.WebhookJSON)
}
