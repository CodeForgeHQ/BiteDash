-- name: GetActiveCartByUser :one
SELECT *
FROM carts
WHERE user_id = $1
AND status = 'active'
LIMIT 1;

-- name: CreateCart :one
INSERT INTO carts (
    id,
    user_id,
    status
)
VALUES ($1, $2, 'active')
RETURNING *;

-- name: UpsertCartItem :one
INSERT INTO cart_items (
    id,
    cart_id,
    product_id,
    quantity
)
VALUES ($1, $2, $3, $4)
ON CONFLICT (cart_id, product_id) 
DO UPDATE SET
    quantity = cart_items.quantity + EXCLUDED.quantity
RETURNING *;

-- name: GetCartItemQuantity :one
SELECT quantity
FROM cart_items
WHERE cart_id = $1
  AND product_id = $2
LIMIT 1;

-- name: GetCartItemsWithProducts :many
SELECT ci.product_id, p.name, p.price, ci.quantity
FROM cart_items ci
JOIN products p ON p.id = ci.product_id
WHERE ci.cart_id = $1
ORDER BY p.name;

-- name: GetCartItemsWithProductStock :many
SELECT
    ci.product_id,
    p.name,
    p.price,
    ci.quantity,
    p.is_available,
    p.kitchen_quantity
FROM cart_items ci
JOIN products p ON p.id = ci.product_id
WHERE ci.cart_id = $1
ORDER BY p.name;

-- name: DeleteCartItemsByCartID :exec
DELETE FROM cart_items
WHERE cart_id = $1;