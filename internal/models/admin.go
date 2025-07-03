package models

// Admin represents system administrator
// IsSuper determines whether admin has super permissions
// Password field is omitted in JSON responses

type Admin struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
	IsSuper  bool   `json:"is_super"`
}
