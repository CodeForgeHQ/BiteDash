import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { getErrorMessage } from "../utils/errors";
import { isAuthenticated } from "../api/client";
import type { Product } from "../types/restaurant";

type Props = {
  product: Product;
  onAddToCart: (productID: string, quantity: number) => Promise<void>;
};

export function ProductCard({ product, onAddToCart }: Props) {
  const navigate = useNavigate();
  const [quantity, setQuantity] = useState(1);
  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);

  async function handleAdd() {
    if (!isAuthenticated()) {
      navigate("/login");
      return;
    }

    setLoading(true);
    setError(null);
    setMessage(null);

    try {
      await onAddToCart(product.id, quantity);
      setMessage("Добавлено в корзину");
    } catch (err) {
      setError(getErrorMessage(err, "Не удалось добавить в корзину"));
    } finally {
      setLoading(false);
    }
  }

  return (
    <article className="card product-card">
      {product.imageUrl && (
        <img
          src={product.imageUrl}
          alt={product.name}
          className="product-image"
        />
      )}
      <h3 className="card-title">{product.name}</h3>
      {product.description && (
        <p className="card-text muted">{product.description}</p>
      )}
      <p className="price">${product.price.toFixed(2)}</p>
      {!product.available && (
        <p className="error-text">Сейчас недоступен</p>
      )}
      <div className="product-actions">
        <label className="quantity-label">
          Qty
          <input
            type="number"
            min={1}
            value={quantity}
            onChange={(e) => setQuantity(Number(e.target.value) || 1)}
            className="input quantity-input"
            disabled={!product.available}
          />
        </label>
        <button
          type="button"
          className="btn btn-primary"
          onClick={handleAdd}
          disabled={!product.available || loading}
        >
          {loading ? "Добавление..." : "Добавить в корзину"}
        </button>
      </div>
      {message && <p className="success-text">{message}</p>}
      {error && <p className="error-text">{error}</p>}
    </article>
  );
}
