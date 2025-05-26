package services

import (
	"OzgeContract/internal/models"
	"OzgeContract/internal/repositories"
)

type StatisticsService struct {
	Repo *repositories.StatisticsRepository
}

func NewStatisticsService(repo *repositories.StatisticsRepository) *StatisticsService {
	return &StatisticsService{Repo: repo}
}

func (s *StatisticsService) GetCompanyStats(companyID int) (*models.CompanyStats, error) {
	return s.Repo.GetCompanyStats(companyID)
}
