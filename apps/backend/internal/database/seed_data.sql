-- Seed Data for Nepify E-commerce Platform
-- This script populates the database with sample data for testing

-- Note: In production, users are created via Clerk webhook
-- For testing, we'll create some sample users with mock clerk_ids

-- Clear existing data (be careful in production!)
-- TRUNCATE TABLE order_items, orders, addresses, product_images, products, categories, shops, users CASCADE;

-- ============================================
-- USERS (Sample test users)
-- ============================================
INSERT INTO users (clerk_id, email, username, first_name, last_name, phone, role, is_active) VALUES
    ('user_test_customer1', 'john.doe@example.com', 'johndoe', 'John', 'Doe', '+977-9841234567', 'customer', TRUE),
    ('user_test_customer2', 'jane.smith@example.com', 'janesmith', 'Jane', 'Smith', '+977-9841234568', 'customer', TRUE),
    ('user_test_vendor1', 'vendor1@nepify.com', 'techstore', 'Ramesh', 'Sharma', '+977-9841234569', 'vendor', TRUE),
    ('user_test_vendor2', 'vendor2@nepify.com', 'fashionhub', 'Sita', 'Thapa', '+977-9841234570', 'vendor', TRUE),
    ('user_test_vendor3', 'vendor3@nepify.com', 'bookworld', 'Prakash', 'Adhikari', '+977-9841234571', 'vendor', TRUE),
    ('user_test_admin', 'admin@nepify.com', 'admin', 'Admin', 'User', '+977-9841234572', 'admin', TRUE)
ON CONFLICT (clerk_id) DO NOTHING;

-- ============================================
-- SHOPS (Vendor stores)
-- ============================================
INSERT INTO shops (vendor_id, shop_name, slug, description, contact_email, contact_phone, address, city, state, country, is_active, is_verified, rating, total_reviews, total_sales)
SELECT 
    u.id,
    'Tech Store Nepal',
    'tech-store-nepal',
    'Your one-stop shop for latest electronics, gadgets, and computer accessories. Authorized dealer of top brands.',
    'contact@techstore.com',
    '+977-9841234569',
    'Putalisadak',
    'Kathmandu',
    'Bagmati',
    'Nepal',
    TRUE,
    TRUE,
    4.5,
    127,
    450
FROM users u WHERE u.email = 'vendor1@nepify.com'
ON CONFLICT (vendor_id) DO NOTHING;

INSERT INTO shops (vendor_id, shop_name, slug, description, contact_email, contact_phone, address, city, state, country, is_active, is_verified, rating, total_reviews, total_sales)
SELECT 
    u.id,
    'Fashion Hub',
    'fashion-hub',
    'Trendy clothing and accessories for men and women. Latest fashion at affordable prices.',
    'info@fashionhub.com',
    '+977-9841234570',
    'New Road',
    'Kathmandu',
    'Bagmati',
    'Nepal',
    TRUE,
    TRUE,
    4.7,
    89,
    320
FROM users u WHERE u.email = 'vendor2@nepify.com'
ON CONFLICT (vendor_id) DO NOTHING;

INSERT INTO shops (vendor_id, shop_name, slug, description, contact_email, contact_phone, address, city, state, country, is_active, is_verified, rating, total_reviews, total_sales)
SELECT 
    u.id,
    'Book World',
    'book-world',
    'Wide collection of books - fiction, non-fiction, academic, and more. Books for every reader.',
    'hello@bookworld.com',
    '+977-9841234571',
    'Bhatbhateni',
    'Kathmandu',
    'Bagmati',
    'Nepal',
    TRUE,
    TRUE,
    4.8,
    210,
    580
FROM users u WHERE u.email = 'vendor3@nepify.com'
ON CONFLICT (vendor_id) DO NOTHING;

-- ============================================
-- CATEGORIES
-- ============================================
INSERT INTO categories (name, slug, description, is_active, display_order) VALUES
    ('Electronics', 'electronics', 'Electronic devices and gadgets', TRUE, 1),
    ('Fashion', 'fashion', 'Clothing and accessories', TRUE, 2),
    ('Books', 'books', 'Books and publications', TRUE, 3),
    ('Home & Living', 'home-living', 'Home decor and furniture', TRUE, 4),
    ('Sports', 'sports', 'Sports equipment and accessories', TRUE, 5)
ON CONFLICT (slug) DO NOTHING;

-- Sub-categories
INSERT INTO categories (name, slug, description, parent_id, is_active, display_order)
SELECT 'Laptops', 'laptops', 'Laptop computers', c.id, TRUE, 1
FROM categories c WHERE c.slug = 'electronics'
ON CONFLICT (slug) DO NOTHING;

