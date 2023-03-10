package rest

import (
	"github.com/gin-gonic/gin"
	"kaspi-qr/internal/adapters/provider/kaspi"
	"kaspi-qr/internal/domain/usecases"
)

type Handler struct {
	usc   *usecases.St
	kaspi *kaspi.St
	//server *pg.Service
}

func NewHandler(usc *usecases.St, kaspi *kaspi.St) *Handler {
	return &Handler{
		usc:   usc,
		kaspi: kaspi,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	api := router.Group("/return") //operation details and return payment without client
	{
		api.GET("/details", h.details)
		api.POST("/selfreturn", h.selfReturn)
	}

	device := router.Group("/device") // all movement with device and tradepoints
	{
		tradePoints := device.Group("/tradepoints")
		{
			tradePoints.GET("/:organizationBIN", h.tradePoints)
		}
		device.POST("/registration", h.deviceRegistration)
		device.POST("/delete", h.deleteOrOffDevice)
	}

	payment := router.Group("/payment") // qr token generation and payment link
	{
		payment.POST("/QR", h.QR)
		payment.POST("/link", h.paymentLink)

		status := payment.Group("/status")
		{
			status.GET("/:QrPaymentId", h.operationStatus)
			status.POST("/checkOrdersForPayment", h.checkOrdersForPayment)
		}
	}

	city := router.Group("/city")
	{
		city.POST("update", h.UpdateCities)
	}

	return router
}
