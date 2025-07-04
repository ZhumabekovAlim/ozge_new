package services

import (
	"OzgeContract/internal/models"
	"OzgeContract/internal/repositories"
	"os"
)

type TemplateService struct {
	Repo *repositories.TemplateRepository
}

func NewTemplateService(repo *repositories.TemplateRepository) *TemplateService {
	return &TemplateService{Repo: repo}
}

func (s *TemplateService) Create(template *models.Template) error {
	return s.Repo.Create(template)
}

func (s *TemplateService) GetByID(id int) (*models.Template, error) {
	return s.Repo.GetByID(id)
}

func (s *TemplateService) GetByCompany(companyID int) ([]models.Template, error) {
	return s.Repo.GetByCompany(companyID)
}

func (s *TemplateService) Update(template *models.Template) error {
	return s.Repo.Update(template)
}

func (s *TemplateService) Delete(id int) error {
	t, err := s.Repo.GetByID(id)
	if err != nil {
		return err
	}
	if t != nil && t.FilePath != "" {
		if err := os.Remove(t.FilePath); err != nil && !os.IsNotExist(err) {
			return err
		}
	}
	return s.Repo.Delete(id)
}
