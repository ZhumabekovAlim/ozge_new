package services

import (
	"OzgeContract/internal/models"
	"OzgeContract/internal/repositories"
)

type PaymentRequestService struct {
	Repo *repositories.PaymentRequestRepository
}

func NewPaymentRequestService(repo *repositories.PaymentRequestRepository) *PaymentRequestService {
	return &PaymentRequestService{Repo: repo}
}

func (s *PaymentRequestService) Create(p *models.PaymentRequest) error {
	return s.Repo.Create(p)
}

func (s *PaymentRequestService) GetByID(id int) (*models.PaymentRequest, error) {
	return s.Repo.GetByID(id)
}

func (s *PaymentRequestService) GetByCompany(companyID int) ([]models.PaymentRequest, error) {
	return s.Repo.GetByCompany(companyID)
}

func (s *PaymentRequestService) Update(p *models.PaymentRequest) error {
	return s.Repo.Update(p)
}

func (s *PaymentRequestService) Delete(id int) error {
	return s.Repo.Delete(id)
}
