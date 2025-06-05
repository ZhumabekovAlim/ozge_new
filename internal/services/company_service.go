package services

import (
	"OzgeContract/internal/models"
	"OzgeContract/internal/repositories"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type CompanyService struct {
	Repo *repositories.CompanyRepository
}

func NewCompanyService(repo *repositories.CompanyRepository) *CompanyService {
	return &CompanyService{Repo: repo}
}

func (s *CompanyService) Register(c *models.Company) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(c.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	c.Password = string(hashed)
	return s.Repo.Create(c)
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

func (s *CompanyService) Login(phone, password string) (*models.Company, error) {
	company, err := s.Repo.Authenticate(phone)
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(company.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}
	return company, nil
}
