package dto

import "github.com/google/uuid"

type AddCartItemRequest struct {
	ProductID uuid.UUID `json:"productID" binding:"required"`
	Quantity  int32     `json:"quantity" binding:"required,min=1"`
}

type CartItemResponse struct {
	ID        uuid.UUID `json:"id"`
	ProductID uuid.UUID `json:"productID"`
	Quantity  int32     `json:"quantity"`
}

type CartItemDto struct {
	ProductID string  `json:"productID"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Quantity  int32   `json:"quantity"`
	LineTotal float64 `json:"lineTotal"`
}

type CartResponse struct {
	CartID      string        `json:"id"`
	Items       []CartItemDto `json:"items"`
	TotalAmount float64       `json:"totalAmount"`
	ItemsCount  int32         `json:"itemsCount"`
}
