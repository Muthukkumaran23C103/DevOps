use reqwest::Client as HttpClient;
use serde::Deserialize;
use std::sync::Arc;
use std::time::Duration;
use thiserror::Error;
use tokio::sync::Mutex;
use tokio::time::{sleep, Instant};

#[derive(Error, Debug)]
pub enum MusicBrainzError {
    #[error("HTTP client error: {0}")]
    Http(#[from] reqwest::Error),
    #[error("Rate limit timeout error")]
    Timeout,
    #[error("API returned non-200 status code: {0}")]
    ApiStatus(u16),
}

#[derive(Debug, Clone)]
pub struct MusicBrainzConfig {
    pub app_name: String,
    pub version: String,
    pub contact_info: String,
}

impl Default for MusicBrainzConfig {
    fn default() -> Self {
        Self {
            app_name: "IndianMusicDB".to_string(),
            version: "0.1.0".to_string(),
            contact_info: "admin@example.com".to_string(),
        }
    }
}

pub struct MusicBrainzClient {
    http_client: HttpClient,
    base_url: String,
    user_agent: String,
    last_request: Arc<Mutex<Option<Instant>>>,
    min_interval: Duration,
}

#[derive(Debug, Deserialize)]
pub struct MBArtistResponse {
    pub id: String,
    pub name: String,
    #[serde(rename = "sort-name")]
    pub sort_name: Option<String>,
    #[serde(rename = "type")]
    pub artist_type: Option<String>,
    pub country: Option<String>,
}

impl MusicBrainzClient {
    pub fn new(config: MusicBrainzConfig) -> Result<Self, MusicBrainzError> {
        let user_agent = format!(
            "{}/{} ( {} )",
            config.app_name, config.version, config.contact_info
        );

        let http_client = HttpClient::builder()
            .timeout(Duration::from_secs(15))
            .build()?;

        Ok(Self {
            http_client,
            base_url: "https://musicbrainz.org/ws/2".to_string(),
            user_agent,
            last_request: Arc::new(Mutex::new(None)),
            min_interval: Duration::from_secs(1), // Strict 1 req/sec compliance
        })
    }

    async fn enforce_rate_limit(&self) {
        let mut last_req = self.last_request.lock().await;
        if let Some(prev) = *last_req {
            let elapsed = prev.elapsed();
            if elapsed < self.min_interval {
                sleep(self.min_interval - elapsed).await;
            }
        }
        *last_req = Some(Instant::now());
    }

    pub async fn fetch_artist(&self, mbid: &str) -> Result<MBArtistResponse, MusicBrainzError> {
        self.enforce_rate_limit().await;

        let url = format!("{}/artist/{}?fmt=json", self.base_url, mbid);
        let resp = self
            .http_client
            .get(&url)
            .header("User-Agent", &self.user_agent)
            .header("Accept", "application/json")
            .send()
            .await?;

        if !resp.status().is_success() {
            return Err(MusicBrainzError::ApiStatus(resp.status().as_u16()));
        }

        let artist: MBArtistResponse = resp.json().await?;
        Ok(artist)
    }
}
