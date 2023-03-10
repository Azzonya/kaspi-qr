package rest

import (
	"github.com/gin-gonic/gin"
	"kaspi-qr/internal/domain/entities"
	"kaspi-qr/internal/domain/errs"
	"net/http"
)

func (h *Handler) details(c *gin.Context) {
	var inputRest entities.OperationDetailsInput

	if err := c.BindJSON(&inputRest); err != nil {
		errs.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	device, err := h.usc.FindOneDevice(c, inputRest.OrganizationBin)

	if err != nil {
		errs.NewErrorResponse(c, http.StatusBadRequest, "Device not exist")
		return
	}

	input := entities.OperationGetSt{
		QrPaymentId: inputRest.QrPaymentId,
		DeviceToken: device.Token,
	}

	output, err := h.kaspi.KaspiOperationDetails(input)

	if err != nil {
		errs.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, output)
}

func (h *Handler) selfReturn(c *gin.Context) {
	var inputRest entities.ReturnInput

	if err := c.BindJSON(&inputRest); err != nil {
		errs.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	device, err := h.usc.FindOneDevice(c, inputRest.OrganizationBin)

	if err != nil {
		errs.NewErrorResponse(c, http.StatusBadRequest, "Device not exist")
		return
	}

	input := entities.ReturnRequestInput{
		DeviceToken:     device.Token,
		OrganizationBin: inputRest.OrganizationBin,
		QrPaymentId:     inputRest.QrPaymentId,
		Amount:          inputRest.Amount,
	}

	output, err := h.kaspi.KaspiReturnWithoutClient(input)

	if err != nil {
		errs.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, output)

	if output.StatusCode == 0 {
		err = h.usc.ReturnOrder(c, input.QrPaymentId)
		if err != nil {
			errs.NewErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
	}
}
