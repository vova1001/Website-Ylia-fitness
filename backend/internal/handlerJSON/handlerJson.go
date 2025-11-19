package handlerJSON

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"

	h "github.com/vova1001/Website-Ylia-fitness/internal/handler"
	m "github.com/vova1001/Website-Ylia-fitness/internal/model"
)

func GetHethJSON(ctx *gin.Context) {
	ctx.JSON(200, "Server is running!")
}

func GetAuthJson(ctx *gin.Context) {
	ctx.JSON(200, "Authorization successful")
}

func PostNewUserJson(ctx *gin.Context) {
	var NewUser m.User
	err := ctx.ShouldBindJSON(&NewUser)
	if err != nil {
		ctx.JSON(400, gin.H{"err": "err json"})
		return
	}
	err = h.RegisterNewUser(NewUser)
	if err != nil {
		ctx.JSON(401, gin.H{"err": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"Sucseseful": "User add table"})

}

func PostAuthJson(ctx *gin.Context) {
	var User m.User
	err := ctx.ShouldBindJSON(&User)
	if err != nil {
		ctx.JSON(400, gin.H{"err": "err json"})
	}
	token, err := h.AuthUser(User)
	if err != nil {
		ctx.JSON(401, gin.H{"err": err.Error()})
	}
	ctx.JSON(200, token)
}

func JWT_Middleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			ctx.JSON(401, gin.H{"err": "Start header not Bearer "})
			ctx.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		secretKey := []byte(os.Getenv("JWT_SECRET"))

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("err SignMethod")
			}
			return secretKey, nil
		})
		if err != nil || !token.Valid {
			ctx.JSON(401, gin.H{"error": "invalid token"})
			ctx.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			ctx.JSON(401, gin.H{"error": "invalid token claims"})
			ctx.Abort()
			return
		}

		if life, ok := claims["timeLife"].(float64); ok {
			if time.Now().Unix() > int64(life) {
				ctx.JSON(401, gin.H{"error": "token dead"})
				ctx.Abort()
				return
			}
		}

		userIDFloat, ok := claims["userID"].(float64)
		if !ok {
			ctx.JSON(401, gin.H{"error": "invalid userID in token"})
			ctx.Abort()
			return
		}

		ctx.Set("userID", int(userIDFloat))
		ctx.Set("userEmail", claims["email"])

		ctx.Next()

	}
}

func FogotPassJSON(ctx *gin.Context) {
	var email m.FogotPass
	err := ctx.ShouldBindJSON(&email)
	if err != nil {
		ctx.JSON(400, gin.H{"err": "err json"})
		return
	}
	token, err := h.FogotPass(email)
	if err != nil {
		ctx.JSON(500, gin.H{"err": err.Error()})
		return
	}
	ctx.JSON(200, token)
}

func ResetPasswordJSON(ctx *gin.Context) {
	var NewPass m.NewPass
	err := ctx.ShouldBindJSON(&NewPass)
	if err != nil {
		ctx.JSON(400, gin.H{"err": "err json"})
		return
	}
	err = h.ResetPassword(NewPass)
	if err != nil {
		ctx.JSON(500, gin.H{"err": err.Error()})
		return
	}
	ctx.JSON(200, "new password has been successfully set")
}

func AddBasketJSON(ctx *gin.Context) {
	var Product m.ProductAddBasket
	err := ctx.ShouldBindJSON(&Product)
	if err != nil {
		ctx.JSON(400, gin.H{"err": "err json"})
		return
	}

	UserID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(401, gin.H{"err": "User not found"})
		return
	}

	Email, exists := ctx.Get("userEmail")
	if !exists {
		ctx.JSON(401, gin.H{"err": "Email not found"})
		return
	}

	UserIDint := UserID.(int)
	EmailStr := Email.(string)
	IdProduct := Product.IdProduct
	res, err := h.ProductAddBasket(UserIDint, IdProduct, EmailStr)
	if err != nil {
		ctx.JSON(500, gin.H{"err": err.Error()})
	}
	ctx.JSON(200, res)

}

func PurchaseJSON(ctx *gin.Context) {
	UserID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(401, gin.H{"err": "User not found"})
		return
	}

	UserIDint := UserID.(int)
	URL_payment, err := h.PurchesRequest(UserIDint)
	if err != nil {
		ctx.JSON(500, gin.H{"err": err.Error()})
		return
	}
	ctx.JSON(200, URL_payment)
}

func WebhookJSON(ctx *gin.Context) {
	var Webhook m.YookassaWebhook
	err := ctx.ShouldBindJSON(&Webhook)
	if err != nil {
		ctx.JSON(400, gin.H{"err": "err json"})
		return
	}
	err = h.WebhookY(Webhook)
	if err != nil {
		ctx.JSON(500, gin.H{"err": err.Error()})
		return
	}
}

func GetBasketJSON(ctx *gin.Context) {
	UserID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(401, gin.H{"err": "User not found"})
		return
	}
	UserIDint := UserID.(int)
	SliceBasket, err := h.GetBasket(UserIDint)
	if err != nil {
		ctx.JSON(500, gin.H{"err": err.Error()})
	}
	ctx.JSON(200, SliceBasket)

}

// func GetCourseJSON(ctx *gin.Context) {
// 	userID, exists := ctx.Get("userID")
// 	if !exists {
// 		ctx.JSON(401, gin.H{"err": "User not found"})
// 		return
// 	}
// 	UserIDint := userID.(int)
// 	ResURL, err := h.GetCourse(UserIDint)
// 	if err != nil {
// 		ctx.JSON(500, gin.H{"err": err.Error()})
// 		return
// 	}
// 	ctx.JSON(200, ResURL)
// }
