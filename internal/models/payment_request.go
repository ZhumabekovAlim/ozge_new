package models

import "time"

type PaymentRequest struct {
	ID           int        `json:"id"`
	CompanyID    int        `json:"company_id"`
	TariffPlanID *int       `json:"tariff_plan_id,omitempty"`
	SMSCount     int        `json:"sms_count"`
	ECPCount     int        `json:"ecp_count"`
	TotalAmount  float64    `json:"total_amount"`
	Status       string     `json:"status"`
	PaymentURL   string     `json:"payment_url"`
	PaymentRef   string     `json:"payment_ref"`
	CreatedAt    time.Time  `json:"created_at"`
	PaidAt       *time.Time `json:"paid_at,omitempty"`
}
