package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"sync"
	"time"

	"devops/musicdb/pkg/catalog"
	"devops/musicdb/pkg/musicbrainz"
)

type Store struct {
	mu      sync.RWMutex
	ratings map[string][]float64 // entity_id -> array of ratings
	reviews map[string][]catalog.Review
	lists   []map[string]interface{}
}

type Server struct {
	mbClient *musicbrainz.Client
	store    *Store
}

func NewStore() *Store {
	return &Store{
		ratings: map[string][]float64{
			"rel-nayakan":        {4.5, 5.0, 4.5, 5.0, 4.0, 4.5, 5.0, 4.5, 4.5, 5.0},
			"rel-thalapathi":     {4.5, 4.5, 5.0, 4.0, 4.5, 4.5, 5.0},
			"rel-roja":           {5.0, 4.5, 4.5, 4.0, 5.0, 4.5},
			"rel-how-to-name-it": {4.5, 4.5, 5.0, 4.5, 4.0},
		},
		reviews: map[string][]catalog.Review{
			"rel-nayakan": {
				{
					ID:         "rev-1",
					UserID:     "usr-1",
					Username:   "CarnaticCrateDigger",
					EntityType: "release",
					EntityID:   "rel-nayakan",
					Content:    "Ilaiyaraaja's score for Nayakan is a monumental masterpiece of Indian cinema. Combining traditional Carnatic melodic structures with lush symphonic string arrangements.",
					IsSpoiler:  false,
					LikesCount: 42,
					CreatedAt:  time.Now().Add(-48 * time.Hour),
				},
			},
		},
		lists: []map[string]interface{}{
			{
				"id":          "list-1",
				"title":       "Top 10 Indian Film Soundtracks of All Time",
				"author":      "IsaiFanatic",
				"description": "Essential film music compositions from South Asia, focusing on Ilaiyaraaja, A.R. Rahman, and R.D. Burman.",
				"item_count":  10,
				"items": []map[string]interface{}{
					{"rank": 1, "release_id": "rel-nayakan", "title": "Nayakan", "artist": "Ilaiyaraaja", "year": 1987, "cover": "assets/images/nayakan.jpg", "remote_cover": "https://is1-ssl.mzstatic.com/image/thumb/Music211/v4/6f/fd/a9/6ffda94e-d67c-00ed-27d2-27c8d2b0e973/196871814317.jpg/600x600bb.jpg"},
					{"rank": 2, "release_id": "rel-thalapathi", "title": "Thalapathi", "artist": "Ilaiyaraaja", "year": 1991, "cover": "assets/images/thalapathi.jpg", "remote_cover": "https://is1-ssl.mzstatic.com/image/thumb/Music211/v4/72/c2/95/72c295f9-ccec-c74e-1ddb-083b649133dd/8905750030999.jpg/600x600bb.jpg"},
					{"rank": 3, "release_id": "rel-roja", "title": "Roja", "artist": "A.R. Rahman", "year": 1992, "cover": "assets/images/roja.jpg", "remote_cover": "https://is1-ssl.mzstatic.com/image/thumb/Music221/v4/09/5d/7b/095d7bd0-001c-695a-4db1-eaa196ed0918/8905750032702.jpg/600x600bb.jpg"},
				},
			},
		},
	}
}

func main() {
	mbClient := musicbrainz.NewClient(musicbrainz.Config{
		AppName:     "IndianMusicDB",
		Version:     "0.1.0",
		ContactInfo: "admin@example.com",
	})

	srv := &Server{
		mbClient: mbClient,
		store:    NewStore(),
	}

	mux := http.NewServeMux()

	// Serve RYM web app frontend on root
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/index.html")
	})

	// Serve static downloaded image assets
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))

	// REST APIs
	mux.HandleFunc("GET /api/v1/releases/{id}", srv.handleGetRelease)
	mux.HandleFunc("GET /api/v1/artists/{id}", srv.handleGetArtist)
	mux.HandleFunc("GET /api/v1/charts", srv.handleGetCharts)
	mux.HandleFunc("GET /api/v1/lists", srv.handleGetLists)
	mux.HandleFunc("POST /api/v1/ratings", srv.handlePostRating)
	mux.HandleFunc("GET /api/v1/reviews", srv.handleGetReviews)
	mux.HandleFunc("POST /api/v1/reviews", srv.handlePostReview)
	mux.HandleFunc("GET /api/v1/search", srv.handleSearch)

	port := ":8080"
	log.Printf("[+] RYM Engine Server listening on http://localhost%s\n", port)
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("[-] Server error: %v", err)
	}
}

func setCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
}

func (s *Server) handleGetRelease(w http.ResponseWriter, r *http.Request) {
	setCORS(w)
	relID := r.PathValue("id")
	if relID == "" {
		relID = "rel-nayakan"
	}

	s.store.mu.RLock()
	ratings := s.store.ratings[relID]
	reviews := s.store.reviews[relID]
	s.store.mu.RUnlock()

	avgRating, totalRatings := calculateStats(ratings)

	payload := map[string]interface{}{
		"id":              relID,
		"title":           "Nayakan",
		"romanized_title": "Nayakan",
		"native_titles": map[string]string{
			"ta": "நாயகன்",
			"hi": "नायगन",
		},
		"release_type":  "film_soundtrack",
		"release_date":  "1987-10-21",
		"label":         "AVM Music",
		"cover_art_url": "assets/images/nayakan.jpg",
		"remote_cover":  "https://is1-ssl.mzstatic.com/image/thumb/Music211/v4/6f/fd/a9/6ffda94e-d67c-00ed-27d2-27c8d2b0e973/196871814317.jpg/600x600bb.jpg",
		"artist": map[string]interface{}{
			"id":             "art-ilaiyaraaja",
			"name":           "Ilaiyaraaja",
			"romanized_name": "Ilaiyaraaja",
			"native_names": map[string]string{
				"ta": "இளையராஜா",
				"hi": "इलैयाराजा",
			},
			"image_url":    "assets/images/ilaiyaraaja.jpg",
			"remote_image": "https://cdn-images.dzcdn.net/images/artist/aeeca2a4b9f808c8e8f5d281b1fb48d0/500x500-000000-80-0-0.jpg",
		},
		"primary_genres":   []string{"Tamil Film Music", "Orchestral Folk"},
		"secondary_genres": []string{"Carnatic Fusion", "Chamber Pop"},
		"descriptors":      []string{"melancholic", "cinematic", "poetic", "orchestral", "nostalgic"},
		"stats": map[string]interface{}{
			"average_rating": avgRating,
			"total_ratings":  totalRatings,
			"total_reviews":  len(reviews),
			"rank_all_time":  "#1 Indian Film Soundtrack",
			"bayesian_score": calculateBayesian(avgRating, totalRatings),
		},
		"tracklist": []map[string]interface{}{
			{"track_number": 1, "title": "Thenpandi Cheemayile", "duration": "4:32", "singers": "Ilaiyaraaja, Kamal Haasan", "ragam": "Mayamalavagowla", "avg_rating": 4.85},
			{"track_number": 2, "title": "Nee Oru Kaadhal Sangeetham", "duration": "4:48", "singers": "K.J. Yesudas, Chithra", "ragam": "Abheri", "avg_rating": 4.90},
			{"track_number": 3, "title": "Naan Sirithal Deepavali", "duration": "4:55", "singers": "M.S. Rajeswari, Jamuna Rani", "ragam": "Bhairavi", "avg_rating": 4.70},
			{"track_number": 4, "title": "Andhi Mazhai Megam", "duration": "4:38", "singers": "T.L. Maharajan, P. Susheela", "ragam": "Kalyani", "avg_rating": 4.75},
		},
		"reviews": reviews,
	}

	json.NewEncoder(w).Encode(payload)
}

