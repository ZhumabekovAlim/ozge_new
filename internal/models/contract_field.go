package models

type ContractField struct {
	ID         int    `json:"id"`
	ContractID int    `json:"contract_id"`
	FieldName  string `json:"field_name"`
	FieldType  string `json:"field_type"` // default: "text"
}
