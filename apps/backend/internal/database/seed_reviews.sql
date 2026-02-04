-- ============================================
-- REVIEW SEED DATA
-- ============================================
-- This script adds sample reviews for existing products

-- Get delivered order and products info
DO $$
DECLARE
    v_customer1_id UUID;
    v_customer2_id UUID;
    v_order1_id UUID;
    v_order2_id UUID;
    v_dell_laptop_id UUID;
    v_macbook_id UUID;
    v_iphone_id UUID;
    v_samsung_id UUID;
    v_earbuds_id UUID;
    v_denim_jacket_id UUID;
    v_tshirt_pack_id UUID;
    v_floral_dress_id UUID;
    v_handbag_id UUID;
    v_psychology_book_id UUID;
    v_atomic_habits_id UUID;
    v_clean_code_id UUID;
    v_alchemist_id UUID;
BEGIN
    -- Get customer IDs
    SELECT id INTO v_customer1_id FROM users WHERE email = 'john.doe@example.com';
    SELECT id INTO v_customer2_id FROM users WHERE email = 'jane.smith@example.com';
    
    -- Get order IDs
    SELECT id INTO v_order1_id FROM orders WHERE order_number = 'NP-2026-00001';
    SELECT id INTO v_order2_id FROM orders WHERE order_number = 'NP-2026-00002';
    
    -- Get product IDs
    SELECT id INTO v_dell_laptop_id FROM products WHERE sku = 'DELL-XPS15-001';
    SELECT id INTO v_macbook_id FROM products WHERE sku = 'APPLE-MBA-M2-001';
    SELECT id INTO v_iphone_id FROM products WHERE sku = 'APPLE-IP15P-001';
    SELECT id INTO v_samsung_id FROM products WHERE sku = 'SAMSUNG-S24U-001';
    SELECT id INTO v_earbuds_id FROM products WHERE sku = 'TECH-EARBUDS-001';
    SELECT id INTO v_denim_jacket_id FROM products WHERE sku = 'FASHION-DJ-001';
    SELECT id INTO v_tshirt_pack_id FROM products WHERE sku = 'FASHION-TS-001';
    SELECT id INTO v_floral_dress_id FROM products WHERE sku = 'FASHION-FD-001';
    SELECT id INTO v_handbag_id FROM products WHERE sku = 'FASHION-HB-001';
    SELECT id INTO v_psychology_book_id FROM products WHERE sku = 'BOOK-PSY-001';
    SELECT id INTO v_atomic_habits_id FROM products WHERE sku = 'BOOK-ATH-001';
    SELECT id INTO v_clean_code_id FROM products WHERE sku = 'BOOK-CC-001';
    SELECT id INTO v_alchemist_id FROM products WHERE sku = 'BOOK-ALC-001';

    -- Update order 1 to delivered status with proper dates
    UPDATE orders 
    SET 
        status = 'delivered',
        delivered_at = CURRENT_TIMESTAMP - INTERVAL '5 days',
        confirmed_at = CURRENT_TIMESTAMP - INTERVAL '10 days',
        shipped_at = CURRENT_TIMESTAMP - INTERVAL '7 days'
    WHERE id = v_order1_id;

    -- Update order 2 to delivered status
    UPDATE orders 
    SET 
        status = 'delivered',
        delivered_at = CURRENT_TIMESTAMP - INTERVAL '1 day',
        confirmed_at = CURRENT_TIMESTAMP - INTERVAL '2 days',
        shipped_at = CURRENT_TIMESTAMP - INTERVAL '2 days'
    WHERE id = v_order2_id;

    -- Add order items for order 1 if they don't exist
    INSERT INTO order_items (order_id, product_id, shop_id, product_name, quantity, unit_price, subtotal)
    SELECT v_order1_id, v_dell_laptop_id, p.shop_id, p.name, 1, p.price, p.price
    FROM products p WHERE p.id = v_dell_laptop_id
    ON CONFLICT DO NOTHING;

    INSERT INTO order_items (order_id, product_id, shop_id, product_name, quantity, unit_price, subtotal)
    SELECT v_order1_id, v_earbuds_id, p.shop_id, p.name, 1, p.price, p.price
    FROM products p WHERE p.id = v_earbuds_id
    ON CONFLICT DO NOTHING;

    -- Add order items for order 2 if they don't exist
    INSERT INTO order_items (order_id, product_id, shop_id, product_name, quantity, unit_price, subtotal)
    SELECT v_order2_id, v_denim_jacket_id, p.shop_id, p.name, 1, p.price, p.price
    FROM products p WHERE p.id = v_denim_jacket_id
    ON CONFLICT DO NOTHING;

    INSERT INTO order_items (order_id, product_id, shop_id, product_name, quantity, unit_price, subtotal)
    SELECT v_order2_id, v_psychology_book_id, p.shop_id, p.name, 1, p.price, p.price
    FROM products p WHERE p.id = v_psychology_book_id
    ON CONFLICT DO NOTHING;

    -- ============================================
    -- DELL XPS 15 REVIEWS (5 reviews)
    -- ============================================
    INSERT INTO reviews (product_id, user_id, order_id, rating, title, comment, is_verified_purchase, is_approved, helpful_count, created_at, updated_at) VALUES
    (v_dell_laptop_id, v_customer1_id, v_order1_id, 5, 'Excellent laptop for professionals!', 'I''ve been using this Dell XPS 15 for two weeks now and I''m extremely impressed. The 4K display is stunning, the performance is top-notch with the i7 processor, and the build quality is excellent. Perfect for my video editing work. Highly recommended!', TRUE, TRUE, 12, CURRENT_TIMESTAMP - INTERVAL '3 days', CURRENT_TIMESTAMP - INTERVAL '3 days'),
    (v_dell_laptop_id, v_customer2_id, NULL, 4, 'Great performance, battery could be better', 'Overall a fantastic laptop. The screen is beautiful and it handles multiple applications smoothly. My only complaint is the battery life - it drains faster than expected when doing intensive tasks. Still worth the price though.', FALSE, TRUE, 8, CURRENT_TIMESTAMP - INTERVAL '5 days', CURRENT_TIMESTAMP - INTERVAL '5 days'),
    (v_dell_laptop_id, v_customer1_id, NULL, 5, 'Best purchase this year', 'Coming from an older laptop, this is a huge upgrade. Everything is so fast and responsive. The keyboard is comfortable for long typing sessions. Storage is plenty for my needs. Very happy with this purchase!', FALSE, TRUE, 15, CURRENT_TIMESTAMP - INTERVAL '7 days', CURRENT_TIMESTAMP - INTERVAL '7 days'),
    (v_dell_laptop_id, v_customer2_id, NULL, 4, 'Solid laptop with premium feel', 'The build quality is impressive - feels very premium. Display is gorgeous. Runs all my development tools without any lag. A bit heavy to carry around but that''s expected with a 15-inch laptop. Good value for money.', FALSE, TRUE, 6, CURRENT_TIMESTAMP - INTERVAL '10 days', CURRENT_TIMESTAMP - INTERVAL '10 days'),
    (v_dell_laptop_id, v_customer1_id, NULL, 5, 'Perfect for content creation', 'As a content creator, this laptop has been a game changer. The color accuracy of the 4K display is excellent, rendering is fast, and it stays cool even under heavy load. Couldn''t be happier!', FALSE, TRUE, 20, CURRENT_TIMESTAMP - INTERVAL '12 days', CURRENT_TIMESTAMP - INTERVAL '12 days');

    -- ============================================
    -- MACBOOK AIR M2 REVIEWS (4 reviews)
    -- ============================================
    INSERT INTO reviews (product_id, user_id, order_id, rating, title, comment, is_verified_purchase, is_approved, helpful_count, created_at, updated_at) VALUES
    (v_macbook_id, v_customer2_id, NULL, 5, 'M2 chip is incredible!', 'Switched from Intel Mac and the difference is night and day. The M2 chip handles everything I throw at it effortlessly. Battery life is amazing - lasts all day on a single charge. The fanless design means it''s completely silent. Love it!', FALSE, TRUE, 25, CURRENT_TIMESTAMP - INTERVAL '4 days', CURRENT_TIMESTAMP - INTERVAL '4 days'),
    (v_macbook_id, v_customer1_id, NULL, 5, 'Perfect for developers', 'Been coding on this for a month now. Xcode runs smoothly, Docker containers work great, and I can have multiple IDEs open without any slowdown. The unified memory architecture really makes a difference. Highly recommend for developers!', FALSE, TRUE, 18, CURRENT_TIMESTAMP - INTERVAL '6 days', CURRENT_TIMESTAMP - INTERVAL '6 days'),
    (v_macbook_id, v_customer2_id, NULL, 4, 'Great but needs more ports', 'Fantastic laptop overall. Fast, quiet, and the screen is beautiful. My only gripe is the limited ports - would have liked more USB-C ports. Had to buy a hub. Otherwise, it''s perfect for daily use and light creative work.', FALSE, TRUE, 10, CURRENT_TIMESTAMP - INTERVAL '8 days', CURRENT_TIMESTAMP - INTERVAL '8 days'),
    (v_macbook_id, v_customer1_id, NULL, 5, 'Best MacBook yet!', 'The M2 MacBook Air is the perfect balance of performance and portability. It''s thin, light, and powerful. Battery life is outstanding. The midnight color looks premium. This is Apple at its best!', FALSE, TRUE, 22, CURRENT_TIMESTAMP - INTERVAL '11 days', CURRENT_TIMESTAMP - INTERVAL '11 days');

    -- ============================================
    -- WIRELESS EARBUDS REVIEWS (6 reviews)
    -- ============================================
    INSERT INTO reviews (product_id, user_id, order_id, rating, title, comment, is_verified_purchase, is_approved, helpful_count, created_at, updated_at) VALUES
    (v_earbuds_id, v_customer1_id, v_order1_id, 5, 'Amazing sound quality!', 'These earbuds sound incredible! The noise cancellation works really well, and the battery life is as advertised. They fit comfortably and stay secure even during workouts. Best wireless earbuds I''ve owned.', TRUE, TRUE, 14, CURRENT_TIMESTAMP - INTERVAL '2 days', CURRENT_TIMESTAMP - INTERVAL '2 days'),
    (v_earbuds_id, v_customer2_id, NULL, 4, 'Great value for money', 'For the price, these are excellent. ANC is effective, sound is clear with good bass. Case is compact and charges quickly. Only minor issue is the touch controls can be a bit finicky sometimes. Still recommend!', FALSE, TRUE, 9, CURRENT_TIMESTAMP - INTERVAL '4 days', CURRENT_TIMESTAMP - INTERVAL '4 days'),
    (v_earbuds_id, v_customer1_id, NULL, 5, 'Better than AirPods!', 'Honestly didn''t expect much but these are better than my friend''s AirPods Pro. Sound quality is excellent, ANC is strong, and battery lasts forever. The case is small enough to fit in my pocket. Very impressed!', FALSE, TRUE, 28, CURRENT_TIMESTAMP - INTERVAL '5 days', CURRENT_TIMESTAMP - INTERVAL '5 days'),
    (v_earbuds_id, v_customer2_id, NULL, 5, 'Perfect for daily commute', 'Use these every day on my commute. The ANC blocks out all the traffic noise. Transparency mode is useful when I need to hear announcements. Battery easily lasts a week with my usage. Love them!', FALSE, TRUE, 11, CURRENT_TIMESTAMP - INTERVAL '7 days', CURRENT_TIMESTAMP - INTERVAL '7 days'),
    (v_earbuds_id, v_customer1_id, NULL, 4, 'Solid earbuds', 'Good sound, comfortable fit, and reliable connection. The app has useful EQ settings. They''re not quite as premium as top-tier brands but for this price point, they''re hard to beat. Happy with the purchase.', FALSE, TRUE, 7, CURRENT_TIMESTAMP - INTERVAL '9 days', CURRENT_TIMESTAMP - INTERVAL '9 days'),
    (v_earbuds_id, v_customer2_id, NULL, 5, 'Exceeded expectations', 'Bought these as a backup pair but they''re so good I use them as my primary earbuds now. Call quality is clear, music sounds great, and the charging case is well-designed. Definitely recommend!', FALSE, TRUE, 16, CURRENT_TIMESTAMP - INTERVAL '13 days', CURRENT_TIMESTAMP - INTERVAL '13 days');

    -- ============================================
    -- DENIM JACKET REVIEWS (4 reviews)
    -- ============================================
    INSERT INTO reviews (product_id, user_id, order_id, rating, title, comment, is_verified_purchase, is_approved, helpful_count, created_at, updated_at) VALUES
    (v_denim_jacket_id, v_customer2_id, v_order2_id, 5, 'Perfect vintage look!', 'This denim jacket is exactly what I was looking for. The vintage wash looks authentic and the quality is great. Fits perfectly and goes well with everything. Got so many compliments already!', TRUE, TRUE, 13, CURRENT_TIMESTAMP - INTERVAL '1 day', CURRENT_TIMESTAMP - INTERVAL '1 day'),
    (v_denim_jacket_id, v_customer1_id, NULL, 4, 'Good quality denim', 'Nice thick denim that feels durable. The fit is slightly oversized which works for the casual look. Wash it before wearing as it''s a bit stiff initially. Good value for the price.', FALSE, TRUE, 8, CURRENT_TIMESTAMP - INTERVAL '3 days', CURRENT_TIMESTAMP - INTERVAL '3 days'),
    (v_denim_jacket_id, v_customer2_id, NULL, 5, 'My new favorite jacket', 'Love this jacket! The color is beautiful, it''s comfortable, and very versatile. Can dress it up or down. True to size. The stitching and quality are impressive. Highly recommend!', FALSE, TRUE, 19, CURRENT_TIMESTAMP - INTERVAL '6 days', CURRENT_TIMESTAMP - INTERVAL '6 days'),
    (v_denim_jacket_id, v_customer1_id, NULL, 4, 'Classic style', 'A wardrobe staple. Good quality construction, nice wash effect. Buttons are sturdy. Would have given 5 stars but the pockets could be a bit deeper. Otherwise, very happy with it.', FALSE, TRUE, 6, CURRENT_TIMESTAMP - INTERVAL '8 days', CURRENT_TIMESTAMP - INTERVAL '8 days');

    -- ============================================
    -- BOOKS REVIEWS
    -- ============================================
    
    -- Psychology of Money (5 reviews)
    INSERT INTO reviews (product_id, user_id, order_id, rating, title, comment, is_verified_purchase, is_approved, helpful_count, created_at, updated_at) VALUES
    (v_psychology_book_id, v_customer2_id, v_order2_id, 5, 'Life-changing book!', 'This book completely changed my perspective on money and wealth. Morgan Housel explains complex financial concepts in simple, relatable stories. Everyone should read this regardless of their financial situation. Highly recommended!', TRUE, TRUE, 42, CURRENT_TIMESTAMP - INTERVAL '1 day', CURRENT_TIMESTAMP - INTERVAL '1 day'),
    (v_psychology_book_id, v_customer1_id, NULL, 5, 'Best finance book I''ve read', 'Unlike other finance books full of jargon, this one is easy to understand and actually enjoyable to read. The lessons are practical and timeless. Already recommended it to my friends and family!', FALSE, TRUE, 31, CURRENT_TIMESTAMP - INTERVAL '3 days', CURRENT_TIMESTAMP - INTERVAL '3 days'),
    (v_psychology_book_id, v_customer2_id, NULL, 5, 'Must-read for everyone', 'Not just about money - it''s about human behavior and decision making. The stories are engaging and the insights are profound. I keep coming back to certain chapters. A true masterpiece!', FALSE, TRUE, 38, CURRENT_TIMESTAMP - INTERVAL '5 days', CURRENT_TIMESTAMP - INTERVAL '5 days'),
    (v_psychology_book_id, v_customer1_id, NULL, 4, 'Insightful and practical', 'Great book with valuable lessons. Some concepts repeated a bit but overall excellent. The chapter on getting wealthy vs staying wealthy was eye-opening. Worth every rupee!', FALSE, TRUE, 24, CURRENT_TIMESTAMP - INTERVAL '7 days', CURRENT_TIMESTAMP - INTERVAL '7 days'),
    (v_psychology_book_id, v_customer2_id, NULL, 5, 'Changed my financial mindset', 'I''m more conscious about my spending and saving habits after reading this. The book is well-written and the examples are relatable. It''s not about getting rich quick but building sustainable wealth. Excellent!', FALSE, TRUE, 29, CURRENT_TIMESTAMP - INTERVAL '10 days', CURRENT_TIMESTAMP - INTERVAL '10 days');

    -- Atomic Habits (4 reviews)
    INSERT INTO reviews (product_id, user_id, order_id, rating, title, comment, is_verified_purchase, is_approved, helpful_count, created_at, updated_at) VALUES
    (v_atomic_habits_id, v_customer1_id, NULL, 5, 'Practical and actionable', 'James Clear breaks down habit formation in a way that''s easy to understand and implement. I''ve already applied several strategies and seen results. This book is a game-changer for self-improvement!', FALSE, TRUE, 35, CURRENT_TIMESTAMP - INTERVAL '2 days', CURRENT_TIMESTAMP - INTERVAL '2 days'),
    (v_atomic_habits_id, v_customer2_id, NULL, 5, 'Best self-help book ever', 'Not your typical motivational fluff - this is based on science and provides real strategies. The 1% improvement concept is powerful. Already building better habits thanks to this book!', FALSE, TRUE, 41, CURRENT_TIMESTAMP - INTERVAL '4 days', CURRENT_TIMESTAMP - INTERVAL '4 days'),
    (v_atomic_habits_id, v_customer1_id, NULL, 5, 'Transform your life', 'This book should be mandatory reading. It''s helped me understand why I struggle with certain habits and how to fix them. The four laws of behavior change are brilliant. Highly recommend!', FALSE, TRUE, 27, CURRENT_TIMESTAMP - INTERVAL '6 days', CURRENT_TIMESTAMP - INTERVAL '6 days'),
    (v_atomic_habits_id, v_customer2_id, NULL, 4, 'Very informative', 'Lots of useful information backed by research. Some parts felt a bit repetitive but the core concepts are solid. Already seeing improvements in my daily routine. Good read!', FALSE, TRUE, 18, CURRENT_TIMESTAMP - INTERVAL '9 days', CURRENT_TIMESTAMP - INTERVAL '9 days');

    -- Clean Code (3 reviews)
    INSERT INTO reviews (product_id, user_id, order_id, rating, title, comment, is_verified_purchase, is_approved, helpful_count, created_at, updated_at) VALUES
    (v_clean_code_id, v_customer1_id, NULL, 5, 'Essential for every developer', 'This book should be required reading for all software developers. Uncle Bob''s principles have made me a better programmer. My code is now more readable and maintainable. A classic!', FALSE, TRUE, 45, CURRENT_TIMESTAMP - INTERVAL '5 days', CURRENT_TIMESTAMP - INTERVAL '5 days'),
    (v_clean_code_id, v_customer2_id, NULL, 4, 'Great principles, some outdated examples', 'The core concepts are timeless and valuable. Some code examples feel dated but the principles still apply. Has significantly improved my code quality. Recommend for junior and mid-level developers.', FALSE, TRUE, 22, CURRENT_TIMESTAMP - INTERVAL '8 days', CURRENT_TIMESTAMP - INTERVAL '8 days'),
    (v_clean_code_id, v_customer1_id, NULL, 5, 'Changed how I write code', 'After reading this, I can''t help but refactor messy code. The book teaches you to think differently about code structure and naming. My team has adopted many of these practices. Excellent resource!', FALSE, TRUE, 33, CURRENT_TIMESTAMP - INTERVAL '12 days', CURRENT_TIMESTAMP - INTERVAL '12 days');

    RAISE NOTICE '============================================';
    RAISE NOTICE 'Review Seed Data Added Successfully!';
    RAISE NOTICE '============================================';
    RAISE NOTICE 'Dell XPS 15: 5 reviews';
    RAISE NOTICE 'MacBook Air M2: 4 reviews';
    RAISE NOTICE 'Wireless Earbuds: 6 reviews';
    RAISE NOTICE 'Denim Jacket: 4 reviews';
    RAISE NOTICE 'Psychology of Money: 5 reviews';
    RAISE NOTICE 'Atomic Habits: 4 reviews';
    RAISE NOTICE 'Clean Code: 3 reviews';
    RAISE NOTICE '============================================';
    RAISE NOTICE 'Total: 31 reviews added';
    RAISE NOTICE 'Orders updated to delivered status';
    RAISE NOTICE '============================================';
END $$;
