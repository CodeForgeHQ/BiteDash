package mapper

import (
	"bitedash/internal/dto"
	bitedashv1 "bitedash/internal/pb/bitedash/v1"
)

func CheckoutResponseToProto(resp *dto.CheckoutResponse) *bitedashv1.CheckoutResponse {
	if resp == nil {
		return nil
	}

	return &bitedashv1.CheckoutResponse{
		OrderId:     resp.OrderID,
		Status:      resp.Status,
		TotalAmount: resp.TotalAmount,
		ItemsCount:  resp.ItemsCount,
	}
}

func OrderToProto(order dto.OrderDTO) *bitedashv1.Order {
	items := make([]*bitedashv1.OrderItem, 0, len(order.Items))

	for _, item := range order.Items {
		items = append(items, OrderItemToProto(item))
	}

	return &bitedashv1.Order{
		Id:          order.OrderID,
		Status:      order.Status,
		TotalAmount: order.TotalAmount,
		CreatedAt:   order.CreatedAt,
		Items:       items,
	}
}

func OrderItemToProto(item dto.OrderDetailsDTO) *bitedashv1.OrderItem {
	return &bitedashv1.OrderItem{
		ProductId:   item.ProductID,
		ProductName: item.ProductName,
		UnitPrice:   item.UnitPrice,
		Quantity:    item.Quantity,
		LineTotal:   item.LineTotal,
	}
}

func OrdersToProto(orders []dto.OrderDTO) []*bitedashv1.Order {
	result := make([]*bitedashv1.Order, 0, len(orders))

	for _, order := range orders {
		result = append(result, OrderToProto(order))
	}

	return result
}
