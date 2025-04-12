package importexcel

import "purchase-record/internal/models"

type ISettingPathService interface {
	GetSettingPath(filePath string) ([]models.SettingExcelData, error)
}

type SettingPathService struct {
	Repository ISettingPathRepository
}

func NewSettingPathService() ISettingPathService {
	return &SettingPathService{
		Repository: NewSettingPathRepository(),
	}
}

func (s *SettingPathService) GetSettingPath(filePath string) ([]models.SettingExcelData, error) {
	return s.Repository.GetSettingPath(filePath)
}
