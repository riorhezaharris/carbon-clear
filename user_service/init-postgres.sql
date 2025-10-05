-- PostgreSQL initialization script for user service
-- This script sets up the database for the user service

-- Create the database if it doesn't exist (handled by POSTGRES_DB env var)
-- CREATE DATABASE carbon_clear_users;

-- Connect to the database
\c carbon_clear_users;

-- Create extensions if needed
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create users table (GORM will handle the actual table creation, but this provides a backup)
-- The actual table structure will be managed by GORM AutoMigrate
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    phone VARCHAR(20),
    is_active BOOLEAN DEFAULT true,
    is_verified BOOLEAN DEFAULT false,
    last_login TIMESTAMP WITH TIME ZONE
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at);
CREATE INDEX IF NOT EXISTS idx_users_is_active ON users(is_active);
CREATE INDEX IF NOT EXISTS idx_users_is_verified ON users(is_verified);

-- Create a function to update the updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create trigger to automatically update updated_at
CREATE TRIGGER update_users_updated_at 
    BEFORE UPDATE ON users 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Insert some sample data (optional - remove in production)
-- INSERT INTO users (email, password_hash, first_name, last_name, is_active, is_verified) 
-- VALUES 
--     ('admin@carbonclear.com', '$2a$10$example_hash', 'Admin', 'User', true, true),
--     ('user@carbonclear.com', '$2a$10$example_hash', 'Regular', 'User', true, false);

-- Grant permissions
GRANT ALL PRIVILEGES ON DATABASE carbon_clear_users TO postgres;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO postgres;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO postgres;

-- Print completion message
DO $$
BEGIN
    RAISE NOTICE 'PostgreSQL initialization completed for user service';
END $$;
