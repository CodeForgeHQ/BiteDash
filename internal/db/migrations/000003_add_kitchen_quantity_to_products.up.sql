ALTER TABLE products
ADD COLUMN kitchen_quantity INT NOT NULL DEFAULT 100 CHECK (kitchen_quantity >= 0);
