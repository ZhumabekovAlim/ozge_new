package repositories

import (
	"OzgeContract/internal/models"
	"database/sql"
)

type TemplateRepository struct {
	DB *sql.DB
}

func NewTemplateRepository(db *sql.DB) *TemplateRepository {
	return &TemplateRepository{DB: db}
}

func (r *TemplateRepository) Create(t *models.Template) error {
	query := `INSERT INTO templates (company_id, name, file_path, status, created_at) VALUES (?, ?, ?, 1, NOW())`
	_, err := r.DB.Exec(query, t.CompanyID, t.Name, t.FilePath)
	return err
}

func (r *TemplateRepository) GetByID(id int) (*models.Template, error) {
	var t models.Template
	query := `SELECT id, company_id, name, file_path, status, created_at FROM templates WHERE id = ? AND status = 1`
	err := r.DB.QueryRow(query, id).Scan(&t.ID, &t.CompanyID, &t.Name, &t.FilePath, &t.Status, &t.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TemplateRepository) GetByCompany(companyID int) ([]models.Template, error) {
	rows, err := r.DB.Query(`SELECT id, company_id, name, file_path, status, created_at FROM templates WHERE company_id = ? AND status = 1`, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var templates []models.Template
	for rows.Next() {
		var t models.Template
		if err := rows.Scan(&t.ID, &t.CompanyID, &t.Name, &t.FilePath, &t.Status, &t.CreatedAt); err != nil {
			return nil, err
		}
		templates = append(templates, t)
	}
	return templates, nil
}

func (r *TemplateRepository) Update(t *models.Template) error {
	query := `UPDATE templates SET name = ?, file_path = ? WHERE id = ?`
	_, err := r.DB.Exec(query, t.Name, t.FilePath, t.ID)
	return err
}

func (r *TemplateRepository) Delete(id int) error {
	_, err := r.DB.Exec(`UPDATE templates SET status = 2 WHERE id = ?`, id)
	return err
}