INSERT INTO categories (name, slug, description, parent_id, is_active, display_order)
SELECT 'Smartphones', 'smartphones', 'Mobile phones', c.id, TRUE, 2
FROM categories c WHERE c.slug = 'electronics'
ON CONFLICT (slug) DO NOTHING;

INSERT INTO categories (name, slug, description, parent_id, is_active, display_order)
SELECT 'Mens Wear', 'mens-wear', 'Clothing for men', c.id, TRUE, 1
FROM categories c WHERE c.slug = 'fashion'
ON CONFLICT (slug) DO NOTHING;

INSERT INTO categories (name, slug, description, parent_id, is_active, display_order)
SELECT 'Womens Wear', 'womens-wear', 'Clothing for women', c.id, TRUE, 2
FROM categories c WHERE c.slug = 'fashion'
ON CONFLICT (slug) DO NOTHING;

-- ============================================
-- PRODUCTS
-- ============================================

-- Tech Store Products
INSERT INTO products (shop_id, category_id, name, slug, description, short_description, sku, price, compare_at_price, stock_quantity, image_url, is_active, is_featured)
SELECT 
    s.id,
    c.id,
    'Dell XPS 15 Laptop',
    'dell-xps-15-laptop',
    'High-performance laptop with Intel i7 processor, 16GB RAM, 512GB SSD, and 15.6" 4K display. Perfect for professionals and creators.',
    'Premium laptop for professionals',
    'DELL-XPS15-001',
    189999.00,
    209999.00,
    15,
    'https://images.unsplash.com/photo-1593642632823-8f785ba67e45?w=800',
    TRUE,
    TRUE
FROM shops s, categories c 
WHERE s.slug = 'tech-store-nepal' AND c.slug = 'laptops'
ON CONFLICT (sku) DO NOTHING;

INSERT INTO products (shop_id, category_id, name, slug, description, short_description, sku, price, compare_at_price, stock_quantity, image_url, is_active, is_featured)
SELECT 
    s.id,
    c.id,
    'MacBook Air M2',
    'macbook-air-m2',
    'Apple MacBook Air with M2 chip, 8GB unified memory, 256GB SSD. Ultra-thin and lightweight design with incredible battery life.',
    'Latest MacBook Air with M2 chip',
    'APPLE-MBA-M2-001',
    159999.00,
    169999.00,
    8,
    'https://images.unsplash.com/photo-1517336714731-489689fd1ca8?w=800',
    TRUE,
    TRUE
FROM shops s, categories c 
WHERE s.slug = 'tech-store-nepal' AND c.slug = 'laptops'
ON CONFLICT (sku) DO NOTHING;

INSERT INTO products (shop_id, category_id, name, slug, description, short_description, sku, price, stock_quantity, image_url, is_active)
SELECT 
    s.id,
    c.id,
    'iPhone 15 Pro',
    'iphone-15-pro',
    'iPhone 15 Pro with A17 Pro chip, advanced camera system, and titanium design. Available in multiple colors.',
    'Latest iPhone with Pro features',
    'APPLE-IP15P-001',
    149999.00,
    12,
    'https://images.unsplash.com/photo-1592286927505-90fd157f0edc?w=800',
    TRUE
FROM shops s, categories c 
WHERE s.slug = 'tech-store-nepal' AND c.slug = 'smartphones'
ON CONFLICT (sku) DO NOTHING;

INSERT INTO products (shop_id, category_id, name, slug, description, short_description, sku, price, compare_at_price, stock_quantity, image_url, is_active)
SELECT 
    s.id,
    c.id,
    'Samsung Galaxy S24 Ultra',
    'samsung-galaxy-s24-ultra',
    'Samsung flagship with 200MP camera, S Pen, and incredible display. Perfect for productivity and photography.',
    'Premium Samsung flagship phone',
    'SAMSUNG-S24U-001',
    139999.00,
    149999.00,
    20,
    'https://images.unsplash.com/photo-1610945415295-d9bbf067e59c?w=800',
    TRUE
FROM shops s, categories c 
WHERE s.slug = 'tech-store-nepal' AND c.slug = 'smartphones'
ON CONFLICT (sku) DO NOTHING;

INSERT INTO products (shop_id, category_id, name, slug, description, short_description, sku, price, stock_quantity, image_url, is_active)
SELECT 
    s.id,
    c.id,
    'Wireless Earbuds Pro',
    'wireless-earbuds-pro',
    'Premium wireless earbuds with active noise cancellation, 30-hour battery life, and crystal clear sound.',
    'Premium wireless earbuds with ANC',
    'TECH-EARBUDS-001',
    8999.00,
    45,
    'https://images.unsplash.com/photo-1590658268037-6bf12165a8df?w=800',
    TRUE
