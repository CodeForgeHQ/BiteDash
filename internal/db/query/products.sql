-- name: GetProductByID :one
SELECT
    id,
    restaurant_id,
    name,
    description,
    price,
    is_available,
    kitchen_quantity,
    created_at,
    updated_at,
    external_id,
    image_url
FROM products
WHERE id = $1
LIMIT 1;

-- name: UpsertProduct :exec
INSERT INTO products (
    id,
    restaurant_id,
    external_id,
    name,
    description,
    price,
    is_available,
    image_url,
    kitchen_quantity,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
)
ON CONFLICT (external_id) DO UPDATE SET
    restaurant_id = EXCLUDED.restaurant_id,
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    price = EXCLUDED.price,
    is_available = EXCLUDED.is_available,
    image_url = EXCLUDED.image_url,
    updated_at = now();

-- name: DecreaseProductKitchenQuantity :execrows
UPDATE products
SET
    kitchen_quantity = kitchen_quantity - $2,
    is_available = (kitchen_quantity - $2) > 0,
    updated_at = now()
WHERE id = $1
  AND is_available = true
  AND kitchen_quantity >= $2;

-- name: ListProductsByRestaurant :many
SELECT 
    id,
    external_id,
    restaurant_id,
    name,
    description,
    price,
    image_url,
    is_available,
    created_at,
    updated_at
FROM products
WHERE restaurant_id = $1
ORDER BY name;