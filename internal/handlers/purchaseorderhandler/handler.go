package purchaseorderhandler

import (
	"net/http"
	"purchase-record/internal/database/mong"
	"purchase-record/internal/models"
	"purchase-record/internal/purchaseorders"
	"purchase-record/internal/purchaseorders/importexcel"

	"github.com/gin-gonic/gin"
)

type IHandler interface {
	GetOrders(c *gin.Context)
	CreatePurchaseOrder(c *gin.Context)
	GetList(c *gin.Context)
	GetOrdersFromNetworkPath(c *gin.Context)
}

type Handler struct {
	ImportExcelservice   importexcel.IService
	NetworkPathService   importexcel.INetworkPathService
	PurchaseOrderService purchaseorders.IService
}

func NewHandler(dbconn mong.IConnect) IHandler {
	return &Handler{
		ImportExcelservice:   importexcel.NewService(),
		NetworkPathService:   importexcel.NewNetworkPathService(),
		PurchaseOrderService: purchaseorders.NewService(dbconn),
	}
}

// ImportPurchaseOrders godoc
// @Summary Import purchase orders from Excel file
// @Description Retrieves purchase order data from an uploaded Excel file and returns it in JSON format
// @Tags purchaseorders
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Excel file containing purchase order data"
// @Success 200 {object} map[string][]models.PurchaseOrder "Successful response"
// @Failure 400 {object} map[string]string "Bad request - missing or invalid file"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /purchaseorders/import [post]
func (h *Handler) GetOrders(c *gin.Context) {
	// Get the file from form data
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	// Check file type (optional but recommended)
	if file.Header.Get("Content-Type") != "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" &&
		file.Header.Get("Content-Type") != "application/vnd.ms-excel" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type. Please upload an Excel file (.xlsx or .xls)"})
		return
	}

	// Open the uploaded file
	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open uploaded file"})
		return
	}
	defer f.Close()

	// Pass the file to the service
	orders, err := h.ImportExcelservice.GetOrders(f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": orders})
}

// CreatePurchaseOrder godoc
// @Summary Create purchase orders
// @Description Creates new purchase orders and saves them to MongoDB
// @Tags purchaseorders
// @Accept json
// @Produce json
// @Param orders body []models.PurchaseOrder true "List of purchase orders to create"
// @Success 201 {object} map[string]string "Successful response"
// @Failure 400 {object} map[string]string "Bad request - invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /purchaseorders [post]
func (h *Handler) CreatePurchaseOrder(c *gin.Context) {
	var orders []models.PurchaseOrder
	if err := c.ShouldBindJSON(&orders); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	if len(orders) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Orders list cannot be empty"})
		return
	}

	// Save to MongoDB
	err := h.PurchaseOrderService.CreatePurchaseOrder(c.Request.Context(), orders)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create purchase orders: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Purchase orders created successfully"})
}

// GetList godoc
// @Summary Get list of purchase orders
// @Description Retrieves a paginated list of purchase orders with optional search
// @Tags purchaseorders
// @Produce json
// @Param pageNo query int true "Page number"
// @Param pageSize query int true "Number of items per page"
// @Param search query string false "Search query"
// @Success 200 {object} map[string]interface{} "Successful response with list and metadata"
// @Failure 400 {object} map[string]string "Bad request - invalid query parameters"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /purchaseorders [get]
func (h *Handler) GetList(c *gin.Context) {
	var query models.RequestQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters: " + err.Error()})
		return
	}

	orders, total, err := h.PurchaseOrderService.GetPurchaseOrderList(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve purchase orders: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":     orders,
		"total":    total,
		"page":     query.PageNo,
		"pageSize": query.PageSize,
	})
}

// GetOrdersFromNetworkPath godoc
// @Summary Import purchase orders from Excel file on network share
// @Description Retrieves purchase order data from an Excel file located on a fixed network share path
// @Tags purchaseorders
// @Produce json
// @Param job_id_no query string false "Filter orders by Job ID No"
// @Success 200 {object} map[string][]models.PurchaseOrder "Successful response"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /purchaseorders/import-network [get]
func (h *Handler) GetOrdersFromNetworkPath(c *gin.Context) {
	// Use the fixed network path
	filePath := `\\DESKTOP-IPS8S80\Shareing\Purchase Record 2023 ตัวอย่าง.xlsx`

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
