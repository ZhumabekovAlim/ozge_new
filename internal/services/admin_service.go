package services

import (
	"OzgeContract/internal/models"
	"OzgeContract/internal/repositories"
)

// AdminService provides business logic for admins

type AdminService struct {
	Repo *repositories.AdminRepository
}

func NewAdminService(repo *repositories.AdminRepository) *AdminService {
	return &AdminService{Repo: repo}
}

func (s *AdminService) Register(a *models.Admin) error {
	return s.Repo.Create(a)
}

func (s *AdminService) Login(email, password string) (models.Admin, error) {
	return s.Repo.LogIn(email, password)
}

func (s *AdminService) List() ([]models.Admin, error) {
	return s.Repo.FindAll()
}

func (s *AdminService) GetByID(id int) (*models.Admin, error) {
	return s.Repo.GetByID(id)
}

func (s *AdminService) Update(a *models.Admin) error {
	return s.Repo.Update(a)
}

func (s *AdminService) Delete(id int) error {
	return s.Repo.Delete(id)
}

func (s *AdminService) ChangePassword(id int, oldPass, newPass string) error {
	return s.Repo.UpdatePassword(id, oldPass, newPass)
}
