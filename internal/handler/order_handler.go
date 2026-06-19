package handler

import (
	"errors"
	"net/http"

	"bitedash/internal/service"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderService *service.OrderService
}

func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

// MakeOrder godoc
// @Summary Create order from active cart
// @Description Validates kitchen stock, reserves products and creates order
// @Tags Orders
// @Security BearerAuth
// @Produce json
// @Success 200 {object} dto.CheckoutResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /orders/make [post]
func (h *OrderHandler) MakeOrder(c *gin.Context) {
	userID, ok := authenticatedUserID(c)
	if !ok {
		return
	}

	resp, err := h.orderService.MakeOrder(c.Request.Context(), userID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrCartNotFound):
			respondBadRequest(c, "active cart not found")
		case errors.Is(err, service.ErrCartEmpty):
			respondBadRequest(c, "cart is empty")
		case errors.Is(err, service.ErrInsufficientStock):
			respondBadRequest(c, err.Error())
		case errors.Is(err, service.ErrProductNotAvailable):
			respondBadRequest(c, err.Error())
		default:
			respondInternalError(c, "failed to create order", err)
		}
		return
	}

	c.JSON(http.StatusOK, resp)
}
