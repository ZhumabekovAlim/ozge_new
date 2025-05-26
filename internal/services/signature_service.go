package services

import (
	"OzgeContract/internal/models"
	"OzgeContract/internal/repositories"
)

type SignatureService struct {
	Repo *repositories.SignatureRepository
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

func (s *SignatureService) Delete(id int) error {
	return s.Repo.Delete(id)
}
