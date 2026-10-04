package main

import (
	"context"
	"fmt"
	"log"

	"devops/musicdb/pkg/catalog"
	"devops/musicdb/pkg/musicbrainz"
)

func main() {
	log.Println("[+] Starting Indian Music Catalog & Rating Server (Go Backend)")

	mbClient := musicbrainz.NewClient(musicbrainz.Config{
		AppName:     "IndianMusicDB",
		Version:     "0.1.0",
		ContactInfo: "admin@example.com",
	})

	log.Println("[+] MusicBrainz rate-limited client initialized (1 req/sec enforced).")

	ctx := context.Background()

	// Query 1: Search for Ilaiyaraaja
	log.Println("[+] Executing rate-limited query 1: Search artist 'Ilaiyaraaja'...")
	res1, err := mbClient.SearchArtist(ctx, "Ilaiyaraaja")
	if err != nil {
		log.Printf("[-] Error searching artist: %v\n", err)
	} else if len(res1.Artists) > 0 {
		artist := res1.Artists[0]
		fmt.Printf("    -> Found Artist: %s | MBID: %s | Type: %s | Country: %s\n",
			artist.Name, artist.ID, artist.Type, artist.Country)
	}

	// Query 2: Search for Tyagaraja (Demonstrating rate limiter delay)
	log.Println("[+] Executing rate-limited query 2: Search artist 'Tyagaraja'...")
	res2, err := mbClient.SearchArtist(ctx, "Tyagaraja")
	if err != nil {
		log.Printf("[-] Error searching artist: %v\n", err)
	} else if len(res2.Artists) > 0 {
		artist := res2.Artists[0]
		fmt.Printf("    -> Found Artist: %s | MBID: %s | Type: %s | Country: %s\n",
			artist.Name, artist.ID, artist.Type, artist.Country)
	}

	log.Println("[+] Local test run complete!")
	_ = catalog.Artist{}
}
