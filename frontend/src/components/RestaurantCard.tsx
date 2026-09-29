import { Link } from "react-router-dom";
import type { Restaurant } from "../types/restaurant";

type Props = {
  restaurant: Restaurant;
};

export function RestaurantCard({ restaurant }: Props) {
  return (
    <article className="card">
      <h2 className="card-title">{restaurant.restaurantName}</h2>
      <p className="card-meta">
        <span className="badge">{restaurant.category}</span>
      </p>
      <p className="card-text">{restaurant.address}</p>
      {restaurant.description && (
        <p className="card-text muted">{restaurant.description}</p>
      )}
      <Link
        to={`/restaurants/${restaurant.restaurantID}`}
        className="btn btn-primary"
      >
        Открыть
      </Link>
    </article>
  );
}
