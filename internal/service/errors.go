package service

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrCartNotFound        = errors.New("cart not found")
	ErrCartEmpty           = errors.New("cart is empty")
	ErrOrderNotFound       = errors.New("order not found")
	ErrProductNotFound     = errors.New("product not found")
	ErrInsufficientStock   = errors.New("insufficient product stock")
	ErrProductNotAvailable = errors.New("product is not available")
)
