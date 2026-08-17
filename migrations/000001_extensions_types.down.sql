-- Filename: 000001_extensions_types.down.sql
BEGIN;

DROP TYPE IF EXISTS job_status;

DROP TYPE IF EXISTS key_status;

DROP TYPE IF EXISTS consumer_status;

DROP EXTENSION IF EXISTS pgcrypto;

DROP EXTENSION IF EXISTS citext;

COMMIT;