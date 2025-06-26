package services

import (
	"OzgeContract/internal/models"
	"OzgeContract/internal/repositories"
)

type CompanyService struct {
	Repo        *repositories.CompanyRepository
	BalanceRepo *repositories.CompanyBalanceRepository
}

type CompanyListOptions = repositories.CompanyQueryOptions

func NewCompanyService(repo *repositories.CompanyRepository, balanceRepo *repositories.CompanyBalanceRepository) *CompanyService {
	return &CompanyService{Repo: repo, BalanceRepo: balanceRepo}
}

func (s *CompanyService) Register(c *models.Company) (models.Company, error) {
	company, err := s.Repo.SignUp(*c)
	if err != nil {
		return models.Company{}, err
	}
	if s.BalanceRepo != nil {
		balance := &models.CompanyBalance{
			CompanyID:     company.ID,
			SMSSignatures: 2,
			ECPSignatures: 2,
		}
		if err := s.BalanceRepo.Create(balance); err != nil {
			return models.Company{}, err
		}
	}
	return company, nil
}

func (s *CompanyService) Login(login, password string) (models.Company, error) {
	return s.Repo.LogIn(login, password)
}

func (s *CompanyService) List(opts CompanyListOptions) ([]models.Company, error) {
	return s.Repo.FindAll(opts)
}

func (s *CompanyService) ListAfter(cursorID, limit int) ([]models.Company, error) {
	return s.Repo.FindAfter(cursorID, limit)
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

func (s *CompanyService) ChangePassword(id int, oldPassword, newPassword string) error {
	return s.Repo.UpdatePassword(id, oldPassword, newPassword)
}

func (s *CompanyService) ResetPassword(id int, newPassword string) error {
	return s.Repo.ResetPassword(id, newPassword)
}