FROM shops s, categories c 
WHERE s.slug = 'tech-store-nepal' AND c.slug = 'electronics'
ON CONFLICT (sku) DO NOTHING;

-- Fashion Hub Products
INSERT INTO products (shop_id, category_id, name, slug, description, short_description, sku, price, compare_at_price, stock_quantity, image_url, is_active, is_featured)
SELECT 
    s.id,
    c.id,
    'Classic Denim Jacket',
    'classic-denim-jacket',
    'Premium denim jacket with vintage wash. Perfect for casual outings and layering. Available in multiple sizes.',
    'Stylish vintage denim jacket',
    'FASHION-DJ-001',
    3999.00,
    4999.00,
    30,
    'https://images.unsplash.com/photo-1551028719-00167b16eac5?w=800',
    TRUE,
    TRUE
FROM shops s, categories c 
WHERE s.slug = 'fashion-hub' AND c.slug = 'mens-wear'
ON CONFLICT (sku) DO NOTHING;

INSERT INTO products (shop_id, category_id, name, slug, description, short_description, sku, price, stock_quantity, image_url, is_active)
SELECT 
    s.id,
    c.id,
    'Cotton T-Shirt Pack',
    'cotton-tshirt-pack',
    'Pack of 3 premium cotton t-shirts in assorted colors. Comfortable fit and breathable fabric.',
    '3-pack premium cotton tees',
    'FASHION-TS-001',
    1499.00,
    100,
    'https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w=800',
    TRUE
FROM shops s, categories c 
WHERE s.slug = 'fashion-hub' AND c.slug = 'mens-wear'
ON CONFLICT (sku) DO NOTHING;

INSERT INTO products (shop_id, category_id, name, slug, description, short_description, sku, price, compare_at_price, stock_quantity, image_url, is_active, is_featured)
SELECT 
    s.id,
    c.id,
    'Floral Summer Dress',
    'floral-summer-dress',
    'Beautiful floral print summer dress. Light and comfortable fabric perfect for warm weather.',
    'Elegant floral summer dress',
    'FASHION-FD-001',
    2999.00,
    3499.00,
    40,
    'https://images.unsplash.com/photo-1572804013309-59a88b7e92f1?w=800',
    TRUE,
    TRUE
FROM shops s, categories c 
WHERE s.slug = 'fashion-hub' AND c.slug = 'womens-wear'
ON CONFLICT (sku) DO NOTHING;

INSERT INTO products (shop_id, category_id, name, slug, description, short_description, sku, price, stock_quantity, image_url, is_active)
SELECT 
    s.id,
    c.id,
    'Leather Handbag',
    'leather-handbag',
    'Genuine leather handbag with multiple compartments. Elegant design suitable for both casual and formal occasions.',
    'Premium leather handbag',
    'FASHION-HB-001',
    5999.00,
    25,
    'https://images.unsplash.com/photo-1584917865442-de89df76afd3?w=800',
    TRUE
FROM shops s, categories c 
WHERE s.slug = 'fashion-hub' AND c.slug = 'womens-wear'
ON CONFLICT (sku) DO NOTHING;

-- Book World Products
INSERT INTO products (shop_id, category_id, name, slug, description, short_description, sku, price, stock_quantity, image_url, is_active, is_featured)
SELECT 
    s.id,
    c.id,
    'The Psychology of Money',
    'psychology-of-money',
    'Best-selling book by Morgan Housel. Timeless lessons on wealth, greed, and happiness.',
    'Best-seller on wealth and happiness',
    'BOOK-PSY-001',
    799.00,
    50,
    'https://images.unsplash.com/photo-1544947950-fa07a98d237f?w=800',
    TRUE,
    TRUE
FROM shops s, categories c 
WHERE s.slug = 'book-world' AND c.slug = 'books'
ON CONFLICT (sku) DO NOTHING;

INSERT INTO products (shop_id, category_id, name, slug, description, short_description, sku, price, stock_quantity, image_url, is_active, is_featured)
SELECT 
    s.id,
    c.id,
    'Atomic Habits',
    'atomic-habits',
    'By James Clear. An easy and proven way to build good habits and break bad ones.',
    'Transform your habits',
    'BOOK-ATH-001',
    899.00,
    60,
    'https://images.unsplash.com/photo-1512820790803-83ca734da794?w=800',
    TRUE,
    TRUE
