package catalog

import (
	"encoding/json"
	"time"
)

// Artist represents a catalog musical entity with image URL support.
type Artist struct {
	ID            string            `json:"id" db:"id"`
	MBID          *string           `json:"mbid,omitempty" db:"mbid"`
	Name          string            `json:"name" db:"name"`
	RomanizedName *string           `json:"romanized_name,omitempty" db:"romanized_name"`
	NativeNames   map[string]string `json:"native_names,omitempty" db:"native_names"`
	Aliases       []string          `json:"aliases,omitempty" db:"aliases"`
	ArtistType    string            `json:"artist_type" db:"artist_type"`
	Bio           *string           `json:"bio,omitempty" db:"bio"`
	ImageURL      *string           `json:"image_url,omitempty" db:"image_url"`
	Source        string            `json:"source" db:"source"`
	Confidence    float64           `json:"confidence" db:"confidence"`
	CreatedAt     time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at" db:"updated_at"`
}

// Work represents a distinct composition (e.g. Tyagaraja kriti, Ghazal, film composition).
type Work struct {
	ID             string            `json:"id" db:"id"`
	MBID           *string           `json:"mbid,omitempty" db:"mbid"`
	Title          string            `json:"title" db:"title"`
	RomanizedTitle *string           `json:"romanized_title,omitempty" db:"romanized_title"`
	NativeTitles   map[string]string `json:"native_titles,omitempty" db:"native_titles"`
	Form           *string           `json:"form,omitempty" db:"form"`         // 'kriti', 'varnam', 'film_song', etc.
	Ragam          *string           `json:"ragam,omitempty" db:"ragam"`       // Raga name
	Talam          *string           `json:"talam,omitempty" db:"talam"`       // Rhythm cycle
	Language       *string           `json:"language,omitempty" db:"language"`
	YearComposed   *int              `json:"year_composed,omitempty" db:"year_composed"`
	Source         string            `json:"source" db:"source"`
	Confidence     float64           `json:"confidence" db:"confidence"`
	CreatedAt      time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at" db:"updated_at"`
}

// Recording represents an audio capture/performance of a work.
type Recording struct {
	ID          string    `json:"id" db:"id"`
	MBID        *string   `json:"mbid,omitempty" db:"mbid"`
	WorkID      *string   `json:"work_id,omitempty" db:"work_id"`
	Title       string    `json:"title" db:"title"`
	ISRC        *string   `json:"isrc,omitempty" db:"isrc"`
	LengthMS    *int      `json:"length_ms,omitempty" db:"length_ms"`
	AudioType   string    `json:"audio_type" db:"audio_type"` // 'studio', 'live', 'field'
	ReleaseDate *string   `json:"release_date,omitempty" db:"release_date"`
	Source      string    `json:"source" db:"source"`
	Confidence  float64   `json:"confidence" db:"confidence"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Release represents an album, EP, single, or film soundtrack release container.
type Release struct {
	ID             string            `json:"id" db:"id"`
	MBID           *string           `json:"mbid,omitempty" db:"mbid"`
	Title          string            `json:"title" db:"title"`
	RomanizedTitle *string           `json:"romanized_title,omitempty" db:"romanized_title"`
	NativeTitles   map[string]string `json:"native_titles,omitempty" db:"native_titles"`
	ReleaseType    string            `json:"release_type" db:"release_type"` // 'album', 'film', 'soundtrack'
	ReleaseDate    *string           `json:"release_date,omitempty" db:"release_date"`
	Label          *string           `json:"label,omitempty" db:"label"`
	CoverArtURL    *string           `json:"cover_art_url,omitempty" db:"cover_art_url"`
	FilmDetails    json.RawMessage   `json:"film_details,omitempty" db:"film_details"`
	Source         string            `json:"source" db:"source"`
	Confidence     float64           `json:"confidence" db:"confidence"`
	CreatedAt      time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at" db:"updated_at"`
}

// Credit attributes artists to works, recordings, or releases.
type Credit struct {
	ID          string    `json:"id" db:"id"`
	EntityType  string    `json:"entity_type" db:"entity_type"`
	EntityID    string    `json:"entity_id" db:"entity_id"`
	ArtistID    string    `json:"artist_id" db:"artist_id"`
	Role        string    `json:"role" db:"role"`
	CreditNotes *string   `json:"credit_notes,omitempty" db:"credit_notes"`
	Source      string    `json:"source" db:"source"`
	Confidence  float64   `json:"confidence" db:"confidence"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// Descriptor represents mood and musical descriptors (e.g. melancholic, polyrhythmic).
type Descriptor struct {
	ID       string `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`
	Category string `json:"category" db:"category"`
}

// TrackRating represents user track-level ratings.
type TrackRating struct {
	ID          string    `json:"id" db:"id"`
	UserID      string    `json:"user_id" db:"user_id"`
	RecordingID string    `json:"recording_id" db:"recording_id"`
	Rating      *float64  `json:"rating,omitempty" db:"rating"`
	IsFavorite  bool      `json:"is_favorite" db:"is_favorite"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// ChartItem represents an entry in RYM dynamic charts with Bayesian score.
type ChartItem struct {
	Rank            int      `json:"rank"`
	Release         Release  `json:"release"`
	Artist          Artist   `json:"artist"`
	WeightedAverage float64  `json:"weighted_average"`
	BayesianScore   float64  `json:"bayesian_score"`
	TotalRatings    int      `json:"total_ratings"`
	TotalReviews    int      `json:"total_reviews"`
	PrimaryGenres   []string `json:"primary_genres"`
	Descriptors     []string `json:"descriptors"`
}

// Rating represents user rating stored isolated in the 'app' schema.
type Rating struct {
	ID         string    `json:"id" db:"id"`
	UserID     string    `json:"user_id" db:"user_id"`
	EntityType string    `json:"entity_type" db:"entity_type"`
	EntityID   string    `json:"entity_id" db:"entity_id"`
	Rating     float64   `json:"rating" db:"rating"` // 0.5 to 5.0
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// Review represents text reviews with spoiler tags.
type Review struct {
	ID         string    `json:"id" db:"id"`
	UserID     string    `json:"user_id" db:"user_id"`
	Username   string    `json:"username,omitempty"`
	EntityType string    `json:"entity_type" db:"entity_type"`
	EntityID   string    `json:"entity_id" db:"entity_id"`
	Content    string    `json:"content" db:"content"`
	IsSpoiler  bool      `json:"is_spoiler" db:"is_spoiler"`
	LikesCount int       `json:"likes_count"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}
