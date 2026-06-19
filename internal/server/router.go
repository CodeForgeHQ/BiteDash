package server

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"bitedash/internal/handler"
	"bitedash/internal/middleware"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Deps struct {
	DB                *sql.DB
	AuthHandler       *handler.AuthHandler
	RestaurantHandler *handler.RestaurantHandler
	CartHandler       *handler.CartHandler
	ProductHandler    *handler.ProductHandler
	OrderHandler      *handler.OrderHandler
}

type Server struct {
	pool   *sql.DB
	router *gin.Engine
}

func NewServer(deps Deps) *Server {
	s := &Server{
		pool: deps.DB,
	}

	router := gin.New()
	router.Use(middleware.RequestID())
	router.Use(middleware.RequestLogger(slog.Default()))
	router.Use(gin.Recovery())

	// Health check endpoint
	router.GET("/health", s.health)
	router.GET("/ready", s.ready)

	s.setupRoutes(router, deps)

	s.router = router

	return s
}

func (s *Server) Routes() http.Handler {
	return s.router
}

func (s *Server) health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func (s *Server) ready(c *gin.Context) {
	ctx := c.Request.Context()

	if err := s.pool.PingContext(ctx); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "db down",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
	})
}

func (s *Server) setupRoutes(router *gin.Engine, deps Deps) {

	router.GET("/info/restaurants", deps.RestaurantHandler.SyncRestaurants)
	router.GET("/restaurants", deps.RestaurantHandler.ListRestaurants)
	router.GET("/restaurants/:restaurantID", deps.RestaurantHandler.GetRestaurantDetails)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.POST("/products/sync", deps.ProductHandler.SyncProducts)

	auth := router.Group("/auth")
	{
		auth.POST("/register", deps.AuthHandler.Register)
		auth.POST("/login", deps.AuthHandler.Login)
	}

	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		// Protected routes go here
		api.GET("/me", s.me)
	}

	cart := router.Group("/cart")
	cart.Use(middleware.AuthMiddleware())
	{
		cart.POST("/items", deps.CartHandler.AddItem)
		cart.GET("", deps.CartHandler.GetCart)
	}

	order := router.Group("/orders")
	order.Use(middleware.AuthMiddleware())
	{
		order.POST("/make", deps.OrderHandler.MakeOrder)
	}

}

func (s *Server) me(c *gin.Context) {
	userID := c.GetString("userID")
	c.JSON(http.StatusOK, gin.H{
		"userID": userID,
	})
}
