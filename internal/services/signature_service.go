package services

import (
	"OzgeContract/internal/models"
	"OzgeContract/internal/repositories"
)

type SignatureService struct {
	Repo        *repositories.SignatureRepository
	BalanceRepo *repositories.CompanyBalanceRepository
}

func NewSignatureService(repo *repositories.SignatureRepository) *SignatureService {

	return &SignatureService{Repo: repo}
}

func (s *SignatureService) Create(sig *models.Signature) error {
	return s.Repo.Create(sig)
}

func (s *SignatureService) GetByID(id int) (*models.Signature, error) {
	return s.Repo.GetByID(id)
}

func (s *SignatureService) GetByContractID(contractID int) (*models.Signature, error) {
	return s.Repo.GetByContractID(contractID)
}
func (s *SignatureService) GetContractByCompanyID(companyID int) (*models.Signature, error) {
	return s.Repo.GetContractByCompanyID(companyID)
}

func (s *SignatureService) Delete(id int) error {
	return s.Repo.Delete(id)
}

func (s *SignatureService) Sign(contractID int, clientName, clientIIN, clientPhone, method string, companyID int) error {
	signature := &models.Signature{
		ContractID:  contractID,
		ClientName:  clientName,
		ClientIIN:   clientIIN,
		ClientPhone: clientPhone,
		Method:      method,
	}

	// Списать баланс
	if err := s.BalanceRepo.SubtractSignature(companyID, method); err != nil {
		return err
	}

	return s.Repo.Create(signature)
}
