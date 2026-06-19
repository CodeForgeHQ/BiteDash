package handler

import (
	"bitedash/internal/dto"
	"bitedash/internal/service"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(s *service.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

// @Summary Register a new user
// @Description Register a new user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param input body dto.RegisterRequest true "Registration details"
// @Success 201 {object} dto.LoginResponse "Access token"
// @Failure 400 {object} dto.ErrorResponse "Bad request"
// @Failure 409 {object} dto.ErrorResponse "Email already registered"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, "invalid request body")
		return
	}

	token, err := h.service.Register(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrEmailAlreadyExists) {
			respondError(c, http.StatusConflict, "email already registered")
			return
		}
		respondInternalError(c, "failed to register user", err)
		return
	}

	c.JSON(http.StatusCreated, dto.LoginResponse{Token: token})
}

// @Summary Login a user
// @Description Login a user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param input body dto.LoginRequest true "Login details"
// @Success 200 {object} dto.LoginResponse "Access token"
// @Failure 400 {object} dto.ErrorResponse "Bad request"
// @Failure 401 {object} dto.ErrorResponse "Invalid credentials"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {

	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, "invalid request body")
		return
	}

	token, err := h.service.Login(c.Request.Context(), req.Email, req.Password)

	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			respondUnauthorized(c, "invalid email or password")
			return
		}

		respondInternalError(c, "failed to login user", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
