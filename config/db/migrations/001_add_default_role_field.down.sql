-- Rollback migration for is_default_for_new_users field
-- This removes the field and associated constraints

-- Step 1: Drop the unique index
DROP INDEX IF EXISTS idx_roles_default_for_new_users;

-- Step 2: Remove the column
ALTER TABLE roles DROP COLUMN IF EXISTS is_default_for_new_users;