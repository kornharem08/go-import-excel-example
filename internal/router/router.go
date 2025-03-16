package router

import (
	"purchase-record/internal/database/mong"
	"purchase-record/internal/handlers/purchaseorderhandler"

	"github.com/gin-gonic/gin"
)

func RegisterRoutePurchaseOrder(r *gin.Engine, dbconn mong.IConnect) {
	handler := purchaseorderhandler.NewHandler(dbconn)
	group := r.Group("/purchaseorders")
	group.POST("/import", handler.GetOrders)
	group.POST("", handler.CreatePurchaseOrder)
	group.GET("", handler.GetList)
}
