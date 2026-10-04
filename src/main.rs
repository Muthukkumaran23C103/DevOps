use musicdb::musicbrainz::{MusicBrainzClient, MusicBrainzConfig};
use tracing::info;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    tracing_subscriber::fmt::init();

    info!("Starting Indian Music Catalog & Rating Engine Service");

    let mb_config = MusicBrainzConfig::default();
    let mb_client = MusicBrainzClient::new(mb_config)?;

    info!(
        "MusicBrainz API rate-limited client initialized (1 req/sec limit compliant)"
    );

    // Placeholder check
    let _ = mb_client;

    Ok(())
}
