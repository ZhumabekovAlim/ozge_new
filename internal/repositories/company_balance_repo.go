package repositories

import (
	"OzgeContract/internal/models"
	"database/sql"
	"errors"
)

type CompanyBalanceRepository struct {
	DB *sql.DB
}

func NewCompanyBalanceRepository(db *sql.DB) *CompanyBalanceRepository {
	return &CompanyBalanceRepository{DB: db}
}

func (r *CompanyBalanceRepository) Create(cb *models.CompanyBalance) error {
	query := `INSERT INTO company_balances (company_id, sms_signatures, ecp_signatures) VALUES (?, ?, ?)`
	_, err := r.DB.Exec(query, cb.CompanyID, cb.SMSSignatures, cb.ECPSignatures)
	return err
}

func (r *CompanyBalanceRepository) GetByCompanyID(companyID int) (*models.CompanyBalance, error) {
	var cb models.CompanyBalance
	query := `SELECT id, company_id, sms_signatures, ecp_signatures, created_at, updated_at FROM company_balances WHERE company_id = ?`
	err := r.DB.QueryRow(query, companyID).Scan(&cb.ID, &cb.CompanyID, &cb.SMSSignatures, &cb.ECPSignatures, &cb.CreatedAt, &cb.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &cb, nil
}

func (r *CompanyBalanceRepository) Update(cb *models.CompanyBalance) error {
	query := `UPDATE company_balances SET sms_signatures = ?, ecp_signatures = ? WHERE company_id = ?`
	_, err := r.DB.Exec(query, cb.SMSSignatures, cb.ECPSignatures, cb.CompanyID)
	return err
}

func (r *CompanyBalanceRepository) Delete(companyID int) error {
	_, err := r.DB.Exec(`DELETE FROM company_balances WHERE company_id = ?`, companyID)
	return err
}

func (r *CompanyBalanceRepository) SubtractSignature(companyID int, method string) error {
	var query string
	if method == "sms" {
		query = `UPDATE company_balances SET sms_signatures = sms_signatures - 1 WHERE company_id = ? AND sms_signatures > 0`
	} else if method == "ecp" {
		query = `UPDATE company_balances SET ecp_signatures = ecp_signatures - 1 WHERE company_id = ? AND ecp_signatures > 0`
	} else {
		return errors.New("invalid method")
	}
	res, err := r.DB.Exec(query, companyID)
	affected, _ := res.RowsAffected()
	if err != nil || affected == 0 {
		return errors.New("insufficient balance")
	}
	return nil
}
