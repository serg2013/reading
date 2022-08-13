package models

type Cred struct {
	Email    string `json:"email" example:"user_@gmail.com"`
	Password string `json:"password" example:"password"`
}
