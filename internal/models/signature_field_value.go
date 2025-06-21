package models

type SignatureFieldValue struct {
	ID              int    `json:"id"`
	SignatureID     int    `json:"signature_id"`
	ContractFieldID int    `json:"contract_field_id"`
	FieldName       string `json:"field_name,omitempty"`
	FieldValue      string `json:"field_value"`
}

// SignatureFieldValueDTO is used when creating signature field values together with a signature.
type SignatureFieldValueDTO struct {
	ContractFieldID int    `json:"contract_field_id"`
	FieldValue      string `json:"field_value"`
}
