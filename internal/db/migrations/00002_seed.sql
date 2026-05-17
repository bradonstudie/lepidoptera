-- +goose Up

-- venues
INSERT INTO venues (id, name, address, city, state)
VALUES
    (gen_random_uuid(), 'The Sidebar', '218 E Lexington St', 'Baltimore', 'MD'),
    (gen_random_uuid(), 'Ottobar', '2549 N Howard St', 'Baltimore', 'MD'),
    (gen_random_uuid(), 'Metro Gallery', '1700 N Charles St', 'Baltimore', 'MD');

-- bands
INSERT INTO bands (id, name, genre, description)
VALUES
    (gen_random_uuid(), 'Turnstile', 'hardcore', 'test description'),
    (gen_random_uuid(), 'Angel Du$t', 'punk', 'test description'),
    (gen_random_uuid(), 'Trapped Under Ice', 'hardcore', 'test description'),
    (gen_random_uuid(), 'Jivebomb', 'punk', 'test description');

-- shows (referencing venues by name for readability)
INSERT INTO shows (id, title, slug, date, venue_id, description, is_published, published_at)
VALUES
    (
        gen_random_uuid(),
        'Turnstile',
        'turnstile-sidebar-2026',
        '2026-06-14 19:00:00-05',
        (SELECT id FROM venues WHERE name = 'The Sidebar'),
        'show description 01',
        true,
        NOW()
    ),
    (
        gen_random_uuid(),
        'Punk Night Vol. 4',
        'punk-night-vol-4-ottobar-2026',
        '2026-07-04 20:00:00-05',
        (SELECT id FROM venues WHERE name = 'Ottobar'),
        'show description 02',
        true,
        NOW()
    ),
    (
        gen_random_uuid(),
        'Trapped Under Ice',
        'trapped-under-ice-metro-2026',
        '2026-08-20 19:30:00-05',
        (SELECT id FROM venues WHERE name = 'Metro Gallery'),
        'show description 03',
        true,
        NOW()
    );

-- show_bands (wire bands to shows)
INSERT INTO show_bands (show_id, band_id, billing_order, is_headliner)
VALUES
    -- Turnstile show
    (
        (SELECT id FROM shows WHERE slug = 'turnstile-sidebar-2026'),
        (SELECT id FROM bands WHERE name = 'Turnstile'),
        0, true
    ),
    (
        (SELECT id FROM shows WHERE slug = 'turnstile-sidebar-2026'),
        (SELECT id FROM bands WHERE name = 'Angel Du$t'),
        1, false
    ),
    -- Punk Night
    (
        (SELECT id FROM shows WHERE slug = 'punk-night-vol-4-ottobar-2026'),
        (SELECT id FROM bands WHERE name = 'Angel Du$t'),
        0, true
    ),
    (
        (SELECT id FROM shows WHERE slug = 'punk-night-vol-4-ottobar-2026'),
        (SELECT id FROM bands WHERE name = 'Jivebomb'),
        1, false
    ),
    -- Trapped Under Ice
    (
        (SELECT id FROM shows WHERE slug = 'trapped-under-ice-metro-2026'),
        (SELECT id FROM bands WHERE name = 'Trapped Under Ice'),
        0, true
    ),
    (
        (SELECT id FROM shows WHERE slug = 'trapped-under-ice-metro-2026'),
        (SELECT id FROM bands WHERE name = 'Jivebomb'),
        1, false
    );

-- +goose Down
DELETE FROM show_bands;
DELETE FROM shows;
DELETE FROM bands;
DELETE FROM venues;