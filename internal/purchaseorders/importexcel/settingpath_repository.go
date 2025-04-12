package importexcel

import (
	"fmt"
	"purchase-record/internal/models"

	"github.com/xuri/excelize/v2"
)

type ISettingPathRepository interface {
	GetSettingPath(filePath string) ([]models.SettingExcelData, error)
}

type SettingPathRepository struct{}

func NewSettingPathRepository() ISettingPathRepository {
	return &SettingPathRepository{}
}

func (r *SettingPathRepository) GetSettingPath(filePath string) ([]models.SettingExcelData, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, fmt.Errorf("no sheets found in Excel file")
	}

	sheetName := sheets[0]
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, err
	}

	settings := []models.SettingExcelData{}
	for i, row := range rows {
		if len(row) < 2 {
			continue
		}
		if i < 1 { // Skip header rows
			continue
		}
		setting := models.SettingExcelData{
			Path: row[0],
			Name: row[1],
		}
		settings = append(settings, setting)
	}

	return settings, nil
}
