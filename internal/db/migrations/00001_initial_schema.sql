-- +goose Up
CREATE TABLE admins (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email       VARCHAR(255) UNIQUE NOT NULL,
    name        VARCHAR(255),
    deleted_at  TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE subscribers (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email           VARCHAR(255) UNIQUE NOT NULL,
    confirmed_at    TIMESTAMPTZ,
    unsubscribed_at TIMESTAMPTZ,
    deleted_at      TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE venues (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(255) NOT NULL,
    address     TEXT,
    city        VARCHAR(100),
    state       VARCHAR(50),
    capacity    INT,
    deleted_at  TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE bands (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(255) NOT NULL,
    genre       VARCHAR(100),
    description TEXT,
    logo_url    TEXT,
    website_url TEXT,
    deleted_at  TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE shows (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title        VARCHAR(255) NOT NULL,
    slug         VARCHAR(255) UNIQUE NOT NULL,
    date         TIMESTAMPTZ NOT NULL,
    venue_id     UUID REFERENCES venues(id),
    description  TEXT,
    flyer_url    TEXT,
    ticket_url   TEXT,
    is_published BOOLEAN NOT NULL DEFAULT false,
    published_at TIMESTAMPTZ,
    created_by   UUID REFERENCES admins(id),
    deleted_at   TIMESTAMPTZ,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE show_bands (
    show_id       UUID REFERENCES shows(id) ON DELETE CASCADE,
    band_id       UUID REFERENCES bands(id) ON DELETE CASCADE,
    billing_order INT NOT NULL DEFAULT 0,
    is_headliner  BOOLEAN NOT NULL DEFAULT false,
    PRIMARY KEY   (show_id, band_id)
);

CREATE TABLE notifications (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    show_id       UUID REFERENCES shows(id),
    subscriber_id UUID REFERENCES subscribers(id),
    type          VARCHAR(50) NOT NULL, -- 'published' or 'reminder'
    sent_at       TIMESTAMPTZ,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS notifications;
DROP TABLE IF EXISTS show_bands;
DROP TABLE IF EXISTS shows;
DROP TABLE IF EXISTS bands;
DROP TABLE IF EXISTS venues;
DROP TABLE IF EXISTS subscribers;
DROP TABLE IF EXISTS admins;
