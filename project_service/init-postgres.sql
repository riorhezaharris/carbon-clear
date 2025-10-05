-- PostgreSQL initialization script for project service
-- This script sets up the database for the project service

-- Create the database if it doesn't exist (handled by POSTGRES_DB env var)
-- CREATE DATABASE carbon_clear_projects;

-- Connect to the database
\c carbon_clear_projects;

-- Create extensions if needed
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create projects table (GORM will handle the actual table creation, but this provides a backup)
-- The actual table structure will be managed by GORM AutoMigrate
CREATE TABLE IF NOT EXISTS projects (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(100),
    region VARCHAR(100),
    country VARCHAR(100),
    verification_standard VARCHAR(100),
    price_per_tonne DECIMAL(10,2),
    total_capacity DECIMAL(15,2),
    available_capacity DECIMAL(15,2),
    project_developer VARCHAR(255),
    project_url VARCHAR(500),
    image_url VARCHAR(500),
    status VARCHAR(50) DEFAULT 'active'
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_projects_title ON projects(title);
CREATE INDEX IF NOT EXISTS idx_projects_category ON projects(category);
CREATE INDEX IF NOT EXISTS idx_projects_region ON projects(region);
CREATE INDEX IF NOT EXISTS idx_projects_country ON projects(country);
CREATE INDEX IF NOT EXISTS idx_projects_status ON projects(status);
CREATE INDEX IF NOT EXISTS idx_projects_price_per_tonne ON projects(price_per_tonne);
CREATE INDEX IF NOT EXISTS idx_projects_created_at ON projects(created_at);
CREATE INDEX IF NOT EXISTS idx_projects_verification_standard ON projects(verification_standard);

-- Create a function to update the updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create trigger to automatically update updated_at
CREATE TRIGGER update_projects_updated_at 
    BEFORE UPDATE ON projects 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Insert some sample data (optional - remove in production)
-- INSERT INTO projects (title, description, category, region, country, verification_standard, price_per_tonne, total_capacity, available_capacity, project_developer, status) 
-- VALUES 
--     ('Solar Farm Project', 'Large-scale solar energy project in California', 'Renewable Energy', 'North America', 'United States', 'VCS', 25.50, 10000.00, 5000.00, 'Green Energy Corp', 'active'),
--     ('Wind Farm Initiative', 'Offshore wind energy project', 'Renewable Energy', 'Europe', 'Germany', 'Gold Standard', 30.75, 15000.00, 8000.00, 'Wind Power Ltd', 'active'),
--     ('Forest Conservation', 'Amazon rainforest protection project', 'Forestry', 'South America', 'Brazil', 'REDD+', 15.25, 5000.00, 3000.00, 'Forest Guardians', 'active');

-- Grant permissions
GRANT ALL PRIVILEGES ON DATABASE carbon_clear_projects TO postgres;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO postgres;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO postgres;

-- Print completion message
DO $$
BEGIN
    RAISE NOTICE 'PostgreSQL initialization completed for project service';
END $$;
