package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	h "github.com/vova1001/Website-Ylia-fitness/internal/handlerJSON"
	metrics "github.com/vova1001/Website-Ylia-fitness/internal/metrics"

	"time"
)

func RegisterRoutes(r *gin.Engine) {
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	r.Use(metrics.PrometheusMiddleware())
	auth := r.Group("/auth")
	auth.Use(h.JWT_Middleware())
	{
		auth.GET("/lifeTime")
		auth.GET("/hi", h.GetAuthJson)
		auth.GET("/purchase", h.GetPurchaseJSON)
		auth.POST("/purchase/extension", h.PostPurchaseExtension)
		auth.GET("/basket", h.GetBasketJSON)
		auth.POST("/basket/add", h.AddBasketJSON)
		auth.DELETE("/basket/item", h.DeleteBasketJSON)
		auth.GET("/getCourse", h.GetCourseJSON)
		auth.POST("/showVideo", h.PostVideoJSON)
	}
	r.POST("/registerUser", h.PostNewUserJson)
	r.POST("/authUser", h.PostAuthJson)
	r.POST("/fogotPassword", h.FogotPassJSON)
	r.POST("/resetPassword", h.ResetPasswordJSON)
	r.GET("/", h.GetHelthJSON)
	r.POST("/webhook", h.WebhookJSON)
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
}
