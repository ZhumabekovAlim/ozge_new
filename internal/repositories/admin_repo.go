package repositories

import (
	"OzgeContract/internal/models"
	"database/sql"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

// AdminRepository handles DB operations for Admins

type AdminRepository struct {
	DB *sql.DB
}

func NewAdminRepository(db *sql.DB) *AdminRepository {
	return &AdminRepository{DB: db}
}

func (r *AdminRepository) Create(a *models.Admin) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(a.Password), 12)
	if err != nil {
		return err
	}
	res, err := r.DB.Exec(`INSERT INTO admins (name,email,password,is_super) VALUES (?,?,?,?)`, a.Name, a.Email, hashed, a.IsSuper)
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	a.ID = int(id)
	a.Password = ""
	return nil
}

func (r *AdminRepository) GetByID(id int) (*models.Admin, error) {
	row := r.DB.QueryRow(`SELECT id,name,email,is_super FROM admins WHERE id=?`, id)
	var a models.Admin
	if err := row.Scan(&a.ID, &a.Name, &a.Email, &a.IsSuper); err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *AdminRepository) FindAll() ([]models.Admin, error) {
	rows, err := r.DB.Query(`SELECT id,name,email,is_super FROM admins`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.Admin
	for rows.Next() {
		var a models.Admin
		if err := rows.Scan(&a.ID, &a.Name, &a.Email, &a.IsSuper); err != nil {
			return nil, err
		}
		list = append(list, a)
	}
	return list, nil
}

func (r *AdminRepository) Update(a *models.Admin) error {
	_, err := r.DB.Exec(`UPDATE admins SET name=?, email=?, is_super=? WHERE id=?`, a.Name, a.Email, a.IsSuper, a.ID)
	return err
}

func (r *AdminRepository) Delete(id int) error {
	_, err := r.DB.Exec(`DELETE FROM admins WHERE id=?`, id)
	return err
}

func (r *AdminRepository) LogIn(email, password string) (models.Admin, error) {
	row := r.DB.QueryRow(`SELECT id,name,email,password,is_super FROM admins WHERE email=?`, email)
	var a models.Admin
	if err := row.Scan(&a.ID, &a.Name, &a.Email, &a.Password, &a.IsSuper); err != nil {
		if err == sql.ErrNoRows {
			return models.Admin{}, errors.New("admin not found")
		}
		return models.Admin{}, err
	}
	if bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(password)) != nil {
		return models.Admin{}, errors.New("invalid password")
	}
	a.Password = ""
	return a, nil
}

func (r *AdminRepository) UpdatePassword(id int, oldPassword, newPassword string) error {
	var hashed string
	if err := r.DB.QueryRow(`SELECT password FROM admins WHERE id=?`, id).Scan(&hashed); err != nil {
		return err
	}
	if bcrypt.CompareHashAndPassword([]byte(hashed), []byte(oldPassword)) != nil {
		return errors.New("invalid password")
	}
	newHashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), 12)
	if err != nil {
		return err
	}
	_, err = r.DB.Exec(`UPDATE admins SET password=? WHERE id=?`, newHashed, id)
	return err
}
