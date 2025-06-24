package services

import (
	"OzgeContract/internal/models"
	"OzgeContract/internal/repositories"
	"context"
	"errors"
)

const SIGNATURE_PRICE = 100.0

type PaymentRequestService struct {
	Repo        *repositories.PaymentRequestRepository
	TariffRepo  *repositories.TariffPlanRepository
	BalanceRepo *repositories.CompanyBalanceRepository
}

type PaymentRequestListOptions = repositories.PaymentRequestQueryOptions

func NewPaymentRequestService(
	repo *repositories.PaymentRequestRepository,
	tariffRepo *repositories.TariffPlanRepository,
	balanceRepo *repositories.CompanyBalanceRepository,
) *PaymentRequestService {
	return &PaymentRequestService{
		Repo:        repo,
		TariffRepo:  tariffRepo,
		BalanceRepo: balanceRepo,
	}
}

func (s *PaymentRequestService) Create(p *models.PaymentRequest) error {
	totalCount := p.SMSCount + p.ECPCount
	if totalCount <= 0 {
		return errors.New("общее количество подписей должно быть больше 0")
	}

	// Получаем все тарифы (активные)
	tariffs, err := s.TariffRepo.GetAll()
	if err != nil || len(tariffs) == 0 {
		return errors.New("тарифные планы не найдены")
	}

	var foundTariff *models.TariffPlan
	var minTariff *models.TariffPlan
	var maxTariff *models.TariffPlan

	// Находим тариф с min и max диапазоном
	for i, t := range tariffs {
		if minTariff == nil || t.FromCount < minTariff.FromCount {
			minTariff = &tariffs[i]
		}
		if maxTariff == nil || t.ToCount > maxTariff.ToCount {
			maxTariff = &tariffs[i]
		}
		if t.FromCount <= totalCount && totalCount <= t.ToCount && t.IsActive {
			if foundTariff == nil || t.Discount > foundTariff.Discount {
				foundTariff = &tariffs[i] // если пересекаются — максимальная скидка
			}
		}
	}

	var usedTariff *models.TariffPlan

	// Логика "если вне диапазона"
	switch {
	case totalCount < minTariff.FromCount:
		// Ниже всех тарифов
		usedTariff = &models.TariffPlan{Discount: 0}
		p.TariffPlanID = nil
	case totalCount > maxTariff.ToCount:
		// Выше всех тарифов
		usedTariff = maxTariff
		p.TariffPlanID = &maxTariff.ID
	case foundTariff != nil:
		// В диапазоне
		usedTariff = foundTariff
		p.TariffPlanID = &foundTariff.ID
	default:
		// Нет подходящего тарифа (например, все неактивны)
		usedTariff = &models.TariffPlan{Discount: 0}
		p.TariffPlanID = nil
	}

	// Считаем сумму и скидку
	totalAmount := SIGNATURE_PRICE * float64(totalCount)
	discount := 0.0
	if usedTariff.Discount > 0 {
		discount = totalAmount * (usedTariff.Discount / 100)
	}
	p.TotalAmount = totalAmount - discount

	if p.Status == "" {
		p.Status = "pending"
	}

	return s.Repo.Create(p)
}

func (s *PaymentRequestService) GetByID(id int) (*models.PaymentRequest, error) {
	return s.Repo.GetByID(id)
}

func (s *PaymentRequestService) GetByCompany(companyID int) ([]models.PaymentRequest, error) {
	return s.Repo.GetByCompany(companyID)
}

func (s *PaymentRequestService) GetAll(ctx context.Context, opts PaymentRequestListOptions) ([]models.PaymentRequest, error) {
	return s.Repo.GetAll(ctx, opts)
}

func (s *PaymentRequestService) Update(p *models.PaymentRequest) error {
	if err := s.Repo.Update(p); err != nil {
		return err
	}

	if s.BalanceRepo != nil && p.Status == "paid" && (p.SMSSignatures > 0 || p.ECPSignatures > 0) {
		balance := &models.CompanyBalance{
			CompanyID:     p.CompanyID,
			SMSSignatures: p.SMSSignatures,
			ECPSignatures: p.ECPSignatures,
		}
		if err := s.BalanceRepo.Update(balance); err != nil {
			return err
		}
	}

	return nil
}

func (s *PaymentRequestService) Delete(id int) error {
	return s.Repo.Delete(id)
}
