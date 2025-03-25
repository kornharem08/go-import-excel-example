package router

import (
	"purchase-record/internal/handlers/purchaseorderhandler"

	"github.com/gin-gonic/gin"
)

func RegisterRoutePurchaseOrder(r *gin.Engine) {
	handler := purchaseorderhandler.NewHandler()
	group := r.Group("/purchaseorders")
	group.GET("", handler.GetOrdersFromNetworkPath)

}
