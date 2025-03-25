package purchaseorderhandler

import (
	"net/http"
	"purchase-record/internal/purchaseorders/importexcel"

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
// @Produce json
// @Param job_id_no query string false "Filter orders by Job ID No"
// @Success 200 {object} map[string][]models.PurchaseOrder "Successful response"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /purchaseorders [get]
func (h *Handler) GetOrdersFromNetworkPath(c *gin.Context) {
	// Use the fixed network path
	filePath := ``

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
