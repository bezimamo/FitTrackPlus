-- FitTrack+ Database Setup Script
-- This script creates the database and user for the FitTrack+ application

-- Connect to PostgreSQL as superuser (postgres)
-- Run this script with: psql -U postgres -f setup-database.sql

-- Create the database
CREATE DATABASE fittrackplus;

-- Create a dedicated user for the application
CREATE USER fittrackplus_user WITH PASSWORD 'fittrackplus_password_2024';

-- Grant privileges to the user
GRANT ALL PRIVILEGES ON DATABASE fittrackplus TO fittrackplus_user;

-- Connect to the new database
\c fittrackplus;

-- Grant schema privileges
GRANT ALL ON SCHEMA public TO fittrackplus_user;

-- Create extensions if needed
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Display success message
\echo '✅ FitTrack+ database setup completed successfully!'
\echo '📊 Database: fittrackplus'
\echo '👤 User: fittrackplus_user'
\echo '🔑 Password: fittrackplus_password_2024'
\echo ''
\echo '📝 Next steps:'
\echo '1. Update your .env file with these credentials'
\echo '2. Run: go run cmd/server/main.go'
\echo '3. Access: http://localhost:8080/swagger/index.html' 