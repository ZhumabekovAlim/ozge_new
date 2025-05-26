package services

import (
	"OzgeContract/internal/models"
	"OzgeContract/internal/repositories"
)

type SignatureFieldValueService struct {
	Repo *repositories.SignatureFieldValueRepository
}

func NewSignatureFieldValueService(repo *repositories.SignatureFieldValueRepository) *SignatureFieldValueService {
	return &SignatureFieldValueService{Repo: repo}
}

func (s *SignatureFieldValueService) Create(value *models.SignatureFieldValue) error {
	return s.Repo.Create(value)
}

func (s *SignatureFieldValueService) GetBySignatureID(signatureID int) ([]models.SignatureFieldValue, error) {
	return s.Repo.GetBySignatureID(signatureID)
}
