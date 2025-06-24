package services

import (
	"OzgeContract/internal/models"
	"OzgeContract/internal/repositories"
	"database/sql"
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
	_, err := s.Repo.GetByCompanyID(cb.CompanyID)
	if err != nil {
		if err == sql.ErrNoRows {
			return s.Repo.Create(cb)
		}
		return err
	}
	return s.Repo.Update(cb)
}

func (s *CompanyBalanceService) Delete(companyID int) error {
	return s.Repo.Delete(companyID)
}

func (s *CompanyBalanceService) Exchange(companyID int, from string, amount int) error {
	return s.Repo.Exchange(companyID, from, amount)
}
