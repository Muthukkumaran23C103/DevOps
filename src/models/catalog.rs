use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use uuid::Uuid;

/// Artist entity mapping to catalog.artists
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Artist {
    pub id: Uuid,
    pub mbid: Option<Uuid>,
    pub name: String,
    pub romanized_name: Option<String>,
    pub native_names: HashMap<String, String>,
    pub aliases: Vec<String>,
    pub artist_type: String,
    pub bio: Option<String>,
    pub source: String,
    pub confidence: f64,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
}

/// Abstract musical composition (e.g. Tyagaraja kriti, Ghazal, film composition)
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Work {
    pub id: Uuid,
    pub mbid: Option<Uuid>,
    pub title: String,
    pub romanized_title: Option<String>,
    pub native_titles: HashMap<String, String>,
    pub form: Option<String>,
    pub ragam: Option<String>,
    pub talam: Option<String>,
    pub language: Option<String>,
    pub year_composed: Option<i32>,
    pub source: String,
    pub confidence: f64,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
}

/// Performance audio capture mapping to catalog.recordings
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Recording {
    pub id: Uuid,
    pub mbid: Option<Uuid>,
    pub work_id: Option<Uuid>,
    pub title: String,
    pub isrc: Option<String>,
    pub length_ms: Option<i32>,
    pub audio_type: String,
    pub release_date: Option<String>,
    pub source: String,
    pub confidence: f64,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
}

/// Album, EP, Single, or Film soundtrack container mapping to catalog.releases
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Release {
    pub id: Uuid,
    pub mbid: Option<Uuid>,
    pub title: String,
    pub romanized_title: Option<String>,
    pub native_titles: HashMap<String, String>,
    pub release_type: String,
    pub release_date: Option<String>,
    pub label: Option<String>,
    pub cover_art_url: Option<String>,
    pub film_details: serde_json::Value,
    pub source: String,
    pub confidence: f64,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
}

/// Granular credit attribution connecting artists to works, recordings, or releases
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Credit {
    pub id: Uuid,
    pub entity_type: String,
    pub entity_id: Uuid,
    pub artist_id: Uuid,
    pub role: String,
    pub credit_notes: Option<String>,
    pub source: String,
    pub confidence: f64,
    pub created_at: DateTime<Utc>,
}

/// Streaming service link verified via ISRC or metadata matching
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct StreamingLink {
    pub id: Uuid,
    pub recording_id: Option<Uuid>,
    pub release_id: Option<Uuid>,
    pub platform: String,
    pub external_id: String,
    pub url: String,
    pub isrc_matched: bool,
    pub created_at: DateTime<Utc>,
}
