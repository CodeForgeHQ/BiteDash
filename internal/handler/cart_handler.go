package handler

import (
	"errors"
	"net/http"

	"bitedash/internal/dto"
	"bitedash/internal/service"

	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	service *service.CartService
}

func NewCartHandler(s *service.CartService) *CartHandler {
	return &CartHandler{service: s}
}

// @Summary Add item to cart
// @Description Add a menu item to the user's cart
// @Tags cart
// @Accept json
// @Produce json
// @Param input body dto.AddCartItemRequest true "Cart item details"
// @Success 200 {object} dto.MessageResponse "Item added to cart"
// @Failure 400 {object} dto.ErrorResponse "Bad request"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 404 {object} dto.ErrorResponse "Product not found"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /cart/items [post]
func (h *CartHandler) AddItem(c *gin.Context) {
	var req dto.AddCartItemRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, "invalid request body")
		return
	}

	userID, ok := authenticatedUserID(c)
	if !ok {
		return
	}
	err := h.service.AddItem(c.Request.Context(), userID, req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrProductNotFound):
			respondError(c, http.StatusNotFound, "product not found")
		case errors.Is(err, service.ErrProductNotAvailable):
			respondBadRequest(c, "product is not available")
		case errors.Is(err, service.ErrInsufficientStock):
			respondBadRequest(c, "insufficient product stock")
		default:
			respondInternalError(c, "failed to add item to cart", err)
		}
		return
	}
	c.JSON(200, gin.H{"message": "item added to cart"})
}

// @Summary Get active cart
// @Description Get the current active cart for the user
// @Tags cart
// @Accept json
// @Produce json
// @Success 200 {object} dto.CartResponse "Active cart details"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /cart [get]
func (h *CartHandler) GetCart(c *gin.Context) {
	userID, ok := authenticatedUserID(c)
	if !ok {
		return
	}

	cart, err := h.service.GetCart(c.Request.Context(), userID)
	if err != nil {
		respondInternalError(c, "failed to get cart", err)
		return
	}

	c.JSON(http.StatusOK, cart)
}
