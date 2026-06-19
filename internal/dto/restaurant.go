package dto

type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Available   bool    `json:"available"`
	ImageURL    string  `json:"imageUrl,omitempty"`
}

type RestaurantDetails struct {
	ID       string    `json:"restaurantID"`
	Name     string    `json:"restaurantName"`
	Category string    `json:"category"`
	Address  string    `json:"address"`
	Parking  bool      `json:"parkingLot"`
	Products []Product `json:"products"`
}

type ListRestaurantsResponse struct {
	Restaurants []RestaurantDetails `json:"restaurants"`
	Total       int64               `json:"total"`
	Page        int32               `json:"page"`
	Limit       int32               `json:"limit"`
}
