package models

import "time"

type PaymentRequest struct {
	ID            int        `json:"id"`
	CompanyID     int        `json:"company_id"`
	TariffPlanID  *int       `json:"tariff_plan_id,omitempty"`
	SMSCount      int        `json:"sms_count"`
	ECPCount      int        `json:"ecp_count"`
	SMSSignatures int        `json:"sms_signatures,omitempty"`
	ECPSignatures int        `json:"ecp_signatures,omitempty"`
	TotalAmount   float64    `json:"total_amount,omitempty"`
	Status        string     `json:"status,omitempty"`
	PaymentURL    string     `json:"payment_url,omitempty"`
	PaymentRef    string     `json:"payment_ref,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	PaidAt        *time.Time `json:"paid_at,omitempty"`
	CompanyName   string     `json:"company_name,omitempty"`
	CompanyIIN    string     `json:"company_iin,omitempty"`
}
