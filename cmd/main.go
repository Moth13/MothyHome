package main

import (
	"log"
	"net/http"
	"os"

	"github.com/moth13/home_client/internal/api"
	"github.com/moth13/home_client/internal/web"
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
	
	// Setup web handlers
	webHandler := web.NewWebHandler(server)
	
	// Web UI routes
	http.HandleFunc("/", webHandler.Home)
	http.HandleFunc("/static/", webHandler.ServeStatic)
	
	addr := ":8080"
	log.Printf("Home client server started on %s", addr)
	log.Printf("Web UI available at http://localhost%s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
