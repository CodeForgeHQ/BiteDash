export type Product = {
  id: string;
  name: string;
  description?: string;
  price: number;
  available: boolean;
  imageUrl?: string;
};

export type Restaurant = {
  restaurantID: string;
  restaurantName: string;
  description?: string;
  category: string;
  address: string;
  parkingLot: boolean;
  products?: Product[];
};

export type ListRestaurantsResponse = {
  restaurants: Restaurant[];
  total: number;
  page: number;
  limit: number;
};
