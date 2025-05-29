package models

type TariffPlan struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Discount  float64 `json:"discount"`
	FromCount int     `json:"from_count"`
	ToCount   int     `json:"to_count"`
	IsActive  bool    `json:"is_active"`
	CreatedAt string  `json:"created_at"`
}
