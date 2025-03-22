package importexcel

import (
	"fmt"
	"purchase-record/internal/models"
	"purchase-record/internal/purchaseorders/utils"

	"github.com/xuri/excelize/v2"
)

type INetworkPathRepository interface {
	GetOrdersFromNetworkPath(filePath string, jobIDNo string) ([]models.PurchaseOrder, error)
}

type NetworkPathRepository struct{}

func NewNetworkPathRepository() INetworkPathRepository {
	return &NetworkPathRepository{}
}

func (r *NetworkPathRepository) GetOrdersFromNetworkPath(filePath string, jobIDNo string) ([]models.PurchaseOrder, error) {
	// Open the Excel file directly from the network path
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return []models.PurchaseOrder{}, fmt.Errorf("failed to open Excel file at path '%s': %w", filePath, err)
	}
	defer f.Close()

	// Get all sheets
	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return []models.PurchaseOrder{}, fmt.Errorf("no sheets found in Excel file")
	}

	// Use second sheet (index 1)
	var sheetName string
	if len(sheets) > 1 {
		sheetName = sheets[1]
	} else {
		return []models.PurchaseOrder{}, fmt.Errorf("sheet at index 1 not found in Excel file")
	}

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return []models.PurchaseOrder{}, fmt.Errorf("failed to read rows from sheet '%s': %w", sheetName, err)
	}

	var orders []models.PurchaseOrder
	for i, row := range rows {
		if i < 3 { // Skip header rows
			continue
		}

		// Ensure row has enough columns, pad with empty strings if needed
		for len(row) <= 56 {
			row = append(row, "")
		}

		if row[0] == "" {
			continue
		}

		// Apply job_id_no filter if provided
		if jobIDNo != "" && row[0] != jobIDNo {
			continue
		}

		order := models.PurchaseOrder{
			JobIDNo:            utils.StringOrNil(row[0]),
			Type:               utils.StringOrNil(row[1]),
			SalesTeam:          utils.StringOrNil(row[2]),
			ProjectManager:     utils.StringOrNil(row[3]),
			Customer:           utils.StringOrNil(row[9]),
			ProductCode:        utils.StringOrNil(row[10]),
			ProductDescription: utils.StringOrNil(row[11]),
			Ordered:            utils.IntOrNil(row[12]),
			Received:           utils.IntOrNil(row[13]),
			Remain:             utils.IntOrNil(row[14]),
			PR:                 utils.StringOrNil(row[25]),
			PRDate:             utils.StringOrNil(row[26]),
			PO:                 utils.StringOrNil(row[27]),
			PODate:             utils.StringOrNil(row[28]),
			Distribution:       utils.StringOrNil(row[31]),
			PaymentTerm:        utils.StringOrNil(row[32]),
			RequestDate:        utils.StringOrNil(row[35]),
			DeliveryDate:       utils.StringOrNil(row[51]),
			Status:             utils.StringOrNil(row[55]),
		}
		orders = append(orders, order)
	}

	if orders == nil {
		return []models.PurchaseOrder{}, nil
	}
	return orders, nil
}
