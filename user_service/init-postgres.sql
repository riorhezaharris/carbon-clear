-- PostgreSQL initialization script for user service
-- This script sets up the database for the user service

-- Create the database if it doesn't exist (handled by POSTGRES_DB env var)
-- CREATE DATABASE carbon_clear_users;

-- Connect to the database
\c carbon_clear_users;

-- Create extensions if needed
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Note: The users table will be created by GORM AutoMigrate
-- GORM will handle table creation, indexes, and constraints automatically
-- GORM also handles the updated_at field automatically

-- Grant permissions
GRANT ALL PRIVILEGES ON DATABASE carbon_clear_users TO postgres;

-- Print completion message
DO $$
BEGIN
    RAISE NOTICE 'PostgreSQL initialization completed for user service';
END $$;
