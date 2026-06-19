ALTER TABLE products
ADD COLUMN external_id INT UNIQUE,
ADD COLUMN image_url TEXT;

CREATE UNIQUE INDEX idx_products_external_id ON products(external_id);
