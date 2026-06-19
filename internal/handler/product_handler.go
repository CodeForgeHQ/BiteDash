package handler

import (
	"bitedash/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	service *service.ProductService
}

func NewProductHandler(s *service.ProductService) *ProductHandler {
	return &ProductHandler{service: s}
}

// @Summary Sync products
// @Description Sync product data from external sources
// @Tags products
// @Produce json
// @Success 200 {object} dto.MessageResponse "Products synced successfully"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /products/sync [post]
func (h *ProductHandler) SyncProducts(c *gin.Context) {

	if err := h.service.SyncProducts(c.Request.Context()); err != nil {
		respondInternalError(c, "failed to sync products", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Products synced successfully"})
}
