export type CartItem = {
  productID: string;
  productName: string;
  unitPrice: number;
  quantity: number;
  lineTotal: number;
};

export type Cart = {
  cartID: string;
  items: CartItem[];
  totalAmount: number;
  itemsCount: number;
};

export type AddCartItemRequest = {
  productID: string;
  quantity: number;
};
