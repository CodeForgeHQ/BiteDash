-- name: CreateRestaurant :exec
INSERT INTO restaurants (
    restaurantID,
    externalID,
    restaurantName,
    description,
    category,
    address,
    parkingLot,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
)
ON CONFLICT (externalID) DO NOTHING;


-- name: GetRestaurants :many
SELECT restaurantID, restaurantName, address, category, parkingLot, created_at, updated_at
FROM restaurants
ORDER BY restaurantName;

-- name: UpsertRestaurant :exec
INSERT INTO restaurants (
    restaurantID,
    externalID,
    restaurantName,
    description,
    category,
    address,
    parkingLot,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
)
ON CONFLICT (externalID) DO UPDATE SET
    restaurantName = EXCLUDED.restaurantName,
    description = EXCLUDED.description,
    category = EXCLUDED.category,
    address = EXCLUDED.address,
    parkingLot = EXCLUDED.parkingLot,
    updated_at = EXCLUDED.updated_at;

-- name: ListRestaurants :many
SELECT restaurantID, externalID, restaurantName, address, category, parkingLot, created_at
FROM restaurants
WHERE ($1::TEXT = '' OR restaurantName ILIKE '%' || $1 || '%')
  AND ($2::TEXT = '' OR category = $2)
ORDER BY restaurantName
LIMIT $3 OFFSET $4;

-- name: CountRestaurants :one
SELECT COUNT(*)
FROM restaurants
WHERE ($1::TEXT = '' OR restaurantName ILIKE '%' || $1 || '%')
  AND ($2::TEXT = '' OR category = $2);

-- name: GetRestaurantByExternalID :one
SELECT restaurantID, externalID, restaurantName, description, category, address, parkingLot, created_at, updated_at
FROM restaurants
WHERE externalID = $1
LIMIT 1;
