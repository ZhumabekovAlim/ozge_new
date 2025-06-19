package repositories

import (
	"OzgeContract/internal/models"
	"database/sql"
	_ "fmt"
	"strings"
)

type SignatureRepository struct {
	DB *sql.DB
}

type SignatureQueryOptions struct {
	Search   string
	Status   *int
	Method   string
	SortBy   string
	Order    string
	CursorID int
	Limit    int
}

func NewSignatureRepository(db *sql.DB) *SignatureRepository {
	return &SignatureRepository{DB: db}
}

func (r *SignatureRepository) Create(s *models.Signature) (int, error) {
	query := `INSERT INTO signatures (contract_id, client_name, client_iin, client_phone, signed_at) VALUES (?, ?, ?, ?, NOW())`
	res, err := r.DB.Exec(query, s.ContractID, s.ClientName, s.ClientIIN, s.ClientPhone)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (r *SignatureRepository) GetByID(id int) (*models.Signature, error) {
	var s models.Signature
	query := `SELECT id, contract_id, client_name, client_iin, client_phone, signed_at, sign_file_path FROM signatures WHERE id = ?`
	err := r.DB.QueryRow(query, id).Scan(&s.ID, &s.ContractID, &s.ClientName, &s.ClientIIN, &s.ClientPhone, &s.SignedAt, &s.SignFilePath)
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

func (r *SignatureRepository) GetContractsByCompanyID(companyID int) ([]models.Signature, error) {
	query := `SELECT s.id, t.name, client_name, client_iin, signed_at, sign_file_path FROM signatures s
    LEFT JOIN contracts c on c.id = s.contract_id
    LEFT JOIN templates t on t.id = c.template_id
    LEFT JOIN signature_field_values sfv on s.id = sfv.signature_id
    WHERE c.company_id = ?`

	rows, err := r.DB.Query(query, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var signatures []models.Signature
	for rows.Next() {
		var s models.Signature
		err := rows.Scan(&s.ID, &s.TemplateName, &s.ClientName, &s.ClientIIN, &s.SignedAt, &s.SignFilePath)
		if err != nil {
			return nil, err
		}
		signatures = append(signatures, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return signatures, nil
}

func (r *SignatureRepository) GetSignaturesAll(opts SignatureQueryOptions) ([]models.Signature, error) {
	query := `
                SELECT
                        s.id,
                        contract_id,
                        t.name,
                        s.client_name,
                        s.client_iin,
                        s.client_phone,
                        s.method,
                        status,
                        s.signed_at,
                        s.sign_file_path,
                        co.name
                FROM signatures s
                LEFT JOIN contracts c ON c.id = s.contract_id
                LEFT JOIN templates t ON t.id = c.template_id
                LEFT JOIN signature_field_values sfv ON s.id = sfv.signature_id
                LEFT JOIN companies co ON c.company_id = co.id
                WHERE s.id >= ?`
	args := []interface{}{opts.CursorID}

	if opts.Search != "" {
		s := "%" + opts.Search + "%"
		query += ` AND (
                        CAST(s.id AS CHAR) LIKE ? OR
                        s.client_name LIKE ? OR
                        s.client_phone LIKE ? OR
                        t.name LIKE ? OR
                        co.name LIKE ? OR
                        s.method LIKE ? OR
                        CAST(status AS CHAR) LIKE ? OR
                        DATE_FORMAT(s.signed_at, '%Y-%m-%d') LIKE ?
                )`
		args = append(args, s, s, s, s, s, s, s, s)
	}
	if opts.Status != nil {
		query += " AND status = ?"
		args = append(args, *opts.Status)
	}
	if opts.Method != "" {
		query += " AND s.method = ?"
		args = append(args, opts.Method)
	}

	orderBy := "s.id"
	switch opts.SortBy {
	case "client_name":
		orderBy = "s.client_name"
	case "signed_at":
		orderBy = "s.signed_at"
	}

	order := "ASC"
	if strings.ToUpper(opts.Order) == "DESC" {
		order = "DESC"
	}

	if opts.Limit == 0 {
		opts.Limit = 20
	}

	query += " ORDER BY " + orderBy + " " + order + " LIMIT ?"
	args = append(args, opts.Limit)

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var signatures []models.Signature
	for rows.Next() {
		var s models.Signature
		err := rows.Scan(
			&s.ID,
			&s.ContractID,
			&s.TemplateName,
			&s.ClientName,
			&s.ClientIIN,
			&s.ClientPhone,
			&s.Method,
			&s.Status,
			&s.SignedAt,
			&s.SignFilePath,
			&s.CompanyName,
		)
		if err != nil {
			return nil, err
		}
		signatures = append(signatures, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return signatures, nil
}

func (r *SignatureRepository) Delete(id int) error {
	_, err := r.DB.Exec(`DELETE FROM signatures WHERE id = ?`, id)
	return err
}

func (r *SignatureRepository) UpdateSignFilePath(signatureID int, path string) error {
	_, err := r.DB.Exec("UPDATE signatures SET sign_file_path = ? WHERE id = ?", path, signatureID)
	return err
}
