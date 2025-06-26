package models

type LastSignedContract struct {
	ContractID   int    `json:"contract_id"`
	ClientName   string `json:"client_name"`
	TemplateName string `json:"template_name"`
	SignedAt     string `json:"signed_at"`
	Status       string `json:"status"` // всегда "Подписан"
}

type CompanyStats struct {
	CompanyID        int                  `json:"company_id"`
	TotalSigned      int                  `json:"total_signed"`
	UniqueSignatures int                  `json:"unique_signatures"`
	SignedThisMonth  int                  `json:"signed_this_month"`
	LastSigned       []LastSignedContract `json:"last_signed_contracts"`
}
