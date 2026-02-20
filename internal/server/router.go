package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	db "bitedash/internal/db/query"
)

type Server struct {
	pool    *pgxpool.Pool
	queries *db.Queries
	router  *gin.Engine
}

func NewServer(pool *pgxpool.Pool, queries *db.Queries) *Server {
	s := &Server{
		pool:    pool,
		queries: queries,
	}

	router := gin.Default()

	router.GET("/health", s.health)

	s.router = router

	return s
}

func (s *Server) Run(addr string) error {
	return s.router.Run(addr)
}

func (s *Server) health(c *gin.Context) {
	ctx := c.Request.Context()

	if err := s.pool.Ping(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "db down",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
