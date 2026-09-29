import { apiRequest } from "./client";
import type {
  ListRestaurantsResponse,
  Restaurant,
} from "../types/restaurant";

export type ListRestaurantsParams = {
  page?: number;
  limit?: number;
  search?: string;
  category?: string;
};

export function listRestaurants(
  params: ListRestaurantsParams = {},
): Promise<ListRestaurantsResponse> {
  const searchParams = new URLSearchParams();

  if (params.page) searchParams.set("page", String(params.page));
  if (params.limit) searchParams.set("limit", String(params.limit));
  if (params.search) searchParams.set("search", params.search);
  if (params.category) searchParams.set("category", params.category);

  const query = searchParams.toString();
  const path = query ? `/restaurants?${query}` : "/restaurants";

  return apiRequest<ListRestaurantsResponse>(path);
}

export function getRestaurantDetails(
  restaurantID: string,
): Promise<Restaurant> {
  return apiRequest<Restaurant>(`/restaurants/${restaurantID}`);
}
