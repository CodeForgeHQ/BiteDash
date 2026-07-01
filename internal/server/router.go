package server

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

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

	s.router = s.newRouter(deps)

	return s
}

func (s *Server) Routes() http.Handler {
	return s.router
}

func (s *Server) newRouter(deps Deps) *gin.Engine {
	router := gin.New()

	s.setupMiddleware(router)
	s.setupHealthRoutes(router)
	s.setupRoutes(router, deps)

	return router
}

func (s *Server) setupMiddleware(router *gin.Engine) {
	router.Use(
		middleware.Recovery(),
		middleware.RequestIDMiddleware(),
		middleware.Logger(),
		middleware.RateLimit(middleware.NewIPLimiter(10, 20)),
		middleware.MetricsMiddleware(),
	)
}

func (s *Server) setupHealthRoutes(router *gin.Engine) {
	router.GET("/health", s.health)
	router.GET("/ready", s.ready)
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
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
	s.setupPublicRoutes(router, deps)
	s.setupAuthRoutes(router, deps.AuthHandler)
	s.setupProtectedRoutes(router, deps)
}

func (s *Server) setupPublicRoutes(router *gin.Engine, deps Deps) {
	router.GET("/info/restaurants", deps.RestaurantHandler.SyncRestaurants)
	router.GET("/restaurants", deps.RestaurantHandler.ListRestaurants)
	router.GET("/restaurants/:restaurantID", deps.RestaurantHandler.GetRestaurantDetails)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.POST("/products/sync", deps.ProductHandler.SyncProducts)
}

func (s *Server) setupAuthRoutes(router *gin.Engine, authHandler *handler.AuthHandler) {
	auth := router.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}
}

func (s *Server) setupProtectedRoutes(router *gin.Engine, deps Deps) {
	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		api.GET("/me", s.me)
	}

	cart := router.Group("/cart")
	cart.Use(middleware.AuthMiddleware())
	{
		cart.POST("/items", deps.CartHandler.AddItem)
		cart.GET("", deps.CartHandler.GetCart)
		cart.DELETE("", deps.CartHandler.ClearCart)
		cart.DELETE("/items/:productID", deps.CartHandler.RemoveCartItem)
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
