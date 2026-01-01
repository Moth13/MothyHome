package main

import (
	"log"
	"os"

	"github.com/moth13/home_client/internal/api"
)

func main() {
	tvIP := os.Getenv("TV_IP")
	tvPSK := os.Getenv("TV_PSK")

	appURIs := map[string]string{
		"netflix": os.Getenv("NETFLIX_URI"),
		"disney":  os.Getenv("DISNEY_URI"),
		"youtube": os.Getenv("YOUTUBE_URI"),
		"dazn":    os.Getenv("DAZN_URI"),
		"tv":      os.Getenv("TV_URI"),
	}

	if tvIP == "" || tvPSK == "" {
		log.Fatal("Missing TV_IP or TV_PSK")
	}

	server := api.NewServer(tvIP, tvPSK, appURIs)
	addr := ":8080"
	server.Start(addr)
}
