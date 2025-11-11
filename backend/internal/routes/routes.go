package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	h "github.com/vova1001/Website-Ylia-fitness/internal/handlerJSON"
	metrics "github.com/vova1001/Website-Ylia-fitness/internal/metrics"
)

func RegisterRoutes(r *gin.Engine) {
	r.Use(metrics.PrometheusMiddleware())
	auth := r.Group("/auth")
	auth.Use(h.JWT_Middleware())
	{
		auth.GET("/hi", h.GetAuthJson)
		// auth.POST("/purchase", h.PurchaseJSON)
		auth.GET("/get-course", h.GetCourseJSON)
		auth.GET("/basket")
		auth.POST("/basket/add")
	}
	r.POST("/registerUser", h.PostNewUserJson)
	r.POST("/authUser", h.PostAuthJson)
	r.POST("/fogotPassword", h.FogotPassJSON)
	r.POST("/resetPassword", h.ResetPasswordJSON)
	r.GET("/", h.GetHethJSON)
	r.POST("/webhook", h.WebhookJSON)
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
}
