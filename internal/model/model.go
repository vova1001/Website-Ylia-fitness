package model

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
