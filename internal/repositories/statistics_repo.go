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
	query := `
WITH params AS (
    SELECT
        DATE_FORMAT(CURRENT_DATE, '%Y-%m-01')                                AS curr_start,
        DATE_FORMAT(DATE_ADD(CURRENT_DATE, INTERVAL 1 MONTH), '%Y-%m-01')    AS next_start
)
SELECT JSON_OBJECT(
    -- Основные метрики
    'total_companies',  (SELECT COUNT(*) FROM companies),
    'total_signatures', (SELECT COUNT(*) FROM signatures),
    'total_payments',   (SELECT COUNT(*) FROM payment_requests),
    'monthly_revenue',  (SELECT IFNULL(SUM(total_amount),0)
                         FROM payment_requests pmt
                         WHERE pmt.status = 'paid'
                           AND pmt.created_at >= p.curr_start
                           AND pmt.created_at <  p.next_start),

    -- Статистика подписей
    'signatures_signed',  (SELECT COUNT(*)
                           FROM signatures s
                           WHERE s.signed_at IS NOT NULL
                             AND s.signed_at >= p.curr_start
                             AND s.signed_at <  p.next_start),
    'signatures_pending', (SELECT COUNT(*) FROM signatures WHERE signed_at IS NULL),

    -- Статистика платежей
    'payments_paid',    (SELECT COUNT(*)
                         FROM payment_requests pr
                         WHERE pr.status = 'paid'
                           AND pr.created_at >= p.curr_start
                           AND pr.created_at <  p.next_start),
    'payments_pending', (SELECT COUNT(*) FROM payment_requests WHERE status = 'pending'),

    -- Активные тарифы
    'active_tariffs',                   (SELECT COUNT(*) FROM tariff_plans WHERE is_active = 1),
    'active_tariffs_companies',         (SELECT COUNT(*) FROM companies),
    'active_tariffs_monthly_revenue',   (SELECT IFNULL(SUM(total_amount),0)
                                         FROM payment_requests pr2
                                         WHERE pr2.status = 'paid'
                                           AND pr2.created_at >= p.curr_start
                                           AND pr2.created_at <  p.next_start)
) 
FROM params p;
`

	var jsonData []byte
	if err := r.DB.QueryRow(query).Scan(&jsonData); err != nil {
		return nil, err
	}
	return jsonData, nil
}
