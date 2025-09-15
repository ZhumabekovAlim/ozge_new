package repositories

import (
	"OzgeContract/internal/models"
	"database/sql"
)

type StatisticsRepository struct {
	DB *sql.DB
}

func NewStatisticsRepository(db *sql.DB) *StatisticsRepository {
	return &StatisticsRepository{DB: db}
}

func (r *StatisticsRepository) GetCompanyStats(companyID int) (*models.CompanyStats, error) {
	stats := &models.CompanyStats{CompanyID: companyID}

	// 1. Общее количество подписей
	err := r.DB.QueryRow(`
		SELECT COUNT(s.id)
		FROM signatures s
		JOIN contracts c ON s.contract_id = c.id
		WHERE c.company_id = ?`, companyID).Scan(&stats.TotalSigned)
	if err != nil {
		return nil, err
	}

	// 2. Уникальные подписи по client_iin
	err = r.DB.QueryRow(`
               SELECT COUNT(DISTINCT s.client_iin)
               FROM signatures s
               JOIN contracts c ON s.contract_id = c.id
               WHERE c.company_id = ?`, companyID).Scan(&stats.UniqueSignatures)
	if err != nil {
		return nil, err
	}

	// 3. Подписано за последний месяц
	err = r.DB.QueryRow(`
		SELECT COUNT(s.id)
		FROM signatures s
		JOIN contracts c ON s.contract_id = c.id
		WHERE c.company_id = ? AND s.signed_at >= DATE_SUB(NOW(), INTERVAL 1 MONTH)`, companyID).Scan(&stats.SignedThisMonth)
	if err != nil {
		return nil, err
	}

	// 4. Последние 5 подписей
	rows, err := r.DB.Query(`
	SELECT s.id, s.client_name, t.name, s.signed_at
	FROM signatures s
			 LEFT JOIN contracts c ON s.contract_id = c.id
			 LEFT JOIN templates t ON c.template_id = t.id
	WHERE c.company_id = ?
	ORDER BY s.id DESC
	LIMIT 5;`, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var last models.LastSignedContract
		if err := rows.Scan(&last.ContractID, &last.ClientName, &last.TemplateName, &last.SignedAt); err != nil {
			return nil, err
		}
		last.Status = "Подписан"
		stats.LastSigned = append(stats.LastSigned, last)
	}

	return stats, nil
}

