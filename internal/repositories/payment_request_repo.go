package repositories

import (
	"OzgeContract/internal/models"
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
)

type PaymentRequestRepository struct {
	DB *sql.DB
}

type PaymentRequestQueryOptions struct {
	Search   string
	Status   string
	SortBy   string
	Order    string
	CursorID int
	Limit    int
}

func NewPaymentRequestRepository(db *sql.DB) *PaymentRequestRepository {
	return &PaymentRequestRepository{DB: db}
}

func (r *PaymentRequestRepository) Create(p *models.PaymentRequest) error {
	query := `
	INSERT INTO payment_requests (company_id, tariff_plan_id, sms_count, ecp_count, total_amount, status, payment_url, payment_ref)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)	
	`
	_, err := r.DB.Exec(query, p.CompanyID, p.TariffPlanID, p.SMSCount, p.ECPCount, p.TotalAmount, p.Status, p.PaymentURL, p.PaymentRef)
	return err
}

func (r *PaymentRequestRepository) GetByID(id int) (*models.PaymentRequest, error) {
	query := `
        SELECT pr.id, pr.company_id, pr.tariff_plan_id, pr.sms_count, pr.ecp_count, pr.total_amount,
               pr.status, pr.payment_url, pr.payment_ref, pr.created_at, pr.paid_at, c.name
        FROM payment_requests pr
        LEFT JOIN companies c ON pr.company_id = c.id
        WHERE pr.id = ?
        `
	row := r.DB.QueryRow(query, id)
	var p models.PaymentRequest
	err := row.Scan(&p.ID, &p.CompanyID, &p.TariffPlanID, &p.SMSCount, &p.ECPCount, &p.TotalAmount, &p.Status, &p.PaymentURL, &p.PaymentRef, &p.CreatedAt, &p.PaidAt, &p.CompanyName)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PaymentRequestRepository) GetByCompany(companyID int) ([]models.PaymentRequest, error) {
	query := `
        SELECT pr.id, pr.company_id, pr.tariff_plan_id, pr.sms_count, pr.ecp_count, pr.total_amount,
               pr.status, pr.payment_url, pr.payment_ref, pr.created_at, pr.paid_at, c.name
        FROM payment_requests pr
        LEFT JOIN companies c ON pr.company_id = c.id
        WHERE pr.company_id = ?
        `
	rows, err := r.DB.Query(query, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.PaymentRequest
	for rows.Next() {
		var p models.PaymentRequest
		err := rows.Scan(&p.ID, &p.CompanyID, &p.TariffPlanID, &p.SMSCount, &p.ECPCount, &p.TotalAmount, &p.Status, &p.PaymentURL, &p.PaymentRef, &p.CreatedAt, &p.PaidAt, &p.CompanyName)
		if err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, nil
}

func (r *PaymentRequestRepository) GetAll(ctx context.Context, opts PaymentRequestQueryOptions) ([]models.PaymentRequest, error) {
	query := `
                SELECT
                        pr.id, pr.company_id, pr.tariff_plan_id, pr.sms_count, pr.ecp_count,
                        pr.total_amount, pr.status, pr.payment_url, pr.payment_ref,
                        pr.created_at, pr.paid_at, c.name
                FROM payment_requests pr
                LEFT JOIN companies c ON pr.company_id = c.id
                WHERE pr.id <= ?`
	args := []interface{}{opts.CursorID}

	if opts.Search != "" {
		s := "%" + opts.Search + "%"
		query += ` AND (
                        CAST(pr.id AS CHAR) LIKE ? OR
                        c.name LIKE ? OR
                        CAST(pr.sms_count AS CHAR) LIKE ? OR
                        CAST(pr.ecp_count AS CHAR) LIKE ? OR
                        CAST(pr.total_amount AS CHAR) LIKE ? OR
                        pr.payment_url LIKE ? OR
                        pr.payment_ref LIKE ? OR
                        DATE_FORMAT(pr.created_at, '%Y-%m-%d') LIKE ? OR
                        DATE_FORMAT(pr.paid_at, '%Y-%m-%d') LIKE ?
                )`
		args = append(args, s, s, s, s, s, s, s, s, s)
	}
	if opts.Status != "" {
		query += " AND pr.status = ?"
		args = append(args, opts.Status)
	}

	orderBy := "pr.id"
	switch opts.SortBy {
	case "company_name":
		orderBy = "c.name"
	case "total_amount":
		orderBy = "pr.total_amount"
	case "created_at":
		orderBy = "pr.created_at"
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

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var list []models.PaymentRequest
	for rows.Next() {
		var p models.PaymentRequest
		err := rows.Scan(
			&p.ID, &p.CompanyID, &p.TariffPlanID, &p.SMSCount, &p.ECPCount,
			&p.TotalAmount, &p.Status, &p.PaymentURL, &p.PaymentRef,
			&p.CreatedAt, &p.PaidAt, &p.CompanyName,
		)
		if err != nil {
			log.Printf("scan error: %v", err)
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		list = append(list, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	log.Printf("fetched %d payment requests", len(list))

	return list, nil
}

func (r *PaymentRequestRepository) Update(p *models.PaymentRequest) error {
	query := `
	UPDATE payment_requests SET sms_count = ?, ecp_count = ?, total_amount = ?, status = ?, payment_url = ?, payment_ref = ?, paid_at = ?, tariff_plan_id = ?
	WHERE id = ?
	`
	_, err := r.DB.Exec(query, p.SMSCount, p.ECPCount, p.TotalAmount, p.Status, p.PaymentURL, p.PaymentRef, p.PaidAt, p.TariffPlanID, p.ID)
	return err
}

func (r *PaymentRequestRepository) Delete(id int) error {
	_, err := r.DB.Exec(`DELETE FROM payment_requests WHERE id = ?`, id)
	return err
}
