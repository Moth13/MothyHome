package web

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/moth13/home_client/internal/api"
	"github.com/moth13/home_client/internal/version"
)

// Action represents a UI action with icon support
type Action struct {
	Type    string
	Name    string
	Display string
	Icon    string
}

// WebHandler handles web UI requests
type WebHandler struct {
	server    *api.Server
	templates *template.Template
}

// NewWebHandler creates a new web handler
func NewWebHandler(server *api.Server) *WebHandler {
	h := &WebHandler{
		server: server,
	}

	// Parse templates from embedded filesystem
	tmpl := template.Must(template.ParseFS(templateFS, "templates/index.html"))
	h.templates = tmpl

	return h
}

// ServeStatic serves static files from embedded filesystem
func (h *WebHandler) ServeStatic(w http.ResponseWriter, r *http.Request) {
	// Get the path after /static/
	path := strings.TrimPrefix(r.URL.Path, "/static/")
	if path == r.URL.Path || path == "" {
		http.NotFound(w, r)
		return
	}
	
	// Serve the file from embedded FS
	// staticFS contains "static/styles.css" so we need to prepend "static/"
	filePath := "static/" + path
	content, err := staticFS.ReadFile(filePath)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	
	// Determine content type
	contentType := "text/plain"
	if strings.HasSuffix(path, ".css") {
		contentType = "text/css"
	} else if strings.HasSuffix(path, ".js") {
		contentType = "application/javascript"
	}
	
	w.Header().Set("Content-Type", contentType)
	w.Write(content)
}

// Home renders the main control interface
func (h *WebHandler) Home(w http.ResponseWriter, r *http.Request) {
	// Get actions from server
	actions := h.server.GetActions()

	// Convert to UI actions with icons
	apps := convertToUIActions(actions.Apps, "app")
	keys := convertToUIActions(actions.Keys, "key")

	data := struct {
		Title   string
		Version string
		Apps    []Action
		Keys    []Action
	}{
		Title:   "Home Client - Contrôle TV",
		Version: version.Version,
		Apps:    apps,
		Keys:    keys,
	}

	if err := h.templates.Execute(w, data); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// convertToUIActions adds icons to actions
func convertToUIActions(apiActions []api.Action, actionType string) []Action {
	uiActions := make([]Action, 0, len(apiActions))

	// App icons
	appIcons := map[string]string{
		"netflix":   "📺",
		"disney":    "🏰",
		"dazn":      "⚽",
		"youtube":   "📹",
		"tv":        "📺",
		"prime":     "📦",
		"amazon":    "📦",
		"primevideo": "📦",
	}

	// Key icons
	keyIcons := map[string]string{
		"power":        "🔌",
		"pause":        "⏸️",
		"play":         "▶️",
		"stop":         "⏹️",
		"rewind":       "⏪",
		"fastforward":  "⏩",
		"volumeup":     "🔊",
		"volumedown":   "🔇",
		"mute":         "🔇",
		"up":           "⬆️",
		"down":         "⬇️",
		"left":         "⬅️",
		"right":        "➡️",
		"ok":           "✅",
		"back":         "↩️",
		"home":         "🏠",
		"menu":         "☰",
		"guide":        "📺",
		"input":        "🔌",
		"hdmi":         "🔌",
		"num0":         "0️⃣",
		"num1":         "1️⃣",
		"num2":         "2️⃣",
		"num3":         "3️⃣",
		"num4":         "4️⃣",
		"num5":         "5️⃣",
		"num6":         "6️⃣",
		"num7":         "7️⃣",
		"num8":         "8️⃣",
		"num9":         "9️⃣",
	}

	for _, a := range apiActions {
		var icon string
		if actionType == "app" {
			icon = appIcons[a.Name]
		} else {
			icon = keyIcons[a.Name]
		}
		if icon == "" {
			icon = string(strings.ToUpper(a.Display)[0])
		}
		uiActions = append(uiActions, Action{
			Type:    a.Type,
			Name:    a.Name,
			Display: a.Display,
			Icon:    icon,
		})
	}

	return uiActions
}
