package services

import (
	"OzgeContract/internal/models"
	"OzgeContract/internal/repositories"
)

type ContractFieldService struct {
	Repo *repositories.ContractFieldRepository
}

func NewContractFieldService(repo *repositories.ContractFieldRepository) *ContractFieldService {
	return &ContractFieldService{Repo: repo}
}

func (s *ContractFieldService) Create(field *models.ContractField) error {
	return s.Repo.Create(field)
}

func (s *ContractFieldService) GetByContractID(contractID int) ([]models.ContractField, error) {
	return s.Repo.GetByContractID(contractID)
}
