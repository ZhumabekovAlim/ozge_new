package repositories

import (
	"OzgeContract/internal/models"
	"database/sql"
)

type ContractFieldRepository struct {
	DB *sql.DB
}

func NewContractFieldRepository(db *sql.DB) *ContractFieldRepository {
	return &ContractFieldRepository{DB: db}
}

func (r *ContractFieldRepository) Create(field *models.ContractField) error {
	query := `INSERT INTO contract_fields (contract_id, field_name, field_type) VALUES (?, ?, ?)`
	_, err := r.DB.Exec(query, field.ContractID, field.FieldName, field.FieldType)
	return err
}

func (r *ContractFieldRepository) GetByContractID(contractID int) ([]models.ContractField, error) {
	rows, err := r.DB.Query(`SELECT id, contract_id, field_name, field_type FROM contract_fields WHERE contract_id = ?`, contractID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fields []models.ContractField
	for rows.Next() {
		var f models.ContractField
		err := rows.Scan(&f.ID, &f.ContractID, &f.FieldName, &f.FieldType)
		if err != nil {
			return nil, err
		}
		fields = append(fields, f)
	}
	return fields, nil
}
