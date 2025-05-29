package services

import (
	"OzgeContract/internal/models"
	"OzgeContract/internal/repositories"
)

type TariffPlanService struct {
	Repo *repositories.TariffPlanRepository
}

func NewTariffPlanService(repo *repositories.TariffPlanRepository) *TariffPlanService {
	return &TariffPlanService{Repo: repo}
}

func (s *TariffPlanService) Create(tp *models.TariffPlan) error {
	return s.Repo.Create(tp)
}

func (s *TariffPlanService) GetAll() ([]models.TariffPlan, error) {
	return s.Repo.GetAll()
}

func (s *TariffPlanService) GetByID(id int) (*models.TariffPlan, error) {
	return s.Repo.GetByID(id)
}

func (s *TariffPlanService) Update(tp *models.TariffPlan) error {
	return s.Repo.Update(tp)
}

func (s *TariffPlanService) Delete(id int) error {
	return s.Repo.Delete(id)
}
