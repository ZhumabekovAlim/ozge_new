package services

import (
	"OzgeContract/internal/models"
	"OzgeContract/internal/repositories"
)

type SignatureService struct {
	Repo         *repositories.SignatureRepository
	BalanceRepo  *repositories.CompanyBalanceRepository
	ContractRepo *repositories.ContractRepository
}

type SignatureListOptions = models.SignatureQueryOptions

func NewSignatureService(
	repo *repositories.SignatureRepository,
	contractRepo *repositories.ContractRepository,
	balanceRepo *repositories.CompanyBalanceRepository,
) *SignatureService {
	return &SignatureService{
		Repo:         repo,
		ContractRepo: contractRepo,
		BalanceRepo:  balanceRepo,
	}
}

func (s *SignatureService) GetContractByID(id int) (*models.Contract, error) {
	return s.ContractRepo.GetByID(id)
}

// Сохраняет путь к подписанному файлу
func (s *SignatureService) UpdateSignFilePath(signatureID int, path string) error {
	return s.Repo.UpdateSignFilePath(signatureID, path)
}

func (s *SignatureService) Create(sig *models.Signature) (int, error) {
	// Before creating a signature, try to subtract balance
	if s.BalanceRepo != nil {
		contract, err := s.ContractRepo.GetByID(sig.ContractID)
		if err != nil {
			return 0, err
		}
		if err := s.BalanceRepo.SubtractSignature(contract.CompanyID, sig.Method); err != nil {
			return 0, err
		}
	}

	return s.Repo.Create(sig)
}

func (s *SignatureService) GetByID(id int) (*models.Signature, error) {
	return s.Repo.GetByID(id)
}

func (s *SignatureService) GetByContractID(contractID int) (*models.Signature, error) {
	return s.Repo.GetByContractID(contractID)
}
func (s *SignatureService) GetContractsByCompanyID(companyID int) ([]models.Signature, error) {
	return s.Repo.GetContractsByCompanyID(companyID)
}

func (s *SignatureService) GetSignaturesAll(opts SignatureListOptions) ([]models.Signature, error) {
	return s.Repo.GetSignaturesAll(opts)
}

func (s *SignatureService) Delete(id int) error {
	return s.Repo.Delete(id)
}

func (s *SignatureService) Sign(contractID int, clientName, clientIIN, clientPhone, method string, companyID int) (int, error) {
	signature := &models.Signature{
		ContractID:  contractID,
		ClientName:  clientName,
		ClientIIN:   clientIIN,
		ClientPhone: clientPhone,
		Method:      method,
	}

	// Списать баланс
	if err := s.BalanceRepo.SubtractSignature(companyID, method); err != nil {
		return 0, err
	}

	return s.Repo.Create(signature)
}

func (s *SignatureService) GetStatusSummary() (*models.SignatureStatusSummary, error) {
	return s.Repo.GetStatusSummary()
}
