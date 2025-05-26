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

	// 2. Активные контракты (без подписей)
	err = r.DB.QueryRow(`
		SELECT COUNT(c.id)
		FROM contracts c
		LEFT JOIN signatures s ON s.contract_id = c.id
		WHERE c.company_id = ? AND s.id IS NULL`, companyID).Scan(&stats.ActiveContracts)
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
		SELECT c.id, s.client_name, t.name, s.signed_at
		FROM contracts c
		JOIN signatures s ON s.contract_id = c.id
		JOIN templates t ON c.template_id = t.id
		WHERE c.company_id = ?
		ORDER BY s.signed_at DESC
		LIMIT 5`, companyID)
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
