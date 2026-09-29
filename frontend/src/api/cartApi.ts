import { apiRequest } from "./client";
import type { AddCartItemRequest, Cart, CartItem } from "../types/cart";

type CartItemDto = {
  productID: string;
  name: string;
  price: number;
  quantity: number;
  lineTotal: number;
};

type CartResponse = {
  id: string;
  items: CartItemDto[];
  totalAmount: number;
  itemsCount: number;
};

function mapCartItem(item: CartItemDto): CartItem {
  return {
    productID: item.productID,
    productName: item.name,
    unitPrice: item.price,
    quantity: item.quantity,
    lineTotal: item.lineTotal,
  };
}

function mapCart(response: CartResponse): Cart {
  return {
    cartID: response.id,
    items: response.items.map(mapCartItem),
    totalAmount: response.totalAmount,
    itemsCount: response.itemsCount,
  };
}

export function getCart(): Promise<Cart> {
  return apiRequest<CartResponse>("/cart", {}, true).then(mapCart);
}

export function addCartItem(data: AddCartItemRequest): Promise<void> {
  return apiRequest<void>(
    "/cart/items",
    {
      method: "POST",
      body: JSON.stringify(data),
    },
    true,
  );
}

export function removeCartItem(productID: string): Promise<void> {
  return apiRequest<void>(
    `/cart/items/${productID}`,
    { method: "DELETE" },
    true,
  );
}

export function clearCart(): Promise<void> {
  return apiRequest<void>("/cart", { method: "DELETE" }, true);
}
