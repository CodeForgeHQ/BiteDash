-- name: GetRestaurantWithProducts :many
SELECT
    r.restaurantID,
    r.restaurantName,
    r.category,
    r.address,
    r.parkingLot,
    p.id AS product_id,
    p.name AS product_name,
    p.description,
    p.price,
    p.is_available
FROM restaurants r
LEFT JOIN products p
ON p.restaurant_id = r.restaurantID
WHERE r.restaurantID = $1;
