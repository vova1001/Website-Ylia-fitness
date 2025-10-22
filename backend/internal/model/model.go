package model

import "time"

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
	Currency     string
	CreateadAt   time.Time
	Status       string
	PaymentID    string
}
