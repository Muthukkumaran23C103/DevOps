-- 004_create_search_indexes.sql
-- Implements trigram search indexes and script normalization functions for multi-script searching
-- (Tamil, Devanagari, Telugu, Kannada, Malayalam, Romanized ASCII).

-- Helper function to convert JSONB native script maps into a single searchable string
CREATE OR REPLACE FUNCTION catalog.native_names_to_text(data JSONB)
RETURNS TEXT IMMUTABLE PARALLEL SAFE AS $$
DECLARE
    val RECORD;
    result TEXT := '';
BEGIN
    IF data IS NULL THEN
        RETURN '';
    END IF;
    FOR val IN SELECT value FROM jsonb_each_text(data) LOOP
        result := result || ' ' || val.value;
    END LOOP;
    RETURN result;
END;
$$ LANGUAGE plpgsql;

-- Immutable text normalization function combining unaccent and lowercasing
CREATE OR REPLACE FUNCTION catalog.search_normalize(txt TEXT)
RETURNS TEXT IMMUTABLE PARALLEL SAFE AS $$
BEGIN
    IF txt IS NULL THEN
        RETURN '';
    END IF;
    RETURN lower(unaccent(txt));
END;
$$ LANGUAGE plpgsql;

-- Combined immutable text generator for artists
CREATE OR REPLACE FUNCTION catalog.artist_search_text(
    name TEXT,
    romanized_name TEXT,
    aliases TEXT[],
    native_names JSONB
) RETURNS TEXT IMMUTABLE PARALLEL SAFE AS $$
BEGIN
    RETURN catalog.search_normalize(name) || ' ' ||
           catalog.search_normalize(romanized_name) || ' ' ||
           array_to_string(aliases, ' ') || ' ' ||
           catalog.native_names_to_text(native_names);
END;
$$ LANGUAGE plpgsql;

-- Combined immutable text generator for works
CREATE OR REPLACE FUNCTION catalog.work_search_text(
    title TEXT,
    romanized_title TEXT,
    native_titles JSONB,
    ragam TEXT
) RETURNS TEXT IMMUTABLE PARALLEL SAFE AS $$
BEGIN
    RETURN catalog.search_normalize(title) || ' ' ||
           catalog.search_normalize(romanized_title) || ' ' ||
           catalog.native_names_to_text(native_titles) || ' ' ||
           catalog.search_normalize(ragam);
END;
$$ LANGUAGE plpgsql;

-- Combined immutable text generator for releases
CREATE OR REPLACE FUNCTION catalog.release_search_text(
    title TEXT,
    romanized_title TEXT,
    native_titles JSONB
) RETURNS TEXT IMMUTABLE PARALLEL SAFE AS $$
BEGIN
    RETURN catalog.search_normalize(title) || ' ' ||
           catalog.search_normalize(romanized_title) || ' ' ||
           catalog.native_names_to_text(native_titles);
END;
$$ LANGUAGE plpgsql;

-- GIN Trigram indexes
CREATE INDEX IF NOT EXISTS idx_artists_search_trgm ON catalog.artists
USING gin (catalog.artist_search_text(name, romanized_name, aliases, native_names) gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_works_search_trgm ON catalog.works
USING gin (catalog.work_search_text(title, romanized_title, native_titles, ragam) gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_releases_search_trgm ON catalog.releases
USING gin (catalog.release_search_text(title, romanized_title, native_titles) gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_recordings_title_trgm ON catalog.recordings
USING gin (catalog.search_normalize(title) gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_recordings_isrc ON catalog.recordings(isrc) WHERE isrc IS NOT NULL;
