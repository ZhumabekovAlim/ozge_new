package models

type Contract struct {
	ID               int    `json:"id"`
	CompanyID        int    `json:"company_id"`
	TemplateID       int    `json:"template_id"`
	ContractToken    string `json:"contract_token"`
	SerialNumber     string `json:"serial_number"`
	GeneratedPDFPath string `json:"generated_pdf_path"`
	Method           string `json:"method"`
	CompanySign      bool   `json:"company_sign"`
	ClientFilled     bool   `json:"client_filled"`
	CreatedAt        string `json:"created_at"`
	CompanyName      string `json:"company_name,omitempty"`
	CompanyIIN       string `json:"company_iin,omitempty"`
	CompanyPhone     string `json:"company_phone,omitempty"`
}

type CreateContractRequest struct {
	CompanyID        int                `json:"company_id"`
	TemplateID       int                `json:"template_id"`
	GeneratedPDFPath string             `json:"generated_pdf_path"`
	SerialNumber     string             `json:"serial_number"`
	ClientFilled     bool               `json:"client_filled"`
	CompanySign      bool               `json:"company_sign"`
	Method           string             `json:"method"`
	Fields           []ContractFieldDTO `json:"fields"` // новая структура (без contract_id)
}
type ContractFieldDTO struct {
	FieldName string `json:"field_name"`
	FieldType string `json:"field_type"`
}

// ContractDetails represents a contract together with its additional fields.
type ContractDetails struct {
	Contract
	Fields []ContractField `json:"fields"`
}
