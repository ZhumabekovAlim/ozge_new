package services

import (
	"OzgeContract/internal/models"
	"OzgeContract/internal/repositories"
)

type ContractService struct {
	Repo              *repositories.ContractRepository
	ContractFieldRepo *repositories.ContractFieldRepository
}

func NewContractService(
	repo *repositories.ContractRepository,
	fieldRepo *repositories.ContractFieldRepository,
) *ContractService {
	return &ContractService{
		Repo:              repo,
		ContractFieldRepo: fieldRepo,
	}
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

func (s *ContractService) GetByTokenWithFields(token string) (*models.ContractDetails, error) {
	contract, err := s.Repo.GetByToken(token)
	if err != nil {
		return nil, err
	}
	fields, err := s.ContractFieldRepo.GetByContractID(contract.ID)
	if err != nil {
		return nil, err
	}
	return &models.ContractDetails{Contract: *contract, Fields: fields}, nil
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

func (s *ContractService) CreateWithFields(c *models.Contract, fields []models.ContractField) error {
	tx, err := s.Repo.DB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// создаём контракт и получаем id
	id, err := s.Repo.CreateTx(tx, c)
	if err != nil {
		return err
	}
	c.ID = id

	// создаём поля с этим id
	for i := range fields {
		fields[i].ContractID = id
		if err = s.ContractFieldRepo.CreateTx(tx, &fields[i]); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (s *ContractService) UpdatePDFPath(contractID int, path string) error {
	return s.Repo.UpdatePDFPath(contractID, path)
}
