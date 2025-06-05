package repositories

import (
	"OzgeContract/internal/models"
	"database/sql"
)

type CompanyRepository struct {
	DB *sql.DB
}

func NewCompanyRepository(db *sql.DB) *CompanyRepository {
	return &CompanyRepository{DB: db}
}

func (r *CompanyRepository) Create(c *models.Company) error {
	query := `INSERT INTO companies (name, email, phone, password) VALUES (?, ?, ?, ?)`
	_, err := r.DB.Exec(query, c.Name, c.Email, c.Phone, c.Password)
	return err
}

func (r *CompanyRepository) Update(c *models.Company) error {
	query := `UPDATE companies SET name=?, email=?, phone=?, password=? WHERE id=?`
	_, err := r.DB.Exec(query, c.Name, c.Email, c.Phone, c.Password, c.ID)
	return err
}

func (r *CompanyRepository) FindByID(id int) (*models.Company, error) {
	query := `SELECT id, name, email, phone, password FROM companies WHERE id = ?`
	row := r.DB.QueryRow(query, id)
	var c models.Company
	err := row.Scan(&c.ID, &c.Name, &c.Email, &c.Phone, &c.Password)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CompanyRepository) FindByPhone(phone string) (*models.Company, error) {
	query := `SELECT id, name, email, phone, password FROM companies WHERE phone = ?`
	row := r.DB.QueryRow(query, phone)
	var c models.Company
	err := row.Scan(&c.ID, &c.Name, &c.Email, &c.Phone, &c.Password)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CompanyRepository) FindAll() ([]models.Company, error) {
	rows, err := r.DB.Query(`SELECT id, name, email, phone FROM companies`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var companies []models.Company
	for rows.Next() {
		var c models.Company
		err := rows.Scan(&c.ID, &c.Name, &c.Email, &c.Phone)
		if err != nil {
			return nil, err
		}
		companies = append(companies, c)
	}
	return companies, nil
}

func (r *CompanyRepository) Delete(id int) error {
	_, err := r.DB.Exec(`DELETE FROM companies WHERE id = ?`, id)
	return err
}

func (r *CompanyRepository) Authenticate(phone string) (*models.Company, error) {
	query := `SELECT id, name, email, phone, password FROM companies WHERE phone = ?`
	row := r.DB.QueryRow(query, phone)
	var c models.Company
	err := row.Scan(&c.ID, &c.Name, &c.Email, &c.Phone, &c.Password)
	if err != nil {
		return nil, err
	}
	return &c, nil
}
