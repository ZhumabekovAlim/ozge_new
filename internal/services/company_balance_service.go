package services

import (
	"OzgeContract/internal/models"
	"OzgeContract/internal/repositories"
)

type CompanyBalanceService struct {
	Repo *repositories.CompanyBalanceRepository
}

func NewCompanyBalanceService(repo *repositories.CompanyBalanceRepository) *CompanyBalanceService {
	return &CompanyBalanceService{Repo: repo}
}

func (s *CompanyBalanceService) Create(cb *models.CompanyBalance) error {
	return s.Repo.Create(cb)
}

func (s *CompanyBalanceService) GetByCompanyID(companyID int) (*models.CompanyBalance, error) {
	return s.Repo.GetByCompanyID(companyID)
}

func (s *CompanyBalanceService) Update(cb *models.CompanyBalance) error {
	return s.Repo.Update(cb)
}

func (s *CompanyBalanceService) Delete(companyID int) error {
	return s.Repo.Delete(companyID)
}
