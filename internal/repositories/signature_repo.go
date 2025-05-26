package repositories

import (
	"OzgeContract/internal/models"
	"database/sql"
)

type SignatureRepository struct {
	DB *sql.DB
}

func NewSignatureRepository(db *sql.DB) *SignatureRepository {
	return &SignatureRepository{DB: db}
}

func (r *SignatureRepository) Create(s *models.Signature) error {
	query := `INSERT INTO signatures (contract_id, client_name, client_iin, client_phone, signed_at) VALUES (?, ?, ?, ?, NOW())`
	_, err := r.DB.Exec(query, s.ContractID, s.ClientName, s.ClientIIN, s.ClientPhone)
	return err
}

func (r *SignatureRepository) GetByID(id int) (*models.Signature, error) {
	var s models.Signature
	query := `SELECT id, contract_id, client_name, client_iin, client_phone, signed_at FROM signatures WHERE id = ?`
	err := r.DB.QueryRow(query, id).Scan(&s.ID, &s.ContractID, &s.ClientName, &s.ClientIIN, &s.ClientPhone, &s.SignedAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *SignatureRepository) GetByContractID(contractID int) (*models.Signature, error) {
	var s models.Signature
	query := `SELECT id, contract_id, client_name, client_iin, client_phone, signed_at FROM signatures WHERE contract_id = ?`
	err := r.DB.QueryRow(query, contractID).Scan(&s.ID, &s.ContractID, &s.ClientName, &s.ClientIIN, &s.ClientPhone, &s.SignedAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *SignatureRepository) Delete(id int) error {
	_, err := r.DB.Exec(`DELETE FROM signatures WHERE id = ?`, id)
	return err
}
