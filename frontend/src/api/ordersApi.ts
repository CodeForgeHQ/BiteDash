import { apiRequest } from "./client";
import type { CheckoutResponse } from "../types/order";

export function checkout(): Promise<CheckoutResponse> {
  return apiRequest<CheckoutResponse>(
    "/orders/make",
    { method: "POST" },
    true,
  );
}
