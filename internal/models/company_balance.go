package models

type CompanyBalance struct {
	ID            int    `json:"id"`
	CompanyID     int    `json:"company_id"`
	SMSSignatures int    `json:"sms_signatures"`
	ECPSignatures int    `json:"ecp_signatures"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}
