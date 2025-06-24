package models

type Company struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	IIN      string `json:"iin"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}