func (s *Server) handleGetArtist(w http.ResponseWriter, r *http.Request) {
	setCORS(w)
	artistID := r.PathValue("id")
	if artistID == "" {
		artistID = "art-ilaiyaraaja"
	}

	payload := map[string]interface{}{
		"id":             artistID,
		"name":           "Ilaiyaraaja",
		"romanized_name": "Ilaiyaraaja",
		"native_names": map[string]string{
			"ta": "இளையராஜா",
			"hi": "इलैयाराजा",
			"te": "ఇళయరాజా",
		},
		"aliases":      []string{"Isaignani", "Ilayaraja", "Maestro"},
		"artist_type":  "person",
		"image_url":    "assets/images/ilaiyaraaja.jpg",
		"remote_image": "https://cdn-images.dzcdn.net/images/artist/aeeca2a4b9f808c8e8f5d281b1fb48d0/500x500-000000-80-0-0.jpg",
		"bio":          "Ilaiyaraaja is an Indian film composer, conductor, arranger, lyricist, and singer who has composed over 7,000 songs and film scores for over 1,400 films.",
		"discography": []map[string]interface{}{
			{"id": "rel-nayakan", "title": "Nayakan", "year": 1987, "type": "Film Soundtrack", "rating": 4.78, "ratings_count": 1420, "cover": "assets/images/nayakan.jpg", "remote_cover": "https://is1-ssl.mzstatic.com/image/thumb/Music211/v4/6f/fd/a9/6ffda94e-d67c-00ed-27d2-27c8d2b0e973/196871814317.jpg/600x600bb.jpg"},
			{"id": "rel-thalapathi", "title": "Thalapathi", "year": 1991, "type": "Film Soundtrack", "rating": 4.72, "ratings_count": 1180, "cover": "assets/images/thalapathi.jpg", "remote_cover": "https://is1-ssl.mzstatic.com/image/thumb/Music211/v4/72/c2/95/72c295f9-ccec-c74e-1ddb-083b649133dd/8905750030999.jpg/600x600bb.jpg"},
			{"id": "rel-how-to-name-it", "title": "How to Name It?", "year": 1986, "type": "Studio Album / Fusion", "rating": 4.65, "ratings_count": 890, "cover": "assets/images/how_to_name_it.jpg", "remote_cover": "https://is1-ssl.mzstatic.com/image/thumb/Music221/v4/a3/4d/09/a34d0902-ffc1-c1cd-96fc-8f693bd9de5e/886448525172.jpg/600x600bb.jpg"},
			{"id": "rel-nothing-but-wind", "title": "Nothing But Wind", "year": 1988, "type": "Studio Album / Instrumental", "rating": 4.58, "ratings_count": 650, "cover": "assets/images/nothing_but_wind.jpg", "remote_cover": "https://is1-ssl.mzstatic.com/image/thumb/Music122/v4/d9/a6/d8/d9a6d8b3-0e47-f088-1edd-ea7546007d82/196871878883.jpg/600x600bb.jpg"},
		},
	}

	json.NewEncoder(w).Encode(payload)
}

func (s *Server) handleGetCharts(w http.ResponseWriter, r *http.Request) {
	setCORS(w)

	charts := []map[string]interface{}{
		{
			"rank":           1,
			"id":             "rel-nayakan",
			"title":          "Nayakan",
			"artist":         "Ilaiyaraaja",
			"artist_id":      "art-ilaiyaraaja",
			"artist_image":   "assets/images/ilaiyaraaja.jpg",
			"remote_artist":  "https://cdn-images.dzcdn.net/images/artist/aeeca2a4b9f808c8e8f5d281b1fb48d0/500x500-000000-80-0-0.jpg",
			"release_year":   1987,
			"release_type":   "Film Soundtrack",
			"cover_art_url":  "assets/images/nayakan.jpg",
			"remote_cover":   "https://is1-ssl.mzstatic.com/image/thumb/Music211/v4/6f/fd/a9/6ffda94e-d67c-00ed-27d2-27c8d2b0e973/196871814317.jpg/600x600bb.jpg",
			"average_rating": 4.78,
			"bayesian_score": 4.72,
			"ratings_count":  1420,
			"reviews_count":  184,
			"primary_genres": []string{"Tamil Film Music", "Orchestral Folk"},
			"descriptors":    []string{"melancholic", "cinematic", "poetic"},
		},
		{
			"rank":           2,
			"id":             "rel-roja",
			"title":          "Roja",
			"artist":         "A.R. Rahman",
			"artist_id":      "art-ar-rahman",
			"artist_image":   "assets/images/ar_rahman.jpg",
			"remote_artist":  "https://cdn-images.dzcdn.net/images/artist//500x500-000000-80-0-0.jpg",
			"release_year":   1992,
			"release_type":   "Film Soundtrack",
			"cover_art_url":  "assets/images/roja.jpg",
			"remote_cover":   "https://is1-ssl.mzstatic.com/image/thumb/Music221/v4/09/5d/7b/095d7bd0-001c-695a-4db1-eaa196ed0918/8905750032702.jpg/600x600bb.jpg",
			"average_rating": 4.75,
			"bayesian_score": 4.68,
			"ratings_count":  1310,
			"reviews_count":  152,
			"primary_genres": []string{"Synth-Pop", "Tamil Film Music"},
			"descriptors":    []string{"lush", "uplifting", "melodic"},
		},
		{
			"rank":           3,
			"id":             "rel-thalapathi",
			"title":          "Thalapathi",
			"artist":         "Ilaiyaraaja",
			"artist_id":      "art-ilaiyaraaja",
			"artist_image":   "assets/images/ilaiyaraaja.jpg",
			"remote_artist":  "https://cdn-images.dzcdn.net/images/artist/aeeca2a4b9f808c8e8f5d281b1fb48d0/500x500-000000-80-0-0.jpg",
			"release_year":   1991,
			"release_type":   "Film Soundtrack",
			"cover_art_url":  "assets/images/thalapathi.jpg",
			"remote_cover":   "https://is1-ssl.mzstatic.com/image/thumb/Music211/v4/72/c2/95/72c295f9-ccec-c74e-1ddb-083b649133dd/8905750030999.jpg/600x600bb.jpg",
			"average_rating": 4.72,
			"bayesian_score": 4.65,
			"ratings_count":  1180,
			"reviews_count":  135,
			"primary_genres": []string{"Carnatic Fusion", "Tamil Film Music"},
			"descriptors":    []string{"rhythmic", "epic", "dramatic"},
		},
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"chart_name": "Top Releases of All Time",
		"items":      charts,
	})
}