FROM shops s, categories c 
WHERE s.slug = 'book-world' AND c.slug = 'books'
ON CONFLICT (sku) DO NOTHING;

INSERT INTO products (shop_id, category_id, name, slug, description, short_description, sku, price, stock_quantity, image_url, is_active)
SELECT 
    s.id,
    c.id,
    'Clean Code',
    'clean-code',
    'A handbook of agile software craftsmanship by Robert C. Martin. Essential reading for developers.',
    'Essential for developers',
    'BOOK-CC-001',
    1499.00,
    35,
    'https://images.unsplash.com/photo-1532012197267-da84d127e765?w=800',
    TRUE
FROM shops s, categories c 
WHERE s.slug = 'book-world' AND c.slug = 'books'
ON CONFLICT (sku) DO NOTHING;

INSERT INTO products (shop_id, category_id, name, slug, description, short_description, sku, price, stock_quantity, is_active)
SELECT 
    s.id,
    c.id,
    'The Alchemist',
    'the-alchemist',
    'Paulo Coelho''s masterpiece about following your dreams. A modern classic that has inspired millions.',
    'Follow your dreams',
    'BOOK-ALC-001',
    599.00,
    80,
    TRUE
FROM shops s, categories c 
WHERE s.slug = 'book-world' AND c.slug = 'books'
ON CONFLICT (sku) DO NOTHING;

-- ============================================
-- ADDRESSES (Sample customer addresses)
-- ============================================
INSERT INTO addresses (user_id, full_name, phone, address_line1, city, state, country, is_default, address_type)
SELECT 
    u.id,
    'John Doe',
    '+977-9841234567',
    'Thamel, Ward No. 26',
    'Kathmandu',
    'Bagmati',
    'Nepal',
    TRUE,
    'both'
FROM users u WHERE u.email = 'john.doe@example.com'
ON CONFLICT DO NOTHING;

INSERT INTO addresses (user_id, full_name, phone, address_line1, city, state, country, is_default, address_type)
SELECT 
    u.id,
    'Jane Smith',
    '+977-9841234568',
    'Lazimpat, Ward No. 3',
    'Kathmandu',
    'Bagmati',
    'Nepal',
    TRUE,
    'both'
FROM users u WHERE u.email = 'jane.smith@example.com'
ON CONFLICT DO NOTHING;

-- ============================================
-- SAMPLE ORDERS (For testing order endpoints)
-- ============================================
INSERT INTO orders (user_id, order_number, status, subtotal, shipping_cost, tax, total, payment_method, payment_status, created_at)
SELECT 
    u.id,
    'NP-2026-00001',
    'delivered',
    190798.00,
    500.00,
    0.00,
    191298.00,
    'esewa',
    'paid',
    CURRENT_TIMESTAMP - INTERVAL '10 days'
FROM users u WHERE u.email = 'john.doe@example.com'
ON CONFLICT (order_number) DO NOTHING;

INSERT INTO orders (user_id, order_number, status, subtotal, shipping_cost, tax, total, payment_method, payment_status, created_at)
SELECT 
    u.id,
    'NP-2026-00002',
    'processing',
    3998.00,
    200.00,
    0.00,
    4198.00,
    'khalti',
    'paid',
    CURRENT_TIMESTAMP - INTERVAL '2 days'
FROM users u WHERE u.email = 'jane.smith@example.com'
ON CONFLICT (order_number) DO NOTHING;

-- ============================================
-- Final Statistics
-- ============================================
-- Display summary of seeded data
DO $$
DECLARE
    user_count INT;
    shop_count INT;
    category_count INT;
    product_count INT;
    address_count INT;
    order_count INT;
BEGIN
    SELECT COUNT(*) INTO user_count FROM users;
    SELECT COUNT(*) INTO shop_count FROM shops;
    SELECT COUNT(*) INTO category_count FROM categories;
    SELECT COUNT(*) INTO product_count FROM products;
    SELECT COUNT(*) INTO address_count FROM addresses;
    SELECT COUNT(*) INTO order_count FROM orders;
    
    RAISE NOTICE '============================================';
    RAISE NOTICE 'Seed Data Summary:';
    RAISE NOTICE '============================================';
    RAISE NOTICE 'Users: %', user_count;
    RAISE NOTICE 'Shops: %', shop_count;
    RAISE NOTICE 'Categories: %', category_count;
    RAISE NOTICE 'Products: %', product_count;
    RAISE NOTICE 'Addresses: %', address_count;
    RAISE NOTICE 'Orders: %', order_count;
    RAISE NOTICE '============================================';
END $$;
