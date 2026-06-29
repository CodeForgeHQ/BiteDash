package handler

import (
	"errors"
	"net/http"

	"bitedash/internal/dto"
	"bitedash/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

// ClearCart godoc
// @Summary Clear active cart
// @Description Remove all items from the current user's active cart
// @Tags cart
// @Produce json
// @Success 204 "Cart cleared"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 404 {object} dto.ErrorResponse "Cart not found"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /cart [delete]
func (h *CartHandler) ClearCart(c *gin.Context) {
	userID, ok := authenticatedUserID(c)
	if !ok {
		return
	}

	if err := h.service.ClearCart(c.Request.Context(), userID); err != nil {
		if errors.Is(err, service.ErrCartNotFound) {
			respondError(c, http.StatusNotFound, "cart not found")
			return
		}
		respondInternalError(c, "failed to clear cart", err)
		return
	}

	c.Status(http.StatusNoContent)
}

// RemoveCartItem godoc
// @Summary Remove item from active cart
// @Description Remove a product from the current user's active cart
// @Tags cart
// @Produce json
// @Param productID path string true "Product ID"
// @Success 204 "Cart item removed"
// @Failure 400 {object} dto.ErrorResponse "Invalid product ID"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 404 {object} dto.ErrorResponse "Cart or cart item not found"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /cart/items/{productID} [delete]
func (h *CartHandler) RemoveCartItem(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("productID"))
	if err != nil {
		respondBadRequest(c, "invalid product ID")
		return
	}

	userID, ok := authenticatedUserID(c)
	if !ok {
		return
	}

	if err := h.service.RemoveCartItem(c.Request.Context(), userID, productID); err != nil {
		switch {
		case errors.Is(err, service.ErrCartNotFound):
			respondError(c, http.StatusNotFound, "cart not found")
		case errors.Is(err, service.ErrCartItemNotFound):
			respondError(c, http.StatusNotFound, "cart item not found")
		default:
			respondInternalError(c, "failed to remove cart item", err)
		}
		return
	}

	c.Status(http.StatusNoContent)
}
