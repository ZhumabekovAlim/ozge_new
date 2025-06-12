package repositories

import (
	"OzgeContract/internal/models"
	"database/sql"
)

type ContractRepository struct {
	DB *sql.DB
}

func NewContractRepository(db *sql.DB) *ContractRepository {
	return &ContractRepository{DB: db}
}

func (r *ContractRepository) Create(c *models.Contract) error {
	query := `INSERT INTO contracts (company_id, template_id, contract_token, generated_file_path, client_filled, method, created_at) VALUES (?, ?, ?, ?, ?, ?, NOW())`
	_, err := r.DB.Exec(query, c.CompanyID, c.TemplateID, c.ContractToken, c.GeneratedPDFPath, c.ClientFilled, c.Method)
	return err
}

func (r *ContractRepository) GetByID(id int) (*models.Contract, error) {
	query := `SELECT id, company_id, template_id, contract_token, generated_file_path, client_filled, method, created_at FROM contracts WHERE id = ?`
	row := r.DB.QueryRow(query, id)
	var c models.Contract
	err := row.Scan(&c.ID, &c.CompanyID, &c.TemplateID, &c.ContractToken, &c.GeneratedPDFPath, &c.ClientFilled, &c.Method, &c.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *ContractRepository) GetByToken(token string) (*models.Contract, error) {
	query := `SELECT id, company_id, template_id, contract_token, generated_file_path, client_filled, method, created_at FROM contracts WHERE contract_token = ?`
	row := r.DB.QueryRow(query, token)
	var c models.Contract
	err := row.Scan(&c.ID, &c.CompanyID, &c.TemplateID, &c.ContractToken, &c.GeneratedPDFPath, &c.ClientFilled, &c.Method, &c.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *ContractRepository) GetByCompanyID(companyID int) ([]models.Contract, error) {
	rows, err := r.DB.Query(`SELECT id, company_id, template_id, contract_token, generated_file_path, client_filled, method, created_at FROM contracts WHERE company_id = ?`, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contracts []models.Contract
	for rows.Next() {
		var c models.Contract
		err := rows.Scan(&c.ID, &c.CompanyID, &c.TemplateID, &c.ContractToken, &c.GeneratedPDFPath, &c.ClientFilled, &c.Method, &c.CreatedAt)
		if err != nil {
			return nil, err
		}
		contracts = append(contracts, c)
	}
	return contracts, nil
}

func (r *ContractRepository) Update(c *models.Contract) error {
	query := `UPDATE contracts SET generated_file_path = ?, client_filled = ? WHERE id = ?`
	_, err := r.DB.Exec(query, c.GeneratedPDFPath, c.ClientFilled, c.ID)
	return err
}

func (r *ContractRepository) Delete(id int) error {
	_, err := r.DB.Exec(`DELETE FROM contracts WHERE id = ?`, id)
	return err
}

func (r *ContractRepository) CreateTx(tx *sql.Tx, c *models.Contract) (int, error) {
	query := `INSERT INTO contracts (company_id, template_id, contract_token, generated_file_path, client_filled, method, created_at)
              VALUES (?, ?, ?, ?, ?, ?, NOW())`
	res, err := tx.Exec(query, c.CompanyID, c.TemplateID, c.ContractToken, c.GeneratedPDFPath, c.ClientFilled, c.Method)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return int(id), err
}

func (r *ContractFieldRepository) CreateTx(tx *sql.Tx, field *models.ContractField) error {
	query := `INSERT INTO contract_fields (contract_id, field_name, field_type) VALUES (?, ?, ?)`
	_, err := tx.Exec(query, field.ContractID, field.FieldName, field.FieldType)
	return err
}
