package api

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/moth13/home_client/internal/utils"
)

func (server *Server) sendAppRequest(w http.ResponseWriter, r *http.Request) {
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

	uri, ok := server.AppURIs[name]
	if !ok || uri == "" {
		http.Error(w, "Unkown app", http.StatusBadRequest)
		return
	}

	reqBody := utils.BraviaRequest{
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

	targetURL := "http://" + server.TVIP + "/sony/appControl"
	req, err := http.NewRequest(http.MethodPost, targetURL, bytes.NewReader(bodyBytes))
	if err != nil {
		log.Printf("Query error: %v", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Auth-PSK", server.TVPSK)

	resp, err := server.Client.Do(req)
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
}
