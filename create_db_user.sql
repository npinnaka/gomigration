-- Create the database
CREATE DATABASE mydb;

-- Create the user
CREATE USER dbuser WITH PASSWORD 'password';

-- Grant all privileges on the database to the user
GRANT ALL PRIVILEGES ON DATABASE mydb TO dbuser;


\c mydb

-- Grant all privileges on all tables in the public schema (optional, if needed)
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO dbuser;
