-- Fix invoice migration: Drop supplier_id, add user_id
-- Run this before migration

-- Step 1: Drop existing invoices and invoice_items if any
TRUNCATE TABLE invoice_items CASCADE;
TRUNCATE TABLE invoices CASCADE;

-- Step 2: Alter invoices table
-- Drop old supplier_id column if exists
ALTER TABLE invoices DROP COLUMN IF EXISTS supplier_id;

-- Add user_id column
-- Note: After running this, run the Go migration which will add the NOT NULL constraint
