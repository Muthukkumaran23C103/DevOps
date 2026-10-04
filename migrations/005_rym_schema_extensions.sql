-- 005_rym_schema_extensions.sql
-- RYM / Sonemic feature extensions: Descriptors, Track Ratings, User Collections, Artist Images.

-- Add artist image URL support
ALTER TABLE catalog.artists ADD COLUMN IF NOT EXISTS image_url TEXT;

-- Mood and music style descriptors (e.g., 'melancholic', 'improvisational', 'polyrhythmic', 'poetic', 'cinematic')
CREATE TABLE IF NOT EXISTS catalog.descriptors (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT UNIQUE NOT NULL,
    category TEXT DEFAULT 'mood' -- 'mood', 'style', 'theme', 'structure'
);

CREATE TABLE IF NOT EXISTS catalog.release_descriptors (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    release_id UUID NOT NULL REFERENCES catalog.releases(id) ON DELETE CASCADE,
    descriptor_id UUID NOT NULL REFERENCES catalog.descriptors(id) ON DELETE CASCADE,
    source TEXT NOT NULL DEFAULT 'manual',
    confidence NUMERIC(3,2) NOT NULL DEFAULT 1.00 CHECK (confidence >= 0.00 AND confidence <= 1.00),
    UNIQUE (release_id, descriptor_id)
);

CREATE TABLE IF NOT EXISTS catalog.work_descriptors (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    work_id UUID NOT NULL REFERENCES catalog.works(id) ON DELETE CASCADE,
    descriptor_id UUID NOT NULL REFERENCES catalog.descriptors(id) ON DELETE CASCADE,
    source TEXT NOT NULL DEFAULT 'manual',
    confidence NUMERIC(3,2) NOT NULL DEFAULT 1.00 CHECK (confidence >= 0.00 AND confidence <= 1.00),
    UNIQUE (work_id, descriptor_id)
);

-- Per-track individual ratings & favorite track highlights
CREATE TABLE IF NOT EXISTS app.track_ratings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES app.users(id) ON DELETE CASCADE,
    recording_id UUID NOT NULL REFERENCES catalog.recordings(id) ON DELETE CASCADE,
    rating NUMERIC(2,1) CHECK (rating >= 0.5 AND rating <= 5.0),
    is_favorite BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, recording_id)
);

-- User Collection tracking (Owned format, Wishlist, Listened)
CREATE TABLE IF NOT EXISTS app.user_collection (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES app.users(id) ON DELETE CASCADE,
    release_id UUID NOT NULL REFERENCES catalog.releases(id) ON DELETE CASCADE,
    status TEXT NOT NULL CHECK (status IN ('owned', 'wishlist', 'listened', 'cataloged')),
    format TEXT DEFAULT 'digital', -- 'vinyl', 'cd', 'cassette', 'digital'
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, release_id)
);
