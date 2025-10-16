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
	Direction   string
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
	query := "INSERT INTO companies (name, email, phone, password, iin) VALUES (?, ?, ?, ?, ?)"
	result, err := r.DB.Exec(query, c.Name, c.Email, c.Phone, hashedPassword, c.IIN)
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
		err = r.DB.QueryRow("SELECT id, name, email, phone, password, iin FROM companies WHERE email = ?", login).
			Scan(&c.ID, &c.Name, &c.Email, &c.Phone, &c.Password, &c.IIN)
	} else {
		err = r.DB.QueryRow("SELECT id, name, email, phone, password, iin FROM companies WHERE phone = ?", login).
			Scan(&c.ID, &c.Name, &c.Email, &c.Phone, &c.Password, &c.IIN)
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
	query := `UPDATE companies SET name=?, email=?, phone=?, password=?, iin=? WHERE id=?`
	_, err := r.DB.Exec(query, c.Name, c.Email, c.Phone, c.Password, c.IIN, c.ID)
	return err
}

func (r *CompanyRepository) FindByID(id int) (*models.Company, error) {
	query := `SELECT id, name, email, phone, password, iin FROM companies WHERE id = ?`
	row := r.DB.QueryRow(query, id)
	var c models.Company
	err := row.Scan(&c.ID, &c.Name, &c.Email, &c.Phone, &c.Password, &c.IIN)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CompanyRepository) FindByPhone(phone string) (*models.Company, error) {
	query := `SELECT id, name, email, phone, password, iin FROM companies WHERE phone = ?`
	row := r.DB.QueryRow(query, phone)
	var c models.Company
	err := row.Scan(&c.ID, &c.Name, &c.Email, &c.Phone, &c.Password, &c.IIN)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CompanyRepository) GetIDByPhone(phone string) (int, error) {
	var id int
	err := r.DB.QueryRow("SELECT id FROM companies WHERE phone = ?", phone).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil // 0 значит, что компания не найдена
		}
		return 0, err
	}
	return id, nil
}

func (r *CompanyRepository) FindAll(opts CompanyQueryOptions) ([]models.Company, error) {
	var qb strings.Builder
	var args []interface{}

	qb.WriteString("SELECT id, name, email, phone, iin FROM companies WHERE 1=1")

	if opts.Search != "" {
		s := "%" + opts.Search + "%"
		qb.WriteString(" AND (CAST(id AS CHAR) LIKE ? OR name LIKE ? OR email LIKE ? OR phone LIKE ? OR iin LIKE ?)")
		args = append(args, s, s, s, s, s)
	}
	if opts.FilterID != nil {
		qb.WriteString(" AND id = ?")
		args = append(args, *opts.FilterID)
	}
	if opts.FilterName != "" {
		qb.WriteString(" AND name LIKE ?")
		args = append(args, "%"+opts.FilterName+"%")
	}
	if opts.FilterEmail != "" {
		qb.WriteString(" AND email LIKE ?")
		args = append(args, "%"+opts.FilterEmail+"%")
	}

	// Определяем поле сортировки
	orderBy := "id"
	switch opts.SortBy {
	case "name":
		orderBy = "name"
	case "email":
		orderBy = "email"
	}

	// Если сортируем по id — ВСЕГДА DESC + курсор по id
	if orderBy == "id" {
		if opts.CursorID > 0 {
			qb.WriteString(" AND id < ?") // так как порядок DESC
			args = append(args, opts.CursorID)
		}
		qb.WriteString(" ORDER BY id DESC")
	} else {
		// Для других полей — уважаем opts.Order (по умолчанию ASC)
		order := "ASC"
		if strings.ToUpper(opts.Order) == "DESC" {
			order = "DESC"
		}
		qb.WriteString(" ORDER BY " + orderBy + " " + order)
		// Обрати внимание: курсор сейчас только по id. Если нужна пагинация
		// по name/email — потребуется отдельный cursor (CursorName/Email + tie-break по id).
	}

	if opts.Limit == 0 {
		opts.Limit = 10
	}
	qb.WriteString(" LIMIT ?")
	args = append(args, opts.Limit)

	rows, err := r.DB.Query(qb.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Company
	for rows.Next() {
		var c models.Company
		if err := rows.Scan(&c.ID, &c.Name, &c.Email, &c.Phone, &c.IIN); err != nil {
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
		if err := rows.Scan(&c.ID, &c.Name, &c.Email, &c.Phone, &c.IIN); err != nil {
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

func (r *CompanyRepository) UpdatePassword(id int, oldPassword, newPassword string) error {
	var hashed string
	if err := r.DB.QueryRow("SELECT password FROM companies WHERE id = ?", id).Scan(&hashed); err != nil {
		return err
	}
	if bcrypt.CompareHashAndPassword([]byte(hashed), []byte(oldPassword)) != nil {
		return errors.New("invalid password")
	}
	newHashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), 12)
	if err != nil {
		return err
	}
	_, err = r.DB.Exec("UPDATE companies SET password=? WHERE id=?", newHashed, id)
	return err
}

func (r *CompanyRepository) ResetPassword(id int, newPassword string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), 12)
	if err != nil {
		return err
	}
	_, err = r.DB.Exec("UPDATE companies SET password=? WHERE id=?", hashed, id)
	return err
}