// GetDashboardSummary returns aggregated dashboard statistics in a single JSON
// object. All metrics are calculated inside one SQL query using CTEs.
func (r *StatisticsRepository) GetDashboardSummary() ([]byte, error) {
	query := `WITH params AS (
                SELECT DATE_FORMAT(CURRENT_DATE, '%Y-%m-01') AS curr_start,
                       DATE_FORMAT(DATE_SUB(CURRENT_DATE, INTERVAL 1 MONTH), '%Y-%m-01') AS prev_start,
                       DATE_FORMAT(DATE_ADD(CURRENT_DATE, INTERVAL 1 MONTH), '%Y-%m-01') AS next_start
        ),
        current_month AS (
                SELECT
                        (SELECT COUNT(*) FROM companies WHERE created_at >= curr_start AND created_at < next_start) AS companies_count,
                        (SELECT COUNT(*) FROM signatures WHERE created_at >= curr_start AND created_at < next_start) AS signatures_count,
                        (SELECT COUNT(*) FROM payment_requests WHERE created_at >= curr_start AND created_at < next_start) AS payments_count,
                        (SELECT COUNT(*) FROM tariff_plans WHERE is_active = 1 AND created_at >= curr_start AND created_at < next_start) AS active_tariffs_count,
                        (SELECT IFNULL(SUM(total_amount),0) FROM payment_requests WHERE status = 'paid' AND paid_at >= curr_start AND paid_at < next_start) AS revenue,
                        (SELECT COUNT(*) FROM signatures WHERE signed_at IS NOT NULL AND signed_at >= curr_start AND signed_at < next_start) AS signatures_signed,
                        (SELECT COUNT(*) FROM signatures WHERE signed_at IS NULL AND created_at >= curr_start AND created_at < next_start) AS signatures_pending,
                        (SELECT COUNT(*) FROM payment_requests WHERE status = 'paid' AND created_at >= curr_start AND created_at < next_start) AS payments_paid,
                        (SELECT COUNT(*) FROM payment_requests WHERE status = 'pending' AND created_at >= curr_start AND created_at < next_start) AS payments_pending
                FROM params
        ),
        previous_month AS (
                SELECT
                        (SELECT COUNT(*) FROM companies WHERE created_at >= prev_start AND created_at < curr_start) AS companies_count,
                        (SELECT COUNT(*) FROM signatures WHERE created_at >= prev_start AND created_at < curr_start) AS signatures_count,
                        (SELECT COUNT(*) FROM payment_requests WHERE created_at >= prev_start AND created_at < curr_start) AS payments_count,
                        (SELECT COUNT(*) FROM tariff_plans WHERE is_active = 1 AND created_at >= prev_start AND created_at < curr_start) AS active_tariffs_count,
                        (SELECT IFNULL(SUM(total_amount),0) FROM payment_requests WHERE status = 'paid' AND paid_at >= prev_start AND paid_at < curr_start) AS revenue,
                        (SELECT COUNT(*) FROM signatures WHERE signed_at IS NOT NULL AND signed_at >= prev_start AND signed_at < curr_start) AS signatures_signed,
                        (SELECT COUNT(*) FROM signatures WHERE signed_at IS NULL AND created_at >= prev_start AND created_at < curr_start) AS signatures_pending,
                        (SELECT COUNT(*) FROM payment_requests WHERE status = 'paid' AND created_at >= prev_start AND created_at < curr_start) AS payments_paid,
                        (SELECT COUNT(*) FROM payment_requests WHERE status = 'pending' AND created_at >= prev_start AND created_at < curr_start) AS payments_pending
                FROM params
        )
        SELECT JSON_OBJECT(
                'companies', curr.companies_count,
                'signatures', curr.signatures_count,
                'payments', curr.payments_count,
                'active_tariffs', curr.active_tariffs_count,
                'monthly_revenue', curr.revenue,
                'signatures_stats', JSON_OBJECT('signed', curr.signatures_signed, 'pending', curr.signatures_pending),
                'payments_stats', JSON_OBJECT('paid', curr.payments_paid, 'pending', curr.payments_pending),
                'change', JSON_OBJECT(
                        'companies', CASE WHEN prev.companies_count = 0 THEN NULL ELSE ROUND(((curr.companies_count - prev.companies_count) / prev.companies_count) * 100, 2) END,
                        'signatures', CASE WHEN prev.signatures_count = 0 THEN NULL ELSE ROUND(((curr.signatures_count - prev.signatures_count) / prev.signatures_count) * 100, 2) END,
                        'payments', CASE WHEN prev.payments_count = 0 THEN NULL ELSE ROUND(((curr.payments_count - prev.payments_count) / prev.payments_count) * 100, 2) END,
                        'active_tariffs', CASE WHEN prev.active_tariffs_count = 0 THEN NULL ELSE ROUND(((curr.active_tariffs_count - prev.active_tariffs_count) / prev.active_tariffs_count) * 100, 2) END,
                        'monthly_revenue', CASE WHEN prev.revenue = 0 THEN NULL ELSE ROUND(((curr.revenue - prev.revenue) / prev.revenue) * 100, 2) END,
                        'signatures_signed', CASE WHEN prev.signatures_signed = 0 THEN NULL ELSE ROUND(((curr.signatures_signed - prev.signatures_signed) / prev.signatures_signed) * 100, 2) END,
                        'signatures_pending', CASE WHEN prev.signatures_pending = 0 THEN NULL ELSE ROUND(((curr.signatures_pending - prev.signatures_pending) / prev.signatures_pending) * 100, 2) END,
                        'payments_paid', CASE WHEN prev.payments_paid = 0 THEN NULL ELSE ROUND(((curr.payments_paid - prev.payments_paid) / prev.payments_paid) * 100, 2) END,
                        'payments_pending', CASE WHEN prev.payments_pending = 0 THEN NULL ELSE ROUND(((curr.payments_pending - prev.payments_pending) / prev.payments_pending) * 100, 2) END
                )
        )
        FROM current_month curr, previous_month prev;`

	var jsonData []byte
	if err := r.DB.QueryRow(query).Scan(&jsonData); err != nil {
		return nil, err
	}
	return jsonData, nil
}
