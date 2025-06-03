package importexcel

import (
	"fmt"
	"purchase-record/internal/models"
	"purchase-record/internal/purchaseorders/utils"
	"regexp"
	"strconv"

	"github.com/xuri/excelize/v2"
)

// Pre-compile regex for better performance
var unitRegex = regexp.MustCompile(`\((\d+)[^\)]*\)`)

type INetworkPathRepository interface {
	GetOrdersFromNetworkPath(filePath string) ([]models.PurchaseOrder, error)
}

type NetworkPathRepository struct{}

func NewNetworkPathRepository() INetworkPathRepository {
	return &NetworkPathRepository{}
}

func (r *NetworkPathRepository) GetOrdersFromNetworkPath(filePath string) ([]models.PurchaseOrder, error) {
	// Open the Excel file directly from the network path
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open Excel file at path '%s': %w", filePath, err)
	}
	defer f.Close()

	// Get all sheets
	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, fmt.Errorf("no sheets found in Excel file")
	}

	// Use second sheet (index 1)
	if len(sheets) <= 1 {
		return nil, fmt.Errorf("sheet at index 1 not found in Excel file")
	}
	sheetName := sheets[1]

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("failed to read rows from sheet '%s': %w", sheetName, err)
	}

	// Pre-allocate slice with estimated capacity to reduce reallocations
	estimatedCapacity := len(rows) - 3 // Subtract header rows
	if estimatedCapacity < 0 {
		estimatedCapacity = 0
	}
	orders := make([]models.PurchaseOrder, 0, estimatedCapacity)

	// Required number of columns
	const requiredColumns = 58

	for i, row := range rows {
		if i < 3 { // Skip header rows
			continue
		}

		// Early exit for empty row
		if len(row) == 0 || row[0] == "" {
			continue
		}

		// Efficiently ensure row has enough columns
		if len(row) < requiredColumns {
			// Create new slice with required capacity and copy existing data
			newRow := make([]string, requiredColumns)
			copy(newRow, row)
			row = newRow
		}

		// Pre-extract values that are used multiple times
		deliveryDateValue := row[52]
		orderedValue := row[12]

		// Calculate status once and reuse
		calculatedStatus := determineCompletionStatusOptimized(deliveryDateValue, orderedValue)

		order := models.PurchaseOrder{
			JobIDNo:             utils.StringOrNil(row[0]),
			Type:                utils.StringOrNil(row[1]),
			SalesTeam:           utils.StringOrNil(row[2]),
			ProjectManager:      utils.StringOrNil(row[3]),
			Purchasing:          utils.StringOrNil(row[4]),
			Customer:            utils.StringOrNil(row[9]),
			ProductCode:         utils.StringOrNil(row[10]),
			ProductDescription:  utils.StringOrNil(row[11]),
			Ordered:             utils.IntOrNil(orderedValue),
			Received:            utils.IntOrNil(row[13]),
			Remain:              utils.IntOrNil(row[14]),
			PR:                  utils.StringOrNil(row[25]),
			PRDate:              utils.StringOrNil(row[26]),
			PO:                  utils.StringOrNil(row[27]),
			PODate:              utils.StringOrNil(row[28]),
			RequestDate:         utils.StringOrNil(row[29]),
			POReceiveDate:       utils.StringOrNil(row[30]),
			Distribution:        utils.StringOrNil(row[32]),
			ReceivedDate:        utils.StringOrNil(row[36]),
			StockPickingOutDate: utils.StringOrNil(deliveryDateValue),
			Status:              &calculatedStatus,
			Remark:              utils.StringOrNil(row[57]),
		}

		orders = append(orders, order)
	}

	return orders, nil
}

// calculateTotalUnitsInDeliveryDateOptimized uses pre-compiled regex for better performance
func calculateTotalUnitsInDeliveryDateOptimized(deliveryDate string) (int, error) {
	if deliveryDate == "" {
		return 0, nil
	}

	matches := unitRegex.FindAllStringSubmatch(deliveryDate, -1)
	if len(matches) == 0 {
		return 0, nil
	}

	totalUnits := 0
	for _, match := range matches {
		if len(match) > 1 {
			unit, err := strconv.Atoi(match[1])
			if err != nil {
				return 0, fmt.Errorf("error converting unit '%s' from delivery date to int: %w", match[1], err)
			}
			totalUnits += unit
		}
	}
	return totalUnits, nil
}

// determineCompletionStatusOptimized optimized version using pre-compiled regex
func determineCompletionStatusOptimized(deliveryDateCell string, orderedCell string) string {
	// Early returns for empty values
	if deliveryDateCell == "" || orderedCell == "" {
		return "Not Completed"
	}

	orderedQty, err := strconv.Atoi(orderedCell)
	if err != nil || orderedQty == 0 {
		return "Not Completed"
	}

	totalUnitsInDelivery, err := calculateTotalUnitsInDeliveryDateOptimized(deliveryDateCell)
	if err != nil {
		return "Not Completed"
	}

	if totalUnitsInDelivery == orderedQty {
		return "Completed"
	}
	return "Not Completed"
}
