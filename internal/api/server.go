package api

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/moth13/home_client/internal/utils"
)

type Server struct {
	TVIP     string
	TVPSK    string
	Client   *http.Client
	AppURIs  map[string]string
	IrrCodes map[string]string
}

func NewServer(tvIP string, tvPSK string, appURIs map[string]string) *Server {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	server := &Server{
		TVIP:    tvIP,
		TVPSK:   tvPSK,
		Client:  client,
		AppURIs: appURIs,
	}

	server.init()

	http.HandleFunc("/app", server.sendAppRequest)
	http.HandleFunc("/key/{id}", server.sendKeyRequest)

	return server
}

func (server *Server) Start(address string) error {
	log.Printf("Home client server started")
	log.Fatal(http.ListenAndServe(address, nil))

	return nil
}

type BundleInfo struct {
	Bundled bool   `json:"bundled"`
	Type    string `json:"type"`
}

type RemoteKey struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type IRResponse struct {
	Result []json.RawMessage `json:"result"`
	ID     int               `json:"id"`
}

func (server *Server) init() {

	reqBody := utils.BraviaRequest{
		Method:  "getRemoteControllerInfo",
		ID:      1,
		Params:  []any{},
		Version: "1.0",
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		log.Printf("JSON error: %v", err)
		return
	}

	targetURL := "http://" + server.TVIP + "/sony/system"
	req, err := http.NewRequest(http.MethodPost, targetURL, bytes.NewReader(bodyBytes))
	if err != nil {
		log.Printf("Query error: %v", err)
		return
	}

	resp, err := server.Client.Do(req)
	if err != nil {
		log.Printf("Query to TV error: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("TV answer %d", resp.StatusCode)
		return
	}

	var irresp IRResponse
	err = json.NewDecoder(resp.Body).Decode(&irresp)
	if err != nil {
		log.Printf("JSON decode error: %v", err)
		return
	}

	var keys []RemoteKey
	if len(irresp.Result) > 1 {
		if err := json.Unmarshal(irresp.Result[1], &keys); err != nil {
			log.Printf("JSON decode error: %v", err)
		}
	}

	irccCodes := make(map[string]string)
	for _, key := range keys {
		irccCodes[strings.ToLower(key.Name)] = key.Value
	}
	server.IrrCodes = irccCodes
}
