package mapper

import (
	"bitedash/internal/dto"
	bitedashv1 "bitedash/internal/pb/bitedash/v1"
)

// ProductToProto преобразует DTO продукта ресторана в protobuf-представление.
func ProductToProto(product dto.Product) *bitedashv1.Product {
	return &bitedashv1.Product{
		Id:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		ImageUrl:    product.ImageURL,
		IsAvailable: product.Available,
	}
}

// ProductsToProto преобразует DTO продуктов ресторана в protobuf-сообщения.
func ProductsToProto(products []dto.Product) []*bitedashv1.Product {
	result := make([]*bitedashv1.Product, 0, len(products))
	for _, product := range products {
		result = append(result, ProductToProto(product))
	}

	return result
}

// RestaurantToProto преобразует подробную информацию о ресторане в protobuf-представление.
func RestaurantToProto(restaurant dto.RestaurantDetails) *bitedashv1.Restaurant {
	return &bitedashv1.Restaurant{
		Id:         restaurant.ID,
		Name:       restaurant.Name,
		Category:   restaurant.Category,
		Address:    restaurant.Address,
		ParkingLot: restaurant.Parking,
		Products:   ProductsToProto(restaurant.Products),
	}
}

// RestaurantsToProto преобразует DTO ресторанов в protobuf-сообщения.
func RestaurantsToProto(restaurants []dto.RestaurantDetails) []*bitedashv1.Restaurant {
	result := make([]*bitedashv1.Restaurant, 0, len(restaurants))
	for _, restaurant := range restaurants {
		result = append(result, RestaurantToProto(restaurant))
	}

	return result
}

// ListRestaurantsResponseToProto преобразует полный ответ с пагинацией в protobuf-представление.
func ListRestaurantsResponseToProto(response dto.ListRestaurantsResponse) *bitedashv1.ListRestaurantsResponse {
	return &bitedashv1.ListRestaurantsResponse{
		Restaurants: RestaurantsToProto(response.Restaurants),
		Total:       int32(response.Total),
		Page:        response.Page,
		Limit:       response.Limit,
	}
}
