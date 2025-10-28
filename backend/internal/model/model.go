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
	ID       int    `json:"id,omitempty"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Token struct {
	JWT_Token string `json:"jwt_token"`
}

type FogotPass struct {
	Email string `json:"email"`
}

type NewPass struct {
	NewPass string `json:"new_pass"`
	Token   string `json:"token"`
}

type PurchaseRequest struct {
	IdProduct int    `json:"id_product"`
	UserToken string `json:"user_token"`
}

type TokenNewPass struct {
	EmailToToken string
	HashToken    string
	TimeLife     time.Time
	Used         bool
}

type Purchase struct {
	UserID       int
	Email        string
	ProductID    int
	ProductName  string
	ProductPrice float64
	CreateadAt   time.Time
	PaymentID    string
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
	Capture      bool   `json:"capture"`
	Description  string `json:"description"`
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
		ID     string `json:"id"`
		Status string `json:"status"`
		Paid   bool   `json:"paid"`
	} `json:"object"`
}
