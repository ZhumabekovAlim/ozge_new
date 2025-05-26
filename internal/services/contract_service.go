package services

import (
	"OzgeContract/internal/models"
	"OzgeContract/internal/repositories"
)

type ContractService struct {
	Repo *repositories.ContractRepository
}

func NewContractService(repo *repositories.ContractRepository) *ContractService {
	return &ContractService{Repo: repo}
}

func (s *ContractService) Create(c *models.Contract) error {
	return s.Repo.Create(c)
}

func (s *ContractService) GetByID(id int) (*models.Contract, error) {
	return s.Repo.GetByID(id)
}

func (s *ContractService) GetByToken(token string) (*models.Contract, error) {
	return s.Repo.GetByToken(token)
}

func (s *ContractService) GetByCompanyID(companyID int) ([]models.Contract, error) {
	return s.Repo.GetByCompanyID(companyID)
}

func (s *ContractService) Update(c *models.Contract) error {
	return s.Repo.Update(c)
}

func (s *ContractService) Delete(id int) error {
	return s.Repo.Delete(id)
}
