package models

type Contract struct {
	ID               int    `json:"id"`
	CompanyID        int    `json:"company_id"`
	TemplateID       int    `json:"template_id"`
	ContractToken    string `json:"contract_token"`
	GeneratedPDFPath string `json:"generated_pdf_path"`
	ClientFilled     bool   `json:"client_filled"`
	CreatedAt        string `json:"created_at"`
}
