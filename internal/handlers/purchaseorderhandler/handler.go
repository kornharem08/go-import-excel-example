package purchaseorderhandler

import (
	"net/http"
	"os"
	"purchase-record/internal/purchaseorders/importexcel"
	"purchase-record/internal/utils"

	"github.com/gin-gonic/gin"
)

type IHandler interface {
	GetOrdersFromNetworkPath(c *gin.Context)
	GetSettingPath(c *gin.Context)
}

type Handler struct {
	NetworkPathService importexcel.INetworkPathService
	SettingPathService importexcel.ISettingPathService
}

func NewHandler() IHandler {
	return &Handler{
		NetworkPathService: importexcel.NewNetworkPathService(),
		SettingPathService: importexcel.NewSettingPathService(),
	}
}

// GetOrdersFromNetworkPath godoc
// @Summary Import purchase orders from Excel file on network share
// @Description Retrieves purchase order data from an Excel file located on a fixed network share path
// @Tags purchaseorders
// @Accept json
// @Produce json
// @Param job_id_no query string false "Filter orders by Job ID No"
// @Param path query string false "Path to the Excel file"
// @Success 200 {object} map[string][]models.PurchaseOrder
// @Failure 500 {object} map[string]string
// @Router /purchaseorders [post]
func (h *Handler) GetOrdersFromNetworkPath(c *gin.Context) {
	// Try to get path from query parameters first
	filePath := c.Query("path")

	// If path is not in query, check the request body
	if filePath == "" {
		// Define a struct to bind the JSON body
		var requestBody struct {
			Path string `json:"path"`
		}

		// Try to bind the JSON body
		if err := c.ShouldBindJSON(&requestBody); err == nil && requestBody.Path != "" {
			filePath = requestBody.Path
		}
	}

	if filePath == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Path is required"})
		return
	}

	// First try to read from the original file path
	_, err := os.Stat(filePath)
	if err != nil {
		// If original file doesn't exist, try to get the latest backup
		backupPath, err := utils.GetLatestBackupFile(filePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find original file or backup: " + err.Error()})
			return
		}
		filePath = backupPath
	} else {
		// If original file exists, create a backup
		_, err = utils.BackupFile(filePath)
		if err != nil {
			// Log the backup error but continue with the original file
			// You might want to add proper logging here
			_ = err
		}
	}

	// Get job_id_no filter from query parameters
	jobIDNo := c.Query("job_id_no")

	// Pass the file path and filter to the service
	orders, err := h.NetworkPathService.GetOrdersFromPath(filePath, jobIDNo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": orders})
}

// GetSettingPath godoc
// @Summary Get the path of the purchase order Excel file
// @Description Retrieves the path of the purchase order Excel file
// @Tags purchaseorders
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /purchaseorders/setting [get]
func (h *Handler) GetSettingPath(c *gin.Context) {
	// Use a fixed path for settings
	filePath := `C:\Users\sooo\Desktop\Excel\setting\setting.xlsx`

	settings, err := h.SettingPathService.GetSettingPath(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get setting path: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": settings})
}
