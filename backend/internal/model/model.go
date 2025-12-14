package model

import (
	"net/http"
	"time"
)

type Task struct {
	Id   int    `json:"id,omitempty"`
	Name string `json:"name"`
	Msg  string `jsom:"msg"`
}

type User struct {
	ID       int    `json:"id,omitempty" `
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Token struct {
	JWT_Token string `json:"jwt_token" binding:"required"`
}

type FogotPass struct {
	Email string `json:"email" binding:"required"`
}

type NewPass struct {
	NewPass string `json:"new_pass" binding:"required"`
	Token   string `json:"token" binding:"required"`
}

type ProductAddBasket struct {
	IdProduct int `json:"id_product" binding:"required"`
}

type TokenNewPass struct {
	EmailToToken string
	HashToken    string
	TimeLife     time.Time
	Used         bool
}

type PurchaseRequest struct {
	ID          int
	UserID      int
	Email       string
	TotalAmount float64
	CreateadAt  time.Time
	PaymentID   string
	Items       []PurchaseItem
}

type PurchaseItem struct {
	ID                int
	PurchaseRequestID int
	ProductID         int
	ProductName       string
	ProductPrice      float64
}

type PurchasePaid struct {
	ID        int
	UserID    int
	Email     string
	PaymentID string
	SubStart  time.Time
	SubEnd    time.Time
}

type Basket struct {
	UserID       int
	Email        string
	ProductID    int
	ProductName  string
	ProductPrice float64
}

type VideoResponse struct {
	URL       string
	VideoName string
}

type DeleteBasketItem struct {
	ID int `json:"delete_item_id"`
}

type YookassaClient struct {
	ShopID  string
	BaseURL string
	ApiKey  string
	Client  *http.Client
}

type YookassaPaymentRequest struct {
	Amount struct {
		Value    string `json:"value"`
		Currency string `json:"currency"`
	} `json:"amount"`
	Capture      bool              `json:"capture"`
	Description  string            `json:"description"`
	Metadata     map[string]string `json:"metadata,omitempty"`
	Confirmation struct {
		Type      string `json:"type"`
		ReturnURL string `json:"return_url"`
	} `json:"confirmation"`
}

type YookassaPaymentResponse struct {
	ID           string `json:"id"`
	Status       string `json:"status"`
	Confirmation struct {
		ConfirmationURL string `json:"confirmation_url"`
	} `json:"confirmation"`
}

type YookassaWebhook struct {
	Event  string `json:"event"`
	Object struct {
		ID       string            `json:"id"`
		Status   string            `json:"status"`
		Paid     bool              `json:"paid"`
		Metadata map[string]string `json:"metadata"`
	} `json:"object"`
}
