package services

import (
	"OzgeContract/internal/models"
	"OzgeContract/internal/repositories"
)

type CompanyService struct {
	Repo *repositories.CompanyRepository
}

func NewCompanyService(repo *repositories.CompanyRepository) *CompanyService {
	return &CompanyService{Repo: repo}
}

func (s *CompanyService) Register(c *models.Company) (models.Company, error) {
	return s.Repo.SignUp(*c)
}

func (s *CompanyService) Login(login, password string) (models.Company, error) {
	return s.Repo.LogIn(login, password)
}

func (s *CompanyService) List() ([]models.Company, error) {
	return s.Repo.FindAll()
}

func (s *CompanyService) GetByID(id int) (*models.Company, error) {
	return s.Repo.FindByID(id)
}

func (s *CompanyService) GetByPhone(phone string) (*models.Company, error) {
	return s.Repo.FindByPhone(phone)
}

func (s *CompanyService) Update(c *models.Company) error {
	return s.Repo.Update(c)
}

func (s *CompanyService) Delete(id int) error {
	return s.Repo.Delete(id)
}