func (s *Server) handleGetLists(w http.ResponseWriter, r *http.Request) {
	setCORS(w)
	s.store.mu.RLock()
	lists := s.store.lists
	s.store.mu.RUnlock()
	json.NewEncoder(w).Encode(lists)
}

func (s *Server) handlePostRating(w http.ResponseWriter, r *http.Request) {
	setCORS(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	var req struct {
		EntityID string  `json:"entity_id"`
		Rating   float64 `json:"rating"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Rating < 0.5 || req.Rating > 5.0 {
		http.Error(w, "Rating must be between 0.5 and 5.0", http.StatusBadRequest)
		return
	}

	s.store.mu.Lock()
	s.store.ratings[req.EntityID] = append(s.store.ratings[req.EntityID], req.Rating)
	ratings := s.store.ratings[req.EntityID]
	s.store.mu.Unlock()

	avg, count := calculateStats(ratings)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":         "success",
		"entity_id":      req.EntityID,
		"new_rating":     req.Rating,
		"average_rating": math.Round(avg*100) / 100,
		"total_ratings":  count,
	})
}

func (s *Server) handleGetReviews(w http.ResponseWriter, r *http.Request) {
	setCORS(w)
	entityID := r.URL.Query().Get("entity_id")
	if entityID == "" {
		entityID = "rel-nayakan"
	}

	s.store.mu.RLock()
	reviews := s.store.reviews[entityID]
	s.store.mu.RUnlock()

	json.NewEncoder(w).Encode(reviews)
}

func (s *Server) handlePostReview(w http.ResponseWriter, r *http.Request) {
	setCORS(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	var req struct {
		EntityID  string `json:"entity_id"`
		Username  string `json:"username"`
		Content   string `json:"content"`
		IsSpoiler bool   `json:"is_spoiler"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Username == "" {
		req.Username = "MusicExplorer"
	}

	rev := catalog.Review{
		ID:         fmt.Sprintf("rev-%d", time.Now().UnixNano()),
		UserID:     "usr-local",
		Username:   req.Username,
		EntityType: "release",
		EntityID:   req.EntityID,
		Content:    req.Content,
		IsSpoiler:  req.IsSpoiler,
		LikesCount: 1,
		CreatedAt:  time.Now(),
	}

	s.store.mu.Lock()
	s.store.reviews[req.EntityID] = append([]catalog.Review{rev}, s.store.reviews[req.EntityID]...)
	s.store.mu.Unlock()

	json.NewEncoder(w).Encode(rev)
}

func (s *Server) handleSearch(w http.ResponseWriter, r *http.Request) {
	setCORS(w)
	q := r.URL.Query().Get("q")
	if q == "" {
		q = "Ilaiyaraaja"
	}

	result, err := s.mbClient.SearchArtist(r.Context(), q)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": %q}`, err.Error()), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(result)
}

func calculateStats(ratings []float64) (float64, int) {
	if len(ratings) == 0 {
		return 0.0, 0
	}
	sum := 0.0
	for _, r := range ratings {
		sum += r
	}
	return sum / float64(len(ratings)), len(ratings)
}

func calculateBayesian(avg float64, count int) float64 {
	m := 5.0  // weight parameter
	C := 3.50 // global average rating
	v := float64(count)
	return math.Round(((v/(v+m))*avg+((m/(v+m))*C))*100) / 100
}
