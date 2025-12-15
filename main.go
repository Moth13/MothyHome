package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type braviaRequest struct {
	Method  string        `json:"method"`
	ID      int           `json:"id"`
	Params  []interface{} `json:"params"`
	Version string        `json:"version"`
}

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

	// Vérification minimale des URIs
	for name, uri := range appURIs {
		if uri == "" {
			log.Printf("No uri app for %s", name)
		}
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	http.HandleFunc("/app", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Unauthorized method", http.StatusMethodNotAllowed)
			return
		}

		q, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			http.Error(w, "Invalid query", http.StatusBadRequest)
			return
		}

		name := strings.ToLower(strings.TrimSpace(q.Get("name")))
		if name == "" {
			http.Error(w, "Missing required parameter name (ex: ?name=netflix)", http.StatusBadRequest)
			return
		}

		uri, ok := appURIs[name]
		if !ok || uri == "" {
			http.Error(w, "Unkown app", http.StatusBadRequest)
			return
		}

		reqBody := braviaRequest{
			Method:  "setActiveApp",
			ID:      1,
			Params:  []interface{}{map[string]string{"uri": uri}},
			Version: "1.0",
		}

		bodyBytes, err := json.Marshal(reqBody)
		if err != nil {
			log.Printf("JSON error: %v", err)
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}

		targetURL := "http://" + tvIP + "/sony/appControl"
		req, err := http.NewRequest(http.MethodPost, targetURL, bytes.NewReader(bodyBytes))
		if err != nil {
			log.Printf("Query error: %v", err)
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Auth-PSK", tvPSK)

		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Query to TV error: %v", err)
			http.Error(w, "Can't reach the TV", http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Printf("TV answer %d", resp.StatusCode)
			http.Error(w, "TV error", http.StatusBadGateway)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})

	addr := ":8080"
	log.Printf("Home client server started")
	log.Fatal(http.ListenAndServe(addr, nil))
}
