package models

type Signature struct {
	ID          int    `json:"id"`
	ContractID  int    `json:"contract_id"`
	ClientName  string `json:"client_name"`
	ClientIIN   string `json:"client_iin"`
	ClientPhone string `json:"client_phone"`
	SignedAt    string `json:"signed_at"`
	Method      string `json:"method"`
}
