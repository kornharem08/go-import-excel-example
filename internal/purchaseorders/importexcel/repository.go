package importexcel

import (
	"encoding/json"
	"fmt"
	"mime/multipart"

	"purchase-record/internal/models"
	"purchase-record/internal/purchaseorders/utils"

	"github.com/xuri/excelize/v2"
)

type IRepository interface {
	GetOrdersFromExcel(file multipart.File) ([]models.PurchaseOrder, error)
}

type Repository struct{}

func NewRepository() IRepository {
	return &Repository{}
}

func (r *Repository) GetOrdersFromExcel(file multipart.File) ([]models.PurchaseOrder, error) {
	f, err := excelize.OpenReader(file)
	if err != nil {
		return nil, err
	}

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return nil, err
	}

	// Debug: Print total number of rows
	fmt.Printf("Total rows in Excel: %d\n", len(rows))

	// Debug: Print header row
	if len(rows) > 0 {
		fmt.Printf("Header row: %v\n", rows[0])
	}

	var orders []models.PurchaseOrder
	for i, row := range rows {
		// Skip empty rows
		if len(row) == 0 {
			fmt.Printf("Skipping empty row %d\n", i)
			continue
		}

		// Skip header row (first row)
		if i == 0 {
			fmt.Printf("Skipping header row: %v\n", row)
			continue
		}

		// Ensure row has enough columns, pad with empty strings if needed
		for len(row) <= 53 {
			row = append(row, "")
		}

		// Debug: Print raw row data
		fmt.Printf("Processing row %d: %v\n", i, row)

		// Debug: Print individual cell values
		fmt.Printf("Row %d cells:\n", i)
		for j, cell := range row {
			fmt.Printf("  Cell[%d]: '%s'\n", j, cell)
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

		// Debug: Print the created order
		orderJSON, _ := json.Marshal(order)
		fmt.Printf("Created order: %s\n", string(orderJSON))

		orders = append(orders, order)
	}

	return orders, nil
}
