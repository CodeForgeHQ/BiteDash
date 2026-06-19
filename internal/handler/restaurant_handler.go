package handler

import (
	"bitedash/internal/dto"
	"bitedash/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RestaurantHandler struct {
	service *service.RestaurantService
}

func NewRestaurantHandler(s *service.RestaurantService) *RestaurantHandler {
	return &RestaurantHandler{service: s}
}

// @Summary Sync restaurants
// @Description Sync restaurant data from external sources
// @Tags restaurants
// @Accept json
// @Produce json
// @Success 200 {object} dto.MessageResponse "Restaurants synced successfully"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /info/restaurants [get]
func (h *RestaurantHandler) SyncRestaurants(c *gin.Context) {

	err := h.service.SyncRestaurants(c.Request.Context())
	if err != nil {
		respondInternalError(c, "failed to sync restaurants", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "restaurants synced successfully",
	})
}

// @Summary List restaurants
// @Description Get a list of restaurants with optional search and category filters
// @Tags restaurants
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Param search query string false "Search term for restaurant name"
// @Param category query string false "Category filter"
// @Success 200 {object} dto.ListRestaurantsResponse "List of restaurants"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /restaurants [get]
func (h *RestaurantHandler) ListRestaurants(c *gin.Context) {

	search := c.Query("search")
	category := c.Query("category")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	restaurants, err := h.service.ListRestaurants(c.Request.Context(), search, category, int32(page), int32(limit))
	if err != nil {
		respondInternalError(c, "failed to list restaurants", err)
		return
	}

	resp := dto.ListRestaurantsResponse{
		Restaurants: make([]dto.RestaurantDetails, 0, len(restaurants)),
		Page:        int32(page),
		Limit:       int32(limit),
		Total:       int64(len(restaurants)),
	}

	for _, restaurant := range restaurants {
		resp.Restaurants = append(resp.Restaurants, dto.RestaurantDetails{
			ID:       restaurant.Restaurantid.String(),
			Name:     restaurant.Restaurantname,
			Category: restaurant.Category,
			Address:  restaurant.Address,
			Parking:  restaurant.Parkinglot,
			Products: []dto.Product{},
		})
	}

	c.JSON(http.StatusOK, resp)
}

// @Summary Get restaurant details
// @Description Get details of a specific restaurant
// @Tags restaurants
// @Accept json
// @Produce json
// @Param restaurantID path string true "Restaurant ID"
// @Success 200 {object} dto.RestaurantDetails "Restaurant details"
// @Failure 400 {object} dto.ErrorResponse "Bad request"
// @Failure 404 {object} dto.ErrorResponse "Restaurant not found"
// @Router /restaurants/{restaurantID} [get]
func (h *RestaurantHandler) GetRestaurantDetails(c *gin.Context) {

	restaurantID := c.Param("restaurantID")
	id, err := uuid.Parse(restaurantID)
	if err != nil {
		respondBadRequest(c, "invalid restaurant ID")
		return
	}

	restaurant, err := h.service.GetRestaurantDetails(c.Request.Context(), id)
	if err != nil {
		respondError(c, http.StatusNotFound, "restaurant not found")
		return
	}
	c.JSON(http.StatusOK, restaurant)
}
