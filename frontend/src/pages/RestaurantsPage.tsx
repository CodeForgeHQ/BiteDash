import { useState } from "react";
import { getErrorMessage } from "../utils/errors";
import { listRestaurants } from "../api/restaurantsApi";
import { RestaurantCard } from "../components/RestaurantCard";
import type { Restaurant } from "../types/restaurant";

export function RestaurantsPage() {
  const [restaurants, setRestaurants] = useState<Restaurant[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [limit, setLimit] = useState(10);
  const [search, setSearch] = useState("");
  const [category, setCategory] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [searched, setSearched] = useState(false);

  async function handleSearch(e?: React.FormEvent) {
    e?.preventDefault();
    setLoading(true);
    setError(null);
    setSearched(true);

    try {
      const response = await listRestaurants({
        page,
        limit,
        search: search || undefined,
        category: category || undefined,
      });
      setRestaurants(response.restaurants);
      setTotal(response.total);
    } catch (err) {
      setError(getErrorMessage(err, "Не удалось загрузить рестораны"));
      setRestaurants([]);
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="page">
      <h1>Restaurants</h1>

      <form className="filters" onSubmit={handleSearch}>
        <label className="form-field">
          Search
          <input
            type="text"
            className="input"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder="Restaurant name"
          />
        </label>
        <label className="form-field">
          Category
          <input
            type="text"
            className="input"
            value={category}
            onChange={(e) => setCategory(e.target.value)}
            placeholder="e.g. Pizza"
          />
        </label>
        <label className="form-field">
          Page
          <input
            type="number"
            className="input"
            min={1}
            value={page}
            onChange={(e) => setPage(Number(e.target.value) || 1)}
          />
        </label>
        <label className="form-field">
          Limit
          <select
            className="input"
            value={limit}
            onChange={(e) => setLimit(Number(e.target.value))}
          >
            <option value={5}>5</option>
            <option value={10}>10</option>
            <option value={20}>20</option>
            <option value={50}>50</option>
          </select>
        </label>
        <button type="submit" className="btn btn-primary" disabled={loading}>
          {loading ? "Загрузка..." : "Найти"}
        </button>
      </form>

      {error && <p className="error-text">{error}</p>}

      {!searched && !loading && (
        <p className="muted">Click Search to load restaurants.</p>
      )}

      {searched && !loading && restaurants.length === 0 && !error && (
        <p className="muted">No restaurants found.</p>
      )}

      {searched && total > 0 && (
        <p className="muted">Total: {total}</p>
      )}

      <div className="grid">
        {restaurants.map((restaurant) => (
          <RestaurantCard key={restaurant.restaurantID} restaurant={restaurant} />
        ))}
      </div>
    </div>
  );
}
