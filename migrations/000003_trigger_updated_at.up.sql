-- Filename: 000003_trigger_updated_at.up.sql
BEGIN;

CREATE
OR REPLACE FUNCTION set_updated_at () RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Consumers
CREATE TRIGGER consumers_updated_at BEFORE
UPDATE ON consumers FOR EACH ROW
EXECUTE FUNCTION set_updated_at ();

-- Jobs
CREATE TRIGGER jobs_updated_at BEFORE
UPDATE ON jobs FOR EACH ROW
EXECUTE FUNCTION set_updated_at ();

COMMIT;