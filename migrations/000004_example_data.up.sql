BEGIN;

-- ====================================================================================
-- CONSUMERS
-- ====================================================================================
INSERT INTO
    consumers (id, name, email, status)
VALUES
    (
        '018f3a00-0000-7000-8000-000000000001',
        'Acme Corp',
        'api@acme.com',
        'active'
    ),
    (
        '018f3a00-0000-7000-8000-000000000002',
        'Starlight Logistics',
        'devs@starlight.io',
        'active'
    ),
    (
        '018f3a00-0000-7000-8000-000000000003',
        'Legacy Systems Inc',
        'admin@legacysys.com',
        'suspended'
    );

-- ====================================================================================
-- API KEYS
-- Raw keys for development testing:
-- 1. gk_live_a1b2c3d4e5f6g7h8i9j0
-- 2. gk_test_x9y8z7w6v5u4t3s2r1q0
-- 3. gk_live_m1n2o3p4q5r6s7t8u9v0
-- ====================================================================================
INSERT INTO
    api_keys (
        id,
        consumer_id,
        key_hash,
        key_prefix,
        status,
        last_used_at,
        expires_at
    )
VALUES
    (
        '018f3a00-0000-7000-8000-000000000101',
        '018f3a00-0000-7000-8000-000000000001',
        encode(
            digest ('gk_a1b2c3d4e5f6g7h8i9j0', 'sha256'),
            'hex'
        ),
        'gk_a1b2',
        'active',
        NOW() - INTERVAL '10 minutes',
        NOW() + INTERVAL '1 year'
    ),
    (
        '018f3a00-0000-7000-8000-000000000102',
        '018f3a00-0000-7000-8000-000000000001',
        encode(
            digest ('gk_x9y8z7w6v5u4t3s2r1q0', 'sha256'),
            'hex'
        ),
        'gk_x9y8',
        'active',
        NULL,
        NOW() + INTERVAL '30 days'
    ),
    (
        '018f3a00-0000-7000-8000-000000000103',
        '018f3a00-0000-7000-8000-000000000002',
        encode(
            digest ('gk_m1n2o3p4q5r6s7t8u9v0', 'sha256'),
            'hex'
        ),
        'gk_m1n2',
        'active',
        NOW() - INTERVAL '1 hour',
        NULL
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
        completed_at
    )
VALUES
    (
        '018f3a00-0000-7000-8000-000000000201',
        'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11',
        '018f3a00-0000-7000-8000-000000000001',
        'export_csv',
        'completed',
        '{"entity": "orders", "date_range": "2026-Q1"}',
        '{"download_url": "https://storage.example.com/exports/orders_q1.csv", "row_count": 1420}',
        NULL,
        NOW() - INTERVAL '2 hours',
        NOW() - INTERVAL '1 hour 58 minutes'
    ),
    (
        '018f3a00-0000-7000-8000-000000000202',
        'b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a22',
        '018f3a00-0000-7000-8000-000000000001',
        'send_webhook',
        'queued',
        '{"event": "invoice.payment_succeeded", "target_url": "https://acme.com/webhooks"}',
        NULL,
        NULL,
        NULL,
        NULL
    ),
    (
        '018f3a00-0000-7000-8000-000000000203',
        'c2eebc99-9c0b-4ef8-bb6d-6bb9bd380a33',
        '018f3a00-0000-7000-8000-000000000002',
        'process_image',
        'failed',
        '{"image_id": "img_9876", "operations": ["resize", "compress"]}',
        NULL,
        'failed to fetch source image: HTTP 404 Not Found',
        NOW() - INTERVAL '30 minutes',
        NOW() - INTERVAL '29 minutes'
    );

COMMIT;