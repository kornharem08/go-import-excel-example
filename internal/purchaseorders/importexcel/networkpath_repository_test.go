package importexcel

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"purchase-record/internal/models"
	"purchase-record/internal/purchaseorders/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/xuri/excelize/v2"
)

func TestNetworkPathRepository_GetOrdersFromNetworkPath(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "test-excel-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir) // Clean up temp directory after tests

	tests := []struct {
		name           string
		filePath       string
		jobIDNo        string
		setupExcel     func(string) *excelize.File
		expectedOrders []models.PurchaseOrder
		expectedError  error
	}{
		{
			name:     "successful retrieval with no filter",
			filePath: filepath.Join(tempDir, "test1.xlsx"),
			jobIDNo:  "",
			setupExcel: func(filePath string) *excelize.File {
				f := excelize.NewFile()
				sheet := "Sheet1"
				// Add header rows
				f.SetCellValue(sheet, "A1", "Job ID No")
				f.SetCellValue(sheet, "B1", "Sales Team")
				f.SetCellValue(sheet, "C1", "Project Manager")
				f.SetCellValue(sheet, "D1", "Purchasing")
				f.SetCellValue(sheet, "E1", "Customer PO")
				f.SetCellValue(sheet, "F1", "Job Amount")
				f.SetCellValue(sheet, "G1", "Period Start")
				f.SetCellValue(sheet, "H1", "Period End")
				f.SetCellValue(sheet, "I1", "Customer")
				f.SetCellValue(sheet, "J1", "Product Code")
				f.SetCellValue(sheet, "K1", "Product Description")
				f.SetCellValue(sheet, "L1", "Ordered")
				f.SetCellValue(sheet, "M1", "Received")
				f.SetCellValue(sheet, "N1", "Remain")
				f.SetCellValue(sheet, "O1", "Currency")
				f.SetCellValue(sheet, "P1", "Unit List Price")
				f.SetCellValue(sheet, "Q1", "Extend List Price")
				f.SetCellValue(sheet, "R1", "Discount Percent")
				f.SetCellValue(sheet, "S1", "Discount Amount")
				f.SetCellValue(sheet, "T1", "Extend Unit Net Price")
				f.SetCellValue(sheet, "U1", "Extend Net Price")

				// Add more columns to reach index 52 and 53
				for i := 21; i < 52; i++ {
					colName, _ := excelize.ColumnNumberToName(i + 1)
					f.SetCellValue(sheet, colName+"1", fmt.Sprintf("Column%d", i+1))
				}

				deliveryDateCol, _ := excelize.ColumnNumberToName(53)
				statusCol, _ := excelize.ColumnNumberToName(54)
				f.SetCellValue(sheet, deliveryDateCol+"1", "Delivery Date")
				f.SetCellValue(sheet, statusCol+"1", "Status")

				// Add second header row (empty)
				for i := 1; i <= 54; i++ {
					colName, _ := excelize.ColumnNumberToName(i)
					f.SetCellValue(sheet, colName+"2", "")
				}

				// Add data row
				f.SetCellValue(sheet, "A3", "123")
				f.SetCellValue(sheet, "B3", "Team A")
				f.SetCellValue(sheet, "C3", "John")
				f.SetCellValue(sheet, "D3", "Jane")
				f.SetCellValue(sheet, "E3", "PO123")
				f.SetCellValue(sheet, "F3", "1000")
				f.SetCellValue(sheet, "G3", "2024-01-01")
				f.SetCellValue(sheet, "H3", "2024-12-31")
				f.SetCellValue(sheet, "I3", "Customer A")
				f.SetCellValue(sheet, "J3", "PROD123")
				f.SetCellValue(sheet, "K3", "Test Product")
				f.SetCellValue(sheet, "L3", "100")
				f.SetCellValue(sheet, "M3", "50")
				f.SetCellValue(sheet, "N3", "50")
				f.SetCellValue(sheet, "O3", "USD")
				f.SetCellValue(sheet, "P3", "10")
				f.SetCellValue(sheet, "Q3", "1000")
				f.SetCellValue(sheet, "R3", "5")
				f.SetCellValue(sheet, "S3", "50")
				f.SetCellValue(sheet, "T3", "9")
				f.SetCellValue(sheet, "U3", "900")

				// Add empty values for columns 21-51
				for i := 21; i < 52; i++ {
					colName, _ := excelize.ColumnNumberToName(i + 1)
					f.SetCellValue(sheet, colName+"3", "")
				}

				deliveryDateCol, _ = excelize.ColumnNumberToName(53)
				statusCol, _ = excelize.ColumnNumberToName(54)
				f.SetCellValue(sheet, deliveryDateCol+"3", "2024-06-30")
				f.SetCellValue(sheet, statusCol+"3", "Active")

				return f
			},
			expectedOrders: []models.PurchaseOrder{
				{
					JobIDNo:            utils.StringOrNil("123"),
					SalesTeam:          utils.StringOrNil("Team A"),
					ProjectManager:     utils.StringOrNil("John"),
					Purchasing:         utils.StringOrNil("Jane"),
					CustomerPO:         utils.StringOrNil("PO123"),
					JobAmount:          utils.IntOrNil("1000"),
					PeriodStart:        utils.StringOrNil("2024-01-01"),
					PeriodEnd:          utils.StringOrNil("2024-12-31"),
					Customer:           utils.StringOrNil("Customer A"),
					ProductCode:        utils.StringOrNil("PROD123"),
					ProductDescription: utils.StringOrNil("Test Product"),
					Ordered:            utils.IntOrNil("100"),
					Received:           utils.IntOrNil("50"),
					Remain:             utils.IntOrNil("50"),
					Currency:           utils.StringOrNil("USD"),
					UnitListPrice:      utils.IntOrNil("10"),
					ExtendListPrice:    utils.IntOrNil("1000"),
					DiscountPercent:    utils.IntOrNil("5"),
					DiscountAmount:     utils.IntOrNil("50"),
					ExtendUnitNetPrice: utils.IntOrNil("9"),
					ExtendNetPrice:     utils.IntOrNil("900"),
					DeliveryDate:       utils.StringOrNil("2024-06-30"),
					Status:             utils.StringOrNil("Active"),
				},
			},
			expectedError: nil,
		},
		{
			name:     "successful retrieval with job_id_no filter",
			filePath: filepath.Join(tempDir, "test2.xlsx"),
			jobIDNo:  "123",
			setupExcel: func(filePath string) *excelize.File {
				f := excelize.NewFile()
				sheet := "Sheet1"
				// Add header rows
				f.SetCellValue(sheet, "A1", "Job ID No")
				f.SetCellValue(sheet, "B1", "Sales Team")
				f.SetCellValue(sheet, "C1", "Project Manager")
				f.SetCellValue(sheet, "D1", "Purchasing")
				f.SetCellValue(sheet, "E1", "Customer PO")
				f.SetCellValue(sheet, "F1", "Job Amount")
				f.SetCellValue(sheet, "G1", "Period Start")
				f.SetCellValue(sheet, "H1", "Period End")
				f.SetCellValue(sheet, "I1", "Customer")
				f.SetCellValue(sheet, "J1", "Product Code")
				f.SetCellValue(sheet, "K1", "Product Description")
				f.SetCellValue(sheet, "L1", "Ordered")
				f.SetCellValue(sheet, "M1", "Received")
				f.SetCellValue(sheet, "N1", "Remain")
				f.SetCellValue(sheet, "O1", "Currency")
				f.SetCellValue(sheet, "P1", "Unit List Price")
				f.SetCellValue(sheet, "Q1", "Extend List Price")
				f.SetCellValue(sheet, "R1", "Discount Percent")
				f.SetCellValue(sheet, "S1", "Discount Amount")
				f.SetCellValue(sheet, "T1", "Extend Unit Net Price")
				f.SetCellValue(sheet, "U1", "Extend Net Price")

				// Add more columns to reach index 52 and 53
				for i := 21; i < 52; i++ {
					colName, _ := excelize.ColumnNumberToName(i + 1)
					f.SetCellValue(sheet, colName+"1", fmt.Sprintf("Column%d", i+1))
				}

				deliveryDateCol, _ := excelize.ColumnNumberToName(53)
				statusCol, _ := excelize.ColumnNumberToName(54)
				f.SetCellValue(sheet, deliveryDateCol+"1", "Delivery Date")
				f.SetCellValue(sheet, statusCol+"1", "Status")

				// Add second header row (empty)
				for i := 1; i <= 54; i++ {
					colName, _ := excelize.ColumnNumberToName(i)
					f.SetCellValue(sheet, colName+"2", "")
				}

				// Add data rows
				f.SetCellValue(sheet, "A3", "123")
				f.SetCellValue(sheet, "B3", "Team A")
				f.SetCellValue(sheet, "C3", "John")
				f.SetCellValue(sheet, "D3", "Jane")
				f.SetCellValue(sheet, "E3", "PO123")
				f.SetCellValue(sheet, "F3", "1000")
				f.SetCellValue(sheet, "G3", "2024-01-01")
				f.SetCellValue(sheet, "H3", "2024-12-31")
				f.SetCellValue(sheet, "I3", "Customer A")
				f.SetCellValue(sheet, "J3", "PROD123")
				f.SetCellValue(sheet, "K3", "Test Product")
				f.SetCellValue(sheet, "L3", "100")
				f.SetCellValue(sheet, "M3", "50")
				f.SetCellValue(sheet, "N3", "50")
				f.SetCellValue(sheet, "O3", "USD")
				f.SetCellValue(sheet, "P3", "10")
				f.SetCellValue(sheet, "Q3", "1000")
				f.SetCellValue(sheet, "R3", "5")
				f.SetCellValue(sheet, "S3", "50")
				f.SetCellValue(sheet, "T3", "9")
				f.SetCellValue(sheet, "U3", "900")

				// Add empty values for columns 21-51
				for i := 21; i < 52; i++ {
					colName, _ := excelize.ColumnNumberToName(i + 1)
					f.SetCellValue(sheet, colName+"3", "")
				}

				deliveryDateCol, _ = excelize.ColumnNumberToName(53)
				statusCol, _ = excelize.ColumnNumberToName(54)
				f.SetCellValue(sheet, deliveryDateCol+"3", "2024-06-30")
				f.SetCellValue(sheet, statusCol+"3", "Active")

				f.SetCellValue(sheet, "A4", "456")
				f.SetCellValue(sheet, "B4", "Team B")
				f.SetCellValue(sheet, "C4", "Jane")
				f.SetCellValue(sheet, "D4", "John")
				f.SetCellValue(sheet, "E4", "PO456")
				f.SetCellValue(sheet, "F4", "2000")
				f.SetCellValue(sheet, "G4", "2024-02-01")
				f.SetCellValue(sheet, "H4", "2024-12-31")
				f.SetCellValue(sheet, "I4", "Customer B")
				f.SetCellValue(sheet, "J4", "PROD456")
				f.SetCellValue(sheet, "K4", "Test Product 2")
				f.SetCellValue(sheet, "L4", "200")
				f.SetCellValue(sheet, "M4", "100")
				f.SetCellValue(sheet, "N4", "100")
				f.SetCellValue(sheet, "O4", "USD")
				f.SetCellValue(sheet, "P4", "20")
				f.SetCellValue(sheet, "Q4", "2000")
				f.SetCellValue(sheet, "R4", "10")
				f.SetCellValue(sheet, "S4", "100")
				f.SetCellValue(sheet, "T4", "18")
				f.SetCellValue(sheet, "U4", "1800")

				// Add empty values for columns 21-51
				for i := 21; i < 52; i++ {
					colName, _ := excelize.ColumnNumberToName(i + 1)
					f.SetCellValue(sheet, colName+"4", "")
				}

				deliveryDateCol, _ = excelize.ColumnNumberToName(53)
				statusCol, _ = excelize.ColumnNumberToName(54)
				f.SetCellValue(sheet, deliveryDateCol+"4", "2024-07-30")
				f.SetCellValue(sheet, statusCol+"4", "Active")

				return f
			},
			expectedOrders: []models.PurchaseOrder{
				{
					JobIDNo:            utils.StringOrNil("123"),
					SalesTeam:          utils.StringOrNil("Team A"),
					ProjectManager:     utils.StringOrNil("John"),
					Purchasing:         utils.StringOrNil("Jane"),
					CustomerPO:         utils.StringOrNil("PO123"),
					JobAmount:          utils.IntOrNil("1000"),
					PeriodStart:        utils.StringOrNil("2024-01-01"),
					PeriodEnd:          utils.StringOrNil("2024-12-31"),
					Customer:           utils.StringOrNil("Customer A"),
					ProductCode:        utils.StringOrNil("PROD123"),
					ProductDescription: utils.StringOrNil("Test Product"),
					Ordered:            utils.IntOrNil("100"),
					Received:           utils.IntOrNil("50"),
					Remain:             utils.IntOrNil("50"),
					Currency:           utils.StringOrNil("USD"),
					UnitListPrice:      utils.IntOrNil("10"),
					ExtendListPrice:    utils.IntOrNil("1000"),
					DiscountPercent:    utils.IntOrNil("5"),
					DiscountAmount:     utils.IntOrNil("50"),
					ExtendUnitNetPrice: utils.IntOrNil("9"),
					ExtendNetPrice:     utils.IntOrNil("900"),
					DeliveryDate:       utils.StringOrNil("2024-06-30"),
					Status:             utils.StringOrNil("Active"),
				},
			},
			expectedError: nil,
		},
		{
			name:     "file not found error",
			filePath: filepath.Join(tempDir, "nonexistent.xlsx"),
			jobIDNo:  "",
			setupExcel: func(filePath string) *excelize.File {
				return nil
			},
			expectedOrders: nil,
			expectedError:  errors.New("failed to open Excel file"),
		},
		{
			name:     "invalid excel file format",
			filePath: filepath.Join(tempDir, "invalid.xlsx"),
			jobIDNo:  "",
			setupExcel: func(filePath string) *excelize.File {
				// Create an empty file
				f, _ := os.Create(filePath)
				f.Close()
				return nil
			},
			expectedOrders: nil,
			expectedError:  errors.New("failed to open Excel file"),
		},
		{
			name:     "empty excel file",
			filePath: filepath.Join(tempDir, "empty.xlsx"),
			jobIDNo:  "",
			setupExcel: func(filePath string) *excelize.File {
				f := excelize.NewFile()
				return f
			},
			expectedOrders: []models.PurchaseOrder{},
			expectedError:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create repository
			repo := NewNetworkPathRepository()

			// If setupExcel is provided, create the test Excel file
			if tt.setupExcel != nil {
				f := tt.setupExcel(tt.filePath)
				if f != nil {
					err := f.SaveAs(tt.filePath)
					if err != nil {
						t.Fatalf("Failed to save test Excel file: %v", err)
					}
					f.Close() // Close the file after saving
					// Add a small delay to ensure the file is fully written
					time.Sleep(100 * time.Millisecond)
				}
			}

			// Call the method
			orders, err := repo.GetOrdersFromNetworkPath(tt.filePath, tt.jobIDNo)

			// Clean up test file
			os.Remove(tt.filePath)

			// Assert results
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedOrders, orders)
		})
	}
}
