-- Filename: 000001_extensions_types.up.sql
BEGIN;

CREATE EXTENSION IF NOT EXISTS citext;

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TYPE consumer_status AS ENUM('active', 'suspended', 'terminated');

CREATE TYPE key_status AS ENUM('active', 'rotating', 'revoked');

CREATE TYPE job_status AS ENUM(
    'queued',
    'processing',
    'completed',
    'failed',
    'cancelled'
);

COMMIT;