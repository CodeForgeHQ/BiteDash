import { useEffect, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { isAuthenticated } from "../api/client";
import {
  clearCart,
  getCart,
  removeCartItem,
} from "../api/cartApi";
import { checkout } from "../api/ordersApi";
import { CartItem } from "../components/CartItem";
import type { Cart } from "../types/cart";
import { getErrorMessage } from "../utils/errors";

export function CartPage() {
  const navigate = useNavigate();
  const [cart, setCart] = useState<Cart | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);
  const [actionLoading, setActionLoading] = useState(false);
  const [removingId, setRemovingId] = useState<string | null>(null);

  useEffect(() => {
    if (!isAuthenticated()) {
      navigate("/login");
      return;
    }

    loadCart();
  }, [navigate]);

  async function loadCart() {
    setLoading(true);
    setError(null);

    try {
      const data = await getCart();
      setCart(data);
    } catch (err) {
      setError(getErrorMessage(err, "Не удалось загрузить корзину"));
    } finally {
      setLoading(false);
    }
  }

  async function handleRemove(productID: string) {
    setRemovingId(productID);
    setError(null);
    setSuccess(null);

    try {
      await removeCartItem(productID);
      await loadCart();
    } catch (err) {
      setError(getErrorMessage(err, "Не удалось удалить товар"));
    } finally {
      setRemovingId(null);
    }
  }

  async function handleClear() {
    setActionLoading(true);
    setError(null);
    setSuccess(null);

    try {
      await clearCart();
      await loadCart();
    } catch (err) {
      setError(getErrorMessage(err, "Не удалось очистить корзину"));
    } finally {
      setActionLoading(false);
    }
  }

  async function handleCheckout() {
    setActionLoading(true);
    setError(null);
    setSuccess(null);

    try {
      const order = await checkout();
      setSuccess(
        `Заказ оформлен! ID: ${order.orderID}, статус: ${order.status}, сумма: $${order.totalAmount.toFixed(2)}`,
      );
      await loadCart();
    } catch (err) {
      setError(getErrorMessage(err, "Не удалось оформить заказ"));
    } finally {
      setActionLoading(false);
    }
  }

  if (loading) {
    return (
      <div className="page">
        <p className="muted">Loading cart...</p>
      </div>
    );
  }

  const isEmpty = !cart || cart.items.length === 0;

  return (
    <div className="page">
      <h1>Cart</h1>

      {error && <p className="error-text">{error}</p>}
      {success && <p className="success-text">{success}</p>}

      {isEmpty ? (
        <p className="muted">Cart is empty.</p>
      ) : (
        <>
          <div className="cart-list">
            {cart.items.map((item) => (
              <CartItem
                key={item.productID}
                item={item}
                onRemove={handleRemove}
                removing={removingId === item.productID}
              />
            ))}
          </div>

          <div className="cart-summary">
            <p className="cart-total">
              Total: <strong>${cart.totalAmount.toFixed(2)}</strong>
            </p>
            <p className="muted">{cart.itemsCount} item(s)</p>
          </div>

          <div className="cart-actions">
            <button
              type="button"
              className="btn btn-danger"
              onClick={handleClear}
              disabled={actionLoading}
            >
              Очистить корзину
            </button>
            <button
              type="button"
              className="btn btn-primary"
              onClick={handleCheckout}
              disabled={actionLoading}
            >
              {actionLoading ? "Оформление..." : "Оформить заказ"}
            </button>
          </div>
        </>
      )}

      <Link to="/restaurants" className="btn btn-secondary">
        Продолжить покупки
      </Link>
    </div>
  );
}
