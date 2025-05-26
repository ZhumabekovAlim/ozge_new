package repositories

import (
	"OzgeContract/internal/models"
	"database/sql"
)

type SignatureFieldValueRepository struct {
	DB *sql.DB
}

func NewSignatureFieldValueRepository(db *sql.DB) *SignatureFieldValueRepository {
	return &SignatureFieldValueRepository{DB: db}
}

func (r *SignatureFieldValueRepository) Create(value *models.SignatureFieldValue) error {
	query := `INSERT INTO signature_field_values (signature_id, contract_field_id, field_value) VALUES (?, ?, ?)`
	_, err := r.DB.Exec(query, value.SignatureID, value.ContractFieldID, value.FieldValue)
	return err
}

func (r *SignatureFieldValueRepository) GetBySignatureID(signatureID int) ([]models.SignatureFieldValue, error) {
	rows, err := r.DB.Query(`
		SELECT sfv.id, sfv.signature_id, sfv.contract_field_id, sfv.field_value, cf.field_name 
		FROM signature_field_values sfv
		JOIN contract_fields cf ON sfv.contract_field_id = cf.id
		WHERE sfv.signature_id = ?`, signatureID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.SignatureFieldValue
	for rows.Next() {
		var v models.SignatureFieldValue
		err := rows.Scan(&v.ID, &v.SignatureID, &v.ContractFieldID, &v.FieldValue, &v.FieldName)
		if err != nil {
			return nil, err
		}
		list = append(list, v)
	}
	return list, nil
}
