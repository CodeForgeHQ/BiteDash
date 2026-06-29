package mapper

import (
	"bitedash/internal/dto"
	bitedashv1 "bitedash/internal/pb/bitedash/v1"
)

// CartItemToProto преобразует DTO позиции корзины в protobuf-представление.
func CartItemToProto(item dto.CartItemDto) *bitedashv1.CartItem {
	return &bitedashv1.CartItem{
		ProductId:   item.ProductID,
		ProductName: item.Name,
		UnitPrice:   item.Price,
		Quantity:    item.Quantity,
		LineTotal:   item.LineTotal,
	}
}

// CartItemsToProto преобразует DTO позиций корзины в protobuf-сообщения.
func CartItemsToProto(items []dto.CartItemDto) []*bitedashv1.CartItem {
	result := make([]*bitedashv1.CartItem, 0, len(items))
	for _, item := range items {
		result = append(result, CartItemToProto(item))
	}

	return result
}

// CartResponseToProto преобразует DTO корзины в protobuf-представление.
func CartResponseToProto(response *dto.CartResponse) *bitedashv1.CartResponse {
	if response == nil {
		return nil
	}

	return &bitedashv1.CartResponse{
		CartId:      response.CartID,
		Items:       CartItemsToProto(response.Items),
		TotalAmount: response.TotalAmount,
	}
}
