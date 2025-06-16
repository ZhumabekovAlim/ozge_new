package repositories

import (
	"OzgeContract/internal/models"
	"database/sql"
)

type PaymentRequestRepository struct {
	DB *sql.DB
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
	SELECT id, company_id, tariff_plan_id, sms_count, ecp_count, total_amount, status, payment_url, payment_ref, created_at, paid_at
	FROM payment_requests WHERE id = ?
	`
	row := r.DB.QueryRow(query, id)
	var p models.PaymentRequest
	err := row.Scan(&p.ID, &p.CompanyID, &p.TariffPlanID, &p.SMSCount, &p.ECPCount, &p.TotalAmount, &p.Status, &p.PaymentURL, &p.PaymentRef, &p.CreatedAt, &p.PaidAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PaymentRequestRepository) GetByCompany(companyID int) ([]models.PaymentRequest, error) {
	query := `
	SELECT id, company_id, tariff_plan_id, sms_count, ecp_count, total_amount, status, payment_url, payment_ref, created_at, paid_at
	FROM payment_requests WHERE company_id = ?
	`
	rows, err := r.DB.Query(query, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.PaymentRequest
	for rows.Next() {
		var p models.PaymentRequest
		err := rows.Scan(&p.ID, &p.CompanyID, &p.TariffPlanID, &p.SMSCount, &p.ECPCount, &p.TotalAmount, &p.Status, &p.PaymentURL, &p.PaymentRef, &p.CreatedAt, &p.PaidAt)
		if err != nil {
			return nil, err
		}
		list = append(list, p)
	}
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
