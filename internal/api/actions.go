package api

import (
	"encoding/json"
	"net/http"
	"strings"
)

type Action struct {
	Type    string `json:"type"`
	Name    string `json:"name"`
	Display string `json:"display"`
	URI     string `json:"uri,omitempty"`
	Code    string `json:"code,omitempty"`
}

type ActionsResponse struct {
	Apps    []Action `json:"apps"`
	Keys    []Action `json:"keys"`
}

// GetActions returns all available actions grouped by type
func (server *Server) GetActions() ActionsResponse {
	response := ActionsResponse{
		Apps: []Action{},
		Keys: []Action{},
	}

	// Add apps
	for name, uri := range server.AppURIs {
		response.Apps = append(response.Apps, Action{
			Type:    "app",
			Name:    name,
			Display: formatDisplayName(name),
			URI:     uri,
		})
	}

	// Add keys (remote control codes)
	for name, code := range server.IrrCodes {
		response.Keys = append(response.Keys, Action{
			Type:    "key",
			Name:    name,
			Display: formatDisplayName(name),
			Code:    code,
		})
	}

	return response
}

// getAvailableActions returns the list of all available actions (GET /api/actions)
func (server *Server) getAvailableActions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	actions := server.GetActions()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(actions); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// executeAction executes an action (POST /api/action/execute/{type}/{name})
func (server *Server) executeAction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract type and name from path: /api/action/execute/{type}/{name}
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 5 {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	actionType := pathParts[4]
	actionName := pathParts[5]

	switch actionType {
	case "app":
		// Reuse existing sendAppRequest logic
		q := r.URL.Query()
		q.Set("name", actionName)
		r.URL.RawQuery = q.Encode()
		server.sendAppRequest(w, r)
	case "key":
		// Reuse existing sendKeyRequest logic
		// The key endpoint expects /key/{id}, so redirect internally
		http.Redirect(w, r, "/key/"+actionName, http.StatusSeeOther)
	default:
		http.Error(w, "Unknown action type", http.StatusBadRequest)
	}
}

// formatDisplayName converts action names to human-readable format
func formatDisplayName(name string) string {
	// Replace underscores with spaces
	name = strings.ReplaceAll(name, "_", " ")
	// Capitalize first letter of each word
	parts := strings.Fields(name)
	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.ToUpper(string(part[0])) + strings.ToLower(part[1:])
		}
	}
	return strings.Join(parts, " ")
}
