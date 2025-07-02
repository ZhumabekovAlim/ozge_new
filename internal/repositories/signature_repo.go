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
	Search    string
	Status    *int
	Method    string
	SortBy    string
	Order     string
	CursorID  int
	Limit     int
	Direction string
}

func NewSignatureRepository(db *sql.DB) *SignatureRepository {
	return &SignatureRepository{DB: db}
}

func (r *SignatureRepository) Create(s *models.Signature) (int, error) {
	query := `INSERT INTO signatures (contract_id, client_name, client_iin, client_phone,method,status, signed_at) VALUES (?, ?, ?, ?, ?,?, NOW())`
	res, err := r.DB.Exec(query, s.ContractID, s.ClientName, s.ClientIIN, s.ClientPhone, s.Method, 1)
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
	query := `SELECT
    IFNULL(s.id, c.id) AS signature_id,
    c.id AS contract_id,
    t.name,
    IFNULL(s.client_name, '') AS client_name,
    IFNULL(s.client_iin, '') AS client_iin,
    IFNULL(DATE_FORMAT(s.signed_at, '%Y-%m-%d %H:%i:%s'), DATE_FORMAT(c.created_at, '%Y-%m-%d %H:%i:%s')) AS signed_at,
    IFNULL(s.sign_file_path, c.generated_file_path) AS sign_file_path,
    IFNULL(s.status, 0) AS status
FROM
    contracts c
        LEFT JOIN
    signatures s ON c.id = s.contract_id
        LEFT JOIN
    templates t ON t.id = c.template_id
WHERE
    c.company_id = ?;`

	rows, err := r.DB.Query(query, companyID)
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
			&s.SignedAt,
			&s.SignFilePath,
			&s.Status,
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

func (r *SignatureRepository) GetSignaturesAll(opts models.SignatureQueryOptions) ([]models.Signature, error) {
	var qb strings.Builder
	var args []interface{}

	qb.WriteString(`
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
            co.name,
            co.iin
        FROM signatures s
        LEFT JOIN contracts c ON c.id = s.contract_id
        LEFT JOIN templates t ON t.id = c.template_id
        LEFT JOIN signature_field_values sfv ON s.id = sfv.signature_id
        LEFT JOIN companies co ON c.company_id = co.id
        WHERE 1=1
    `)

	// Поиск
	if opts.Search != "" {
		s := "%" + opts.Search + "%"
		qb.WriteString(`
            AND (
                CAST(s.id AS CHAR) LIKE ? OR
                s.client_name LIKE ? OR
                s.client_phone LIKE ? OR
                t.name LIKE ? OR
                co.name LIKE ? OR
                co.iin LIKE ? OR
                s.method LIKE ? OR
                CAST(status AS CHAR) LIKE ? OR
                DATE_FORMAT(s.signed_at, '%Y-%m-%d') LIKE ?
            )
        `)
		args = append(args, s, s, s, s, s, s, s, s, s)
	}

	if opts.Status != nil {
		qb.WriteString(" AND status = ?")
		args = append(args, *opts.Status)
	}

	if opts.Method != "" {
		qb.WriteString(" AND s.method = ?")
		args = append(args, opts.Method)
	}

	// Сортировка и направление
	var orderBy string
	switch opts.SortBy {
	case "client_name":
		orderBy = "s.client_name"
	case "signed_at":
		orderBy = "s.signed_at"
	case "id", "", "default":
		fallthrough
	default:
		orderBy = "s.id"
	}

	// Устанавливаем order и comparator
	order := "ASC"
	comparator := ">"
	switch strings.ToUpper(opts.Order) {
	case "DESC":
		order = "DESC"
		comparator = "<" // очень важно!
	case "ASC":
		order = "ASC"
		comparator = ">"
	default:
		order = "ASC"
		comparator = ">"
	}

	// Если direction = prev — инвертируем порядок и comparator
	if opts.Direction == "prev" {
		if order == "ASC" {
			comparator = "<"
			order = "DESC"
		} else {
			comparator = ">"
			order = "ASC"
		}
	}

	// Добавить условие курсора
	if opts.CursorID > 0 {
		qb.WriteString(" AND s.id " + comparator + " ?")
		args = append(args, opts.CursorID)
	}

	if opts.Limit == 0 {
		opts.Limit = 20
	}

	qb.WriteString(" ORDER BY " + orderBy + " " + order)
	qb.WriteString(" LIMIT ?")
	args = append(args, opts.Limit)

	rows, err := r.DB.Query(qb.String(), args...)
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
			&s.CompanyIIN,
		)
		if err != nil {
			return nil, err
		}
		signatures = append(signatures, s)
	}

	// Реверс при prev, чтобы вернуть в правильном порядке
	if opts.Direction == "prev" {
		for i, j := 0, len(signatures)-1; i < j; i, j = i+1, j-1 {
			signatures[i], signatures[j] = signatures[j], signatures[i]
		}
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

func (r *SignatureRepository) GetStatusSummary() (*models.SignatureStatusSummary, error) {
	summary := &models.SignatureStatusSummary{}
	query := `SELECT 
                COUNT(id) AS total,
                SUM(CASE WHEN status = 1 THEN 1 ELSE 0 END) AS signed,
                SUM(CASE WHEN status = 0 THEN 1 ELSE 0 END) AS pending,
                SUM(CASE WHEN status = 3 THEN 1 ELSE 0 END) AS declined
        	FROM signatures`
	row := r.DB.QueryRow(query)
	if err := row.Scan(&summary.Total, &summary.Signed, &summary.Pending, &summary.Declined); err != nil {
		return nil, err
	}
	return summary, nil
}
