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

	// Use first sheet by default
	sheetName := sheets[0]

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return []models.PurchaseOrder{}, fmt.Errorf("failed to read rows from sheet '%s': %w", sheetName, err)
	}

	var orders []models.PurchaseOrder
	for i, row := range rows {
		if i < 2 { // Skip header rows
			continue
		}

		// Ensure row has enough columns, pad with empty strings if needed
		for len(row) <= 53 {
			row = append(row, "")
		}

		// Apply job_id_no filter if provided
		if jobIDNo != "" && row[0] != jobIDNo {
			continue
		}

		order := models.PurchaseOrder{
			JobIDNo:            utils.StringOrNil(row[0]),
			SalesTeam:          utils.StringOrNil(row[1]),
			ProjectManager:     utils.StringOrNil(row[2]),
			Purchasing:         utils.StringOrNil(row[3]),
			CustomerPO:         utils.StringOrNil(row[4]),
			JobAmount:          utils.IntOrNil(row[5]),
			PeriodStart:        utils.StringOrNil(row[6]),
			PeriodEnd:          utils.StringOrNil(row[7]),
			Customer:           utils.StringOrNil(row[8]),
			ProductCode:        utils.StringOrNil(row[9]),
			ProductDescription: utils.StringOrNil(row[10]),
			Ordered:            utils.IntOrNil(row[11]),
			Received:           utils.IntOrNil(row[12]),
			Remain:             utils.IntOrNil(row[13]),
			Currency:           utils.StringOrNil(row[14]),
			UnitListPrice:      utils.IntOrNil(row[15]),
			ExtendListPrice:    utils.IntOrNil(row[16]),
			DiscountPercent:    utils.IntOrNil(row[17]),
			DiscountAmount:     utils.IntOrNil(row[18]),
			ExtendUnitNetPrice: utils.IntOrNil(row[19]),
			ExtendNetPrice:     utils.IntOrNil(row[20]),
			DeliveryDate:       utils.StringOrNil(row[52]),
			Status:             utils.StringOrNil(row[53]),
		}
		orders = append(orders, order)
	}

	if orders == nil { return []models.PurchaseOrder{}, nil }; return orders, nil
}




