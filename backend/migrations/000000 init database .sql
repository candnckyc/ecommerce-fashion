-- Create database if it doesn't exist
-- Run this as postgres superuser

-- Connect to default postgres database first
\c postgres

-- Create database
CREATE DATABASE ecommerce;

-- Create user if doesn't exist
DO
$$
BEGIN
   IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'ecommerce') THEN
      CREATE USER ecommerce WITH PASSWORD 'ecommerce123';
   END IF;
END
$$;

-- Grant privileges
GRANT ALL PRIVILEGES ON DATABASE ecommerce TO ecommerce;

-- Connect to the new database
\c ecommerce

-- Grant schema privileges
GRANT ALL ON SCHEMA public TO ecommerce;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO ecommerce;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO ecommerce;

-- Enable full-text search extension
CREATE EXTENSION IF NOT EXISTS pg_trgm;

\echo 'Database setup complete!'
\echo 'Run migrations with: migrate -path migrations -database "postgresql://ecommerce:ecommerce123@localhost:5432/ecommerce?sslmode=disable" up'