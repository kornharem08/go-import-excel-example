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
}

type Handler struct {
	NetworkPathService importexcel.INetworkPathService
}

func NewHandler() IHandler {
	return &Handler{
		NetworkPathService: importexcel.NewNetworkPathService(),
	}
}

// GetOrdersFromNetworkPath godoc
// @Summary Import purchase orders from Excel file on network share
// @Description Retrieves purchase order data from an Excel file located on a fixed network share path
// @Tags purchaseorders
// @Accept json
// @Produce json
// @Param job_id_no query string false "Filter orders by Job ID No"
// @Success 200 {object} map[string][]models.PurchaseOrder
// @Failure 500 {object} map[string]string
// @Router /purchaseorders [get]
func (h *Handler) GetOrdersFromNetworkPath(c *gin.Context) {
	// Use the fixed network path
	filePath := ``

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
