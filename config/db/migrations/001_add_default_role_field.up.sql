-- Add is_default_for_new_users field to roles table
-- This migration adds the field for setting default role for new users

-- Step 1: Add the column with default value
ALTER TABLE roles ADD COLUMN is_default_for_new_users BOOLEAN NOT NULL DEFAULT false;

-- Step 2: Add unique partial index to ensure only one default role
CREATE UNIQUE INDEX IF NOT EXISTS idx_roles_default_for_new_users
ON roles (is_default_for_new_users)
WHERE is_default_for_new_users = true AND deleted = false;

-- Step 3: Set initial default role (common_user)
-- Only set if no default role exists yet
UPDATE roles
SET is_default_for_new_users = true
WHERE front_id = 'common_user'
AND deleted = false
AND NOT EXISTS (
    SELECT 1 FROM roles
    WHERE is_default_for_new_users = true
    AND deleted = false
);