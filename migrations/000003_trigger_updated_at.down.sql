-- Filename: 000003_trigger_updated_at.down.sql
BEGIN;

DROP TRIGGER IF EXISTS jobs_updated_at ON jobs;

DROP TRIGGER IF EXISTS consumers_updated_at ON consumers;

DROP FUNCTION IF EXISTS set_updated_at ();

COMMIT;