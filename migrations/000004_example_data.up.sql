BEGIN;

-- ====================================================================================
-- CONSUMERS
-- ====================================================================================
INSERT INTO
    consumers (
        id,
        name,
        email,
        status,
        version,
        created_at,
        updated_at
    )
VALUES
    (
        '0198f000-0000-7000-8000-000000000001',
        'Belize City Council',
        'api@belizecity.example',
        'active',
        1,
        now() - interval '90 days',
        now() - interval '2 days'
    ),
    (
        '0198f000-0000-7000-8000-000000000002',
        'National Meteorological Service',
        'api@hydromet.example',
        'active',
        1,
        now() - interval '60 days',
        now() - interval '1 day'
    ),
    (
        '0198f000-0000-7000-8000-000000000003',
        'Legacy Test Consumer',
        'legacy@example.test',
        'suspended',
        1,
        now() - interval '120 days',
        now() - interval '30 days'
    );

-- ====================================================================================
-- API KEYS
-- ====================================================================================
-- These are deliberately non-secret demonstration hashes, not usable API keys.
INSERT INTO
    api_keys (
        id,
        consumer_id,
        key_hash,
        key_prefix,
        status,
        last_used_at,
        expires_at,
        created_at
    )
VALUES
    (
        '0198f000-0000-7000-8000-000000000101',
        '0198f000-0000-7000-8000-000000000001',
        'sample_hash_city_active',
        'gk_city_',
        'active',
        now() - interval '2 hours',
        now() + interval '180 days',
        now() - interval '90 days'
    ),
    (
        '0198f000-0000-7000-8000-000000000102',
        '0198f000-0000-7000-8000-000000000001',
        'sample_hash_city_revoked',
        'gk_old_',
        'revoked',
        now() - interval '45 days',
        now() - interval '30 days',
        now() - interval '120 days'
    ),
    (
        '0198f000-0000-7000-8000-000000000103',
        '0198f000-0000-7000-8000-000000000002',
        'sample_hash_met_active',
        'gk_met__',
        'active',
        now() - interval '15 minutes',
        now() + interval '180 days',
        now() - interval '60 days'
    ),
    (
        '0198f000-0000-7000-8000-000000000104',
        '0198f000-0000-7000-8000-000000000003',
        'sample_hash_legacy_revoked',
        'gk_leg__',
        'revoked',
        now() - interval '35 days',
        now() - interval '5 days',
        now() - interval '120 days'
    );

-- ====================================================================================
-- JOBS
-- ====================================================================================
INSERT INTO
    jobs (
        id,
        public_id,
        consumer_id,
        job_type,
        status,
        payload,
        result,
        error_message,
        started_at,
        completed_at,
        created_at
    )
VALUES
    (
        '0198f000-0000-7000-8000-000000000201',
        '1198f000-0000-4000-8000-000000000201',
        '0198f000-0000-7000-8000-000000000001',
        'data_export',
        'completed',
        '{"format":"csv"}',
        '{"rows":1250}',
        NULL,
        now() - interval '3 days 2 minutes',
        now() - interval '3 days',
        now() - interval '3 days 3 minutes'
    ),
    (
        '0198f000-0000-7000-8000-000000000202',
        '1198f000-0000-4000-8000-000000000202',
        '0198f000-0000-7000-8000-000000000001',
        'data_export',
        'failed',
        '{"format":"json"}',
        NULL,
        'source data unavailable',
        now() - interval '1 day 1 minute',
        now() - interval '1 day',
        now() - interval '1 day 2 minutes'
    ),
    (
        '0198f000-0000-7000-8000-000000000203',
        '1198f000-0000-4000-8000-000000000203',
        '0198f000-0000-7000-8000-000000000002',
        'data_export',
        'completed',
        '{"format":"csv"}',
        '{"rows":480}',
        NULL,
        now() - interval '8 hours 1 minute',
        now() - interval '8 hours',
        now() - interval '8 hours 2 minutes'
    ),
    (
        '0198f000-0000-7000-8000-000000000204',
        '1198f000-0000-4000-8000-000000000204',
        '0198f000-0000-7000-8000-000000000002',
        'data_export',
        'processing',
        '{"format":"csv"}',
        NULL,
        NULL,
        now() - interval '1 minute',
        NULL,
        now() - interval '2 minutes'
    );

COMMIT;