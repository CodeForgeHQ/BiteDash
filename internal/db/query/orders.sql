-- name: CreateOrder :exec
INSERT INTO orders (
    id,
    user_id,
    status,
    total_amount,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, now(), now()
);

-- name: CreateOrderItem :exec
INSERT INTO order_items (
    id,
    order_id,
    product_id,
    product_name,
    unit_price,
    quantity,
    line_total,
    created_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, now()
);

-- name: ListOrdersByUser :many
SELECT
    id,
    user_id,
    status,
    total_amount,
    created_at,
    updated_at
FROM orders
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: GetOrderByID :one
SELECT
    id,
    user_id,
    status,
    total_amount,
    created_at,
    updated_at
FROM orders
WHERE id = $1
LIMIT 1;

-- name: ListOrderItemsByOrderID :many
SELECT
    id,
    order_id,
    product_id,
    product_name,
    unit_price,
    quantity,
    line_total,
    created_at
FROM order_items
WHERE order_id = $1
ORDER BY created_at ASC;