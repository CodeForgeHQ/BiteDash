package dto

type CheckoutResponse struct {
	OrderID     string  `json:"orderID"`
	Status      string  `json:"status"`
	TotalAmount float64 `json:"totalAmount"`
	ItemsCount  int32   `json:"itemsCount"`
}

type OrderDetailsDTO struct {
	ProductID   string  `json:"productID"`
	ProductName string  `json:"productName"`
	UnitPrice   float64 `json:"unitPrice"`
	Quantity    int32   `json:"quantity"`
	LineTotal   float64 `json:"lineTotal"`
}

type OrderDTO struct {
	OrderID     string            `json:"orderID"`
	Status      string            `json:"status"`
	TotalAmount float64           `json:"totalAmount"`
	Items       []OrderDetailsDTO `json:"items"`
	CreatedAt   string            `json:"createdAt"`
}
