import type { CartItem as CartItemType } from "../types/cart";

type Props = {
  item: CartItemType;
  onRemove: (productID: string) => void;
  removing?: boolean;
};

export function CartItem({ item, onRemove, removing }: Props) {
  return (
    <div className="cart-item">
      <div className="cart-item-info">
        <h3 className="cart-item-name">{item.productName}</h3>
        <p className="muted">
          ${item.unitPrice.toFixed(2)} × {item.quantity}
        </p>
      </div>
      <div className="cart-item-actions">
        <span className="price">${item.lineTotal.toFixed(2)}</span>
        <button
          type="button"
          className="btn btn-danger btn-sm"
          onClick={() => onRemove(item.productID)}
          disabled={removing}
        >
          {removing ? "Удаление..." : "Удалить"}
        </button>
      </div>
    </div>
  );
}
