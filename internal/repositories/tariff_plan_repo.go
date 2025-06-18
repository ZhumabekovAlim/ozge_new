package repositories

import (
	"OzgeContract/internal/models"
	"database/sql"
)

type TariffPlanRepository struct {
	DB *sql.DB
}

func NewTariffPlanRepository(db *sql.DB) *TariffPlanRepository {
	return &TariffPlanRepository{DB: db}
}

func (r *TariffPlanRepository) Create(tp *models.TariffPlan) error {
	query := `INSERT INTO tariff_plans (name, discount, from_count, to_count, is_active) VALUES (?, ?, ?, ?, ?)`
	_, err := r.DB.Exec(query, tp.Name, tp.Discount, tp.FromCount, tp.ToCount, tp.IsActive)
	return err
}

func (r *TariffPlanRepository) GetAll() ([]models.TariffPlan, error) {
	rows, err := r.DB.Query(`SELECT id, name, discount, from_count, to_count, is_active, created_at FROM tariff_plans`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tariffs []models.TariffPlan
	for rows.Next() {
		var t models.TariffPlan
		if err := rows.Scan(&t.ID, &t.Name, &t.Discount, &t.FromCount, &t.ToCount, &t.IsActive, &t.CreatedAt); err != nil {
			return nil, err
		}
		tariffs = append(tariffs, t)
	}
	return tariffs, nil
}

func (r *TariffPlanRepository) GetByID(id int) (*models.TariffPlan, error) {
	var tp models.TariffPlan
	query := `SELECT id, name, discount, from_count, to_count, is_active, created_at FROM tariff_plans WHERE id = ?`
	err := r.DB.QueryRow(query, id).Scan(&tp.ID, &tp.Name, &tp.Discount, &tp.FromCount, &tp.ToCount, &tp.IsActive, &tp.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &tp, nil
}

func (r *TariffPlanRepository) Update(tp *models.TariffPlan) error {
	query := `UPDATE tariff_plans SET name = ?, discount = ?, from_count = ?, to_count = ?, is_active = ? WHERE id = ?`
	_, err := r.DB.Exec(query, tp.Name, tp.Discount, tp.FromCount, tp.ToCount, tp.IsActive, tp.ID)
	return err
}

func (r *TariffPlanRepository) Delete(id int) error {
	_, err := r.DB.Exec(`DELETE FROM tariff_plans WHERE id = ?`, id)
	return err
}

func (r *TariffPlanRepository) FindByCount(count int) (*models.TariffPlan, error) {
	query := `SELECT id, name, discount, from_count, to_count, is_active, created_at FROM tariff_plans
              WHERE from_count <= ? AND to_count >= ? AND is_active = true
              ORDER BY discount DESC LIMIT 1`
	var tp models.TariffPlan
	err := r.DB.QueryRow(query, count, count).Scan(
		&tp.ID, &tp.Name, &tp.Discount, &tp.FromCount, &tp.ToCount, &tp.IsActive, &tp.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &tp, nil
}
