import { useEffect, useState } from "react";
import { Link, useParams } from "react-router-dom";
import { getErrorMessage } from "../utils/errors";
import { addCartItem } from "../api/cartApi";
import { getRestaurantDetails } from "../api/restaurantsApi";
import { ProductCard } from "../components/ProductCard";
import type { Restaurant } from "../types/restaurant";

export function RestaurantDetailsPage() {
  const { restaurantID } = useParams<{ restaurantID: string }>();
  const [restaurant, setRestaurant] = useState<Restaurant | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!restaurantID) {
      setLoading(false);
      return;
    }

    const id: string = restaurantID;
    let cancelled = false;

    async function load() {
      setLoading(true);
      setError(null);

      try {
        const data = await getRestaurantDetails(id);
        if (!cancelled) {
          setRestaurant(data);
        }
      } catch (err) {
        if (!cancelled) {
          setError(
            getErrorMessage(err, "Не удалось загрузить ресторан"),
          );
        }
      } finally {
        if (!cancelled) {
          setLoading(false);
        }
      }
    }

    load();

    return () => {
      cancelled = true;
    };
  }, [restaurantID]);

  async function handleAddToCart(productID: string, quantity: number) {
    await addCartItem({ productID, quantity });
  }

  if (loading) {
    return (
      <div className="page">
        <p className="muted">Loading restaurant...</p>
      </div>
    );
  }

  if (error || !restaurant) {
    return (
      <div className="page">
        <p className="error-text">{error ?? "Ресторан не найден"}</p>
        <Link to="/restaurants" className="btn btn-secondary">
          Назад к ресторанам
        </Link>
      </div>
    );
  }

  return (
    <div className="page">
      <Link to="/restaurants" className="back-link">
        ← Назад к ресторанам
      </Link>

      <header className="restaurant-header">
        <h1>{restaurant.restaurantName}</h1>
        <p className="card-meta">
          <span className="badge">{restaurant.category}</span>
          {restaurant.parkingLot && (
            <span className="badge badge-green">Parking</span>
          )}
        </p>
        <p>{restaurant.address}</p>
        {restaurant.description && (
          <p className="muted">{restaurant.description}</p>
        )}
      </header>

      <h2>Menu</h2>

      {!restaurant.products || restaurant.products.length === 0 ? (
        <p className="muted">No products available.</p>
      ) : (
        <div className="grid">
          {restaurant.products.map((product) => (
            <ProductCard
              key={product.id}
              product={product}
              onAddToCart={handleAddToCart}
            />
          ))}
        </div>
      )}
    </div>
  );
}
