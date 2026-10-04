-- 002_create_catalog_tables.sql
-- Catalog schema tables: Artists, Works, Recordings, Releases, Credits, Genres, Streaming Links.
-- All metadata fields support LLM/external data auditability via 'source' and 'confidence' columns.

CREATE TABLE IF NOT EXISTS catalog.artists (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    mbid UUID UNIQUE,
    name TEXT NOT NULL,
    romanized_name TEXT,
    native_names JSONB DEFAULT '{}'::jsonb,
    aliases TEXT[] DEFAULT '{}',
    artist_type TEXT DEFAULT 'person',
    bio TEXT,
    source TEXT NOT NULL DEFAULT 'manual',
    confidence NUMERIC(3,2) NOT NULL DEFAULT 1.00 CHECK (confidence >= 0.00 AND confidence <= 1.00),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS catalog.works (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    mbid UUID UNIQUE,
    title TEXT NOT NULL,
    romanized_title TEXT,
    native_titles JSONB DEFAULT '{}'::jsonb,
    form TEXT, -- e.g., 'kriti', 'varnam', 'thillana', 'ghazal', 'film_song'
    ragam TEXT, -- e.g., 'Mayamalavagowla', 'Bhairavi', 'Yaman'
    talam TEXT, -- e.g., 'Adi', 'Rupakam', 'Teental'
    language TEXT,
    year_composed INT,
    source TEXT NOT NULL DEFAULT 'manual',
    confidence NUMERIC(3,2) NOT NULL DEFAULT 1.00 CHECK (confidence >= 0.00 AND confidence <= 1.00),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS catalog.recordings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    mbid UUID UNIQUE,
    work_id UUID REFERENCES catalog.works(id) ON DELETE SET NULL,
    title TEXT NOT NULL,
    isrc TEXT,
    length_ms INT,
    audio_type TEXT DEFAULT 'studio', -- 'studio', 'live', 'field', 'broadcast'
    release_date DATE,
    source TEXT NOT NULL DEFAULT 'manual',
    confidence NUMERIC(3,2) NOT NULL DEFAULT 1.00 CHECK (confidence >= 0.00 AND confidence <= 1.00),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS catalog.releases (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    mbid UUID UNIQUE,
    title TEXT NOT NULL,
    romanized_title TEXT,
    native_titles JSONB DEFAULT '{}'::jsonb,
    release_type TEXT NOT NULL DEFAULT 'album', -- 'album', 'film', 'soundtrack', 'ep', 'single', 'compilation'
    release_date DATE,
    label TEXT,
    cover_art_url TEXT,
    film_details JSONB DEFAULT '{}'::jsonb, -- e.g. {"director": "Mani Ratnam", "language": "ta", "year": 1991}
    source TEXT NOT NULL DEFAULT 'manual',
    confidence NUMERIC(3,2) NOT NULL DEFAULT 1.00 CHECK (confidence >= 0.00 AND confidence <= 1.00),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS catalog.release_tracks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    release_id UUID NOT NULL REFERENCES catalog.releases(id) ON DELETE CASCADE,
    recording_id UUID NOT NULL REFERENCES catalog.recordings(id) ON DELETE RESTRICT,
    disc_number INT NOT NULL DEFAULT 1,
    track_number INT NOT NULL,
    track_title TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (release_id, disc_number, track_number)
);

CREATE TABLE IF NOT EXISTS catalog.credits (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    entity_type TEXT NOT NULL CHECK (entity_type IN ('work', 'recording', 'release')),
    entity_id UUID NOT NULL,
    artist_id UUID NOT NULL REFERENCES catalog.artists(id) ON DELETE CASCADE,
    role TEXT NOT NULL, -- 'composer', 'lyricist', 'music_director', 'vocalist', 'mridangam', 'violin', 'tabla', 'producer', etc.
    credit_notes TEXT,
    source TEXT NOT NULL DEFAULT 'manual',
    confidence NUMERIC(3,2) NOT NULL DEFAULT 1.00 CHECK (confidence >= 0.00 AND confidence <= 1.00),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS catalog.genres (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT UNIQUE NOT NULL,
    parent_id UUID REFERENCES catalog.genres(id) ON DELETE SET NULL,
    description TEXT
);

CREATE TABLE IF NOT EXISTS catalog.entity_genres (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    entity_type TEXT NOT NULL CHECK (entity_type IN ('work', 'recording', 'release', 'artist')),
    entity_id UUID NOT NULL,
    genre_id UUID NOT NULL REFERENCES catalog.genres(id) ON DELETE CASCADE,
    source TEXT NOT NULL DEFAULT 'manual',
    confidence NUMERIC(3,2) NOT NULL DEFAULT 1.00 CHECK (confidence >= 0.00 AND confidence <= 1.00),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (entity_type, entity_id, genre_id)
);

CREATE TABLE IF NOT EXISTS catalog.streaming_links (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    recording_id UUID REFERENCES catalog.recordings(id) ON DELETE CASCADE,
    release_id UUID REFERENCES catalog.releases(id) ON DELETE CASCADE,
    platform TEXT NOT NULL CHECK (platform IN ('spotify', 'apple_music', 'youtube_music', 'bandcamp', 'archive_org')),
    external_id TEXT NOT NULL,
    url TEXT NOT NULL,
    isrc_matched BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT stream_link_target CHECK (recording_id IS NOT NULL OR release_id IS NOT NULL)
);
