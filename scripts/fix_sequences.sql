-- Fix PostgreSQL sequences after manual insert
-- Run this SQL script when you get "duplicate key value violates unique constraint" error
-- This happens when you insert data directly to database without using the application

-- Fix product_stocks sequence
SELECT setval('product_stocks_id_seq', (SELECT COALESCE(MAX(id), 1) FROM product_stocks));

-- Fix other tables sequences (run these if you have similar issues)
SELECT setval('brands_id_seq', (SELECT COALESCE(MAX(id), 1) FROM brands));
SELECT setval('categories_id_seq', (SELECT COALESCE(MAX(id), 1) FROM categories));
SELECT setval('files_id_seq', (SELECT COALESCE(MAX(id), 1) FROM files));
SELECT setval('locations_id_seq', (SELECT COALESCE(MAX(id), 1) FROM locations));
SELECT setval('product_batches_id_seq', (SELECT COALESCE(MAX(id), 1) FROM product_batches));
SELECT setval('product_batch_tracks_id_seq', (SELECT COALESCE(MAX(id), 1) FROM product_batch_tracks));
SELECT setval('product_items_id_seq', (SELECT COALESCE(MAX(id), 1) FROM product_items));
SELECT setval('product_item_tracks_id_seq', (SELECT COALESCE(MAX(id), 1) FROM product_item_tracks));
SELECT setval('products_id_seq', (SELECT COALESCE(MAX(id), 1) FROM products));
SELECT setval('product_stock_tracks_id_seq', (SELECT COALESCE(MAX(id), 1) FROM product_stock_tracks));
SELECT setval('product_units_id_seq', (SELECT COALESCE(MAX(id), 1) FROM product_units));
SELECT setval('product_unit_tracks_id_seq', (SELECT COALESCE(MAX(id), 1) FROM product_unit_tracks));
SELECT setval('roles_id_seq', (SELECT COALESCE(MAX(id), 1) FROM roles));
SELECT setval('users_id_seq', (SELECT COALESCE(MAX(id), 1) FROM users));

-- Verify the sequences are now correct
SELECT 'brands' as table_name, currval('brands_id_seq') as current_sequence, MAX(id) as max_id FROM brands
UNION ALL
SELECT 'categories', currval('categories_id_seq'), MAX(id) FROM categories
UNION ALL
SELECT 'files', currval('files_id_seq'), MAX(id) FROM files
UNION ALL
SELECT 'locations', currval('locations_id_seq'), MAX(id) FROM locations
UNION ALL
SELECT 'product_batches', currval('product_batches_id_seq'), MAX(id) FROM product_batches
UNION ALL
SELECT 'product_batch_tracks', currval('product_batch_tracks_id_seq'), MAX(id) FROM product_batch_tracks
UNION ALL
SELECT 'product_items', currval('product_items_id_seq'), MAX(id) FROM product_items
UNION ALL
SELECT 'product_item_tracks', currval('product_item_tracks_id_seq'), MAX(id) FROM product_item_tracks
UNION ALL
SELECT 'products', currval('products_id_seq'), MAX(id) FROM products
UNION ALL
SELECT 'product_stocks', currval('product_stocks_id_seq'), MAX(id) FROM product_stocks
UNION ALL
SELECT 'product_stock_tracks', currval('product_stock_tracks_id_seq'), MAX(id) FROM product_stock_tracks
UNION ALL
SELECT 'product_units', currval('product_units_id_seq'), MAX(id) FROM product_units
UNION ALL
SELECT 'product_unit_tracks', currval('product_unit_tracks_id_seq'), MAX(id) FROM product_unit_tracks
UNION ALL
SELECT 'roles', currval('roles_id_seq'), MAX(id) FROM roles
UNION ALL
SELECT 'users', currval('users_id_seq'), MAX(id) FROM users;
