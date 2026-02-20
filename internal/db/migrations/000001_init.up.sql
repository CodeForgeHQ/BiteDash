CREATE TABLE users (
    id UUID PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX idx_users_email ON users(email);

CREATE TABLE restaurants (
    id UUID PRIMARY KEY,
    external_id INT UNIQUE,
    name TEXT NOT NULL,
    description TEXT,
    category TEXT NOT NULL,
    address TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX idx_restaurants_category ON restaurants(category);
CREATE INDEX idx_restaurants_name ON restaurants(name);

CREATE TABLE products (
    id UUID PRIMARY KEY,
    restaurant_id UUID NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    description TEXT,
    price NUMERIC(10,2) NOT NULL CHECK (price >= 0),
    is_available BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX idx_products_restaurant_id ON products(restaurant_id);
CREATE INDEX idx_products_name ON products(name);

CREATE TYPE cart_status AS ENUM ('active', 'ordered', 'abandoned');

CREATE TABLE carts (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status cart_status NOT NULL DEFAULT 'active',
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX idx_carts_user_id ON carts(user_id);

-- Только одна активная корзина
CREATE UNIQUE INDEX unique_active_cart_per_user
ON carts(user_id)
WHERE status = 'active';

CREATE TABLE cart_items (
    id UUID PRIMARY KEY,
    cart_id UUID NOT NULL REFERENCES carts(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id),
    quantity INT NOT NULL CHECK (quantity > 0),
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    UNIQUE(cart_id, product_id)
);

CREATE INDEX idx_cart_items_cart_id ON cart_items(cart_id);

CREATE TYPE order_status AS ENUM (
    'pending',
    'paid',
    'cancelled',
    'delivered'
);

CREATE TABLE orders (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    status order_status NOT NULL DEFAULT 'pending',
    total_amount NUMERIC(10,2) NOT NULL CHECK (total_amount >= 0),
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_orders_status ON orders(status);

CREATE TABLE order_items (
    id UUID PRIMARY KEY,
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id UUID NOT NULL,
    product_name TEXT NOT NULL,
    price NUMERIC(10,2) NOT NULL CHECK (price >= 0),
    quantity INT NOT NULL CHECK (quantity > 0)
);

CREATE INDEX idx_order_items_order_id ON order_items(order_id);

CREATE TYPE payment_status AS ENUM (
    'pending',
    'success',
    'failed'
);

CREATE TABLE payments (
    id UUID PRIMARY KEY,
    order_id UUID NOT NULL REFERENCES orders(id),
    status payment_status NOT NULL DEFAULT 'pending',
    provider TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX idx_payments_order_id ON payments(order_id);

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = now();
   RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_restaurants_updated_at
BEFORE UPDATE ON restaurants
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_products_updated_at
BEFORE UPDATE ON products
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_carts_updated_at
BEFORE UPDATE ON carts
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_orders_updated_at
BEFORE UPDATE ON orders
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

