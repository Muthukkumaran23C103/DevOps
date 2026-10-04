-- 001_init_extensions_and_schemas.sql
-- Enables required PostgreSQL extensions and creates isolated schemas for catalog vs app user data.

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";
CREATE EXTENSION IF NOT EXISTS "unaccent";

-- Schema for public music metadata (artists, works, recordings, releases, credits)
CREATE SCHEMA IF NOT EXISTS catalog;

-- Schema for user-generated application data (users, ratings, reviews, lists)
CREATE SCHEMA IF NOT EXISTS app;
