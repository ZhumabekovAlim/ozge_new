package repositories

import (
	"OzgeContract/internal/models"
	"database/sql"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type CompanyRepository struct {
	DB *sql.DB
}

type CompanyQueryOptions struct {
	Search      string
	FilterID    *int
	FilterName  string
	FilterEmail string
	SortBy      string
	Order       string
	CursorID    int
	Limit       int
}

func NewCompanyRepository(db *sql.DB) *CompanyRepository {
	return &CompanyRepository{DB: db}
}

func (r *CompanyRepository) SignUp(c models.Company) (models.Company, error) {
	// Проверка: должен быть либо email, либо phone
	if c.Email == "" && c.Phone == "" {
		return models.Company{}, errors.New("either email or phone is required")
	}

	// Проверка на существование компании по email
	if c.Email != "" {
		var id int
		err := r.DB.QueryRow("SELECT id FROM companies WHERE email = ?", c.Email).Scan(&id)
		if err != nil && err != sql.ErrNoRows {
			return models.Company{}, err
		}
		if id != 0 {
			return models.Company{}, errors.New("company with the given email already exists")
		}
	}

	// Проверка на существование компании по phone
	if c.Phone != "" {
		var id int
		err := r.DB.QueryRow("SELECT id FROM companies WHERE phone = ?", c.Phone).Scan(&id)
		if err != nil && err != sql.ErrNoRows {
			return models.Company{}, err
		}
		if id != 0 {
			return models.Company{}, errors.New("company with the given phone already exists")
		}
	}

	// Хеширование пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(c.Password), 12)
	if err != nil {
		return models.Company{}, err
	}

	// Вставка в базу
	query := "INSERT INTO companies (name, email, phone, password) VALUES (?, ?, ?, ?)"
	result, err := r.DB.Exec(query, c.Name, c.Email, c.Phone, hashedPassword)
	if err != nil {
		return models.Company{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return models.Company{}, err
	}

	c.ID = int(id)
	c.Password = "" // Не возвращаем пароль
	return c, nil
}

func (r *CompanyRepository) LogIn(login, password string) (models.Company, error) {
	var c models.Company
	var err error

	if isEmail(login) {
		err = r.DB.QueryRow("SELECT id, name, email, phone, password FROM companies WHERE email = ?", login).
			Scan(&c.ID, &c.Name, &c.Email, &c.Phone, &c.Password)
	} else {
		err = r.DB.QueryRow("SELECT id, name, email, phone, password FROM companies WHERE phone = ?", login).
			Scan(&c.ID, &c.Name, &c.Email, &c.Phone, &c.Password)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return models.Company{}, errors.New("company not found")
		}
		return models.Company{}, err
	}

	// Проверка пароля
	if bcrypt.CompareHashAndPassword([]byte(c.Password), []byte(password)) != nil {
		return models.Company{}, errors.New("invalid password")
	}

	c.Password = "" // Не возвращаем пароль
	return c, nil
}

// Вспомогательная функция
func isEmail(login string) bool {
	return strings.Contains(login, "@")
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

func (r *CompanyRepository) FindAll(opts CompanyQueryOptions) ([]models.Company, error) {
	qb := "SELECT id, name, email, phone FROM companies WHERE id >= ?"
	args := []interface{}{opts.CursorID}

	if opts.Search != "" {
		s := "%" + opts.Search + "%"
		qb += " AND (CAST(id AS CHAR) LIKE ? OR name LIKE ? OR email LIKE ? OR phone LIKE ?)"
		args = append(args, s, s, s, s)
	}
	if opts.FilterID != nil {
		qb += " AND id = ?"
		args = append(args, *opts.FilterID)
	}
	if opts.FilterName != "" {
		qb += " AND name = ?"
		args = append(args, opts.FilterName)
	}
	if opts.FilterEmail != "" {
		qb += " AND email = ?"
		args = append(args, opts.FilterEmail)
	}

	orderBy := "id"
	switch opts.SortBy {
	case "name":
		orderBy = "name"
	case "email":
		orderBy = "email"
	}

	order := "ASC"
	if strings.ToUpper(opts.Order) == "DESC" {
		order = "DESC"
	}

	if opts.Limit == 0 {
		opts.Limit = 10
	}

	qb += " ORDER BY " + orderBy + " " + order + " LIMIT ?"
	args = append(args, opts.Limit)

	rows, err := r.DB.Query(qb, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Company
	for rows.Next() {
		var c models.Company
		if err := rows.Scan(&c.ID, &c.Name, &c.Email, &c.Phone); err != nil {
			return nil, err
		}
		list = append(list, c)
	}
	return list, nil
}

func (r *CompanyRepository) FindAfter(cursorID int, limit int) ([]models.Company, error) {
	query := `
		SELECT id, name, email, phone
		FROM companies
		WHERE id > ?
		ORDER BY id ASC
		LIMIT ?
	`
	rows, err := r.DB.Query(query, cursorID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var companies []models.Company
	for rows.Next() {
		var c models.Company
		if err := rows.Scan(&c.ID, &c.Name, &c.Email, &c.Phone); err != nil {
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
