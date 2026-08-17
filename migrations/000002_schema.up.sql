-- Filename: 000002_schema.up.sql
BEGIN;

-- ====================================================================================
-- CONSUMERS
-- ====================================================================================
CREATE TABLE IF NOT EXISTS
    consumers (
        id UUID PRIMARY KEY DEFAULT uuidv7 (),
        name TEXT NOT NULL,
        email CITEXT NOT NULL UNIQUE,
        status consumer_status NOT NULL DEFAULT 'active',
        version INTEGER NOT NULL DEFAULT 1,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );

-- ====================================================================================
-- API KEYS
-- ====================================================================================
CREATE TABLE IF NOT EXISTS
    api_keys (
        id UUID PRIMARY KEY DEFAULT uuidv7 (),
        consumer_id UUID NOT NULL REFERENCES consumers (id) ON DELETE CASCADE,
        key_hash TEXT NOT NULL UNIQUE,
        key_prefix TEXT NOT NULL,
        status key_status NOT NULL DEFAULT 'active',
        last_used_at TIMESTAMPTZ,
        expires_at TIMESTAMPTZ,
        created_at TIMESTAMPTZ NOT NULL DEFAULT now()
    );

-- ====================================================================================
-- JOBS
-- ====================================================================================
CREATE TABLE IF NOT EXISTS
    jobs (
        id UUID PRIMARY KEY DEFAULT uuidv7 (),
        public_id UUID NOT NULL DEFAULT uuidv4 (),
        consumer_id UUID NOT NULL REFERENCES consumers (id) ON DELETE CASCADE,
        job_type TEXT NOT NULL,
        status job_status NOT NULL DEFAULT 'queued',
        payload JSONB NOT NULL DEFAULT '{}',
        result JSONB,
        error_message TEXT,
        started_at TIMESTAMPTZ,
        completed_at TIMESTAMPTZ,
        created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );

COMMIT;