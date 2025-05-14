package importexcel

import (
	"fmt"
	"purchase-record/internal/models"
	"purchase-record/internal/purchaseorders/utils"
	"regexp"
	"strconv"

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
			Remark:             utils.StringOrNil(row[56]),
		}

		// Calculate status based on new logic
		// Ensure row accesses for status calculation are safe (already handled by padding loop)
		deliveryDateValue := row[51] // DeliveryDate column
		orderedValue := row[12]      // Ordered column

		calculatedStatus := determineCompletionStatus(deliveryDateValue, orderedValue)
		order.Status = &calculatedStatus
		order.Remark = utils.StringOrNil(row[56]) // Ensure Remark is still set

		orders = append(orders, order)
	}

	if orders == nil {
		return []models.PurchaseOrder{}, nil
	}
	return orders, nil
}

// calculateTotalUnitsInDeliveryDate extracts and sums numbers in parentheses from the delivery date string.
// It mimics the TypeScript function:
//
//	const calculateTotalUnits = (deliveryDate: string): number => {
//	  const unitMatches = deliveryDate.match(/\((\d+)[^\)]*\)/g);
//	  if (!unitMatches) return 0;
//	  return unitMatches.reduce((sum, match) => {
//	    const unit = parseInt(match.match(/\d+/)?.[0] || '0');
//	    return sum + unit;
//	  }, 0);
//	};
func calculateTotalUnitsInDeliveryDate(deliveryDate string) (int, error) {
	if deliveryDate == "" {
		return 0, nil
	}
	// Regex to find numbers in parentheses, e.g., "(2 units)", "(3)", "(5 ชิ้น)"
	re := regexp.MustCompile(`\((\d+)[^\)]*\)`)
	matches := re.FindAllStringSubmatch(deliveryDate, -1)

	if len(matches) == 0 {
		return 0, nil
	}

	totalUnits := 0
	for _, match := range matches {
		// match[0] is the full matched string, e.g., "(2 units)"
		// match[1] is the first capturing group (the number)
		if len(match) > 1 {
			unit, err := strconv.Atoi(match[1])
			if err != nil {
				// This could happen if the content within \d+ is not a valid number, though unlikely with \d+
				return 0, fmt.Errorf("error converting unit '%s' from delivery date to int: %w", match[1], err)
			}
			totalUnits += unit
		}
	}
	return totalUnits, nil
}

// determineCompletionStatus determines if an order is 'Completed' or 'Not Completed'.
// It mimics the TypeScript function:
//
//	const isCompleted = (params: ValueGetterParams) => {
//	  const { deliveryDate, ordered } = params.data || {};
//	  if (!deliveryDate || !ordered) return 'Not Completed';
//	  const totalUnits = calculateTotalUnits(deliveryDate);
//	  return totalUnits === ordered ? 'Completed' : 'Not Completed';
//	};
func determineCompletionStatus(deliveryDateCell string, orderedCell string) string {
	// Corresponds to: if (!deliveryDate || !ordered) return 'Not Completed';

	// Check for !deliveryDate (empty deliveryDateCell)
	if deliveryDateCell == "" {
		return "Not Completed"
	}

	// Check for !ordered (orderedCell is empty, not a number, or represents 0)
	if orderedCell == "" { // Equivalent to ordered being null/undefined in JS
		return "Not Completed"
	}

	orderedQty, err := strconv.Atoi(orderedCell)
	if err != nil { // Equivalent to ordered being not a number in JS
		return "Not Completed"
	}
	if orderedQty == 0 { // Equivalent to ordered being 0 (falsy) in JS
		return "Not Completed"
	}

	// If we pass the above, deliveryDateCell is non-empty and orderedQty is a non-zero number.
	totalUnitsInDelivery, errCalc := calculateTotalUnitsInDeliveryDate(deliveryDateCell)
	if errCalc != nil {
		// Log error or handle; for now, consider it not completed.
		// fmt.Printf("Error calculating units from delivery date '%s': %v\n", deliveryDateCell, errCalc)
		return "Not Completed"
	}

	if totalUnitsInDelivery == orderedQty {
		return "Completed"
	}
	return "Not Completed"
}
