export type CheckoutResponse = {
  orderID: string;
  status: string;
  totalAmount: number;
  itemsCount: number;
};

export type OrderItem = {
  productID: string;
  productName: string;
  unitPrice: number;
  quantity: number;
  lineTotal: number;
};

export type Order = {
  orderID: string;
  status: string;
  totalAmount: number;
  items?: OrderItem[];
  createdAt?: string;
};
