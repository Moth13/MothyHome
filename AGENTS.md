# AGENTS.md - Home Client Project

## Overview
**home_client** is a Go application that provides HTTP API endpoints to control Sony Bravia TVs remotely. It now includes a **web UI** built with HTMX and Go's standard `html/template` for dynamic control.

**Current Version: 0.1.0**

## Project Structure
```
.
├── cmd/
│   └── main.go              # Entry point - starts HTTP server
├── internal/
│   ├── api/
│   │   ├── app.go           # App launching endpoint (/app) - GET/POST
│   │   ├── key.go           # Key/Remote control endpoint (/key/{id}) - GET/POST
│   │   ├── server.go        # Server configuration and initialization
│   │   └── actions.go       # UI API endpoints (/api/actions, /api/action/execute)
│   ├── utils/
│   │   └── const.go         # BraviaRequest struct
│   ├── version/
│   │   └── version.go       # Application version
│   └── web/
│       ├── handler.go        # Web UI handler for home page
│       ├── assets.go         # Embedded templates and static files
│       ├── templates/
│       │   └── index.html    # Main UI template with version display
│       └── static/
│           └── styles.css    # CSS styles for the UI
├── CHANGELOG.md             # Version history
├── Dockerfile
├── Readme.md
├── go.mod
└── .gitignore
```

## Core Functionality
- **TV Control**: Communicates with Sony Bravia TVs via their REST API
- **App Launching**: Opens specific applications (Netflix, YouTube, etc.) on the TV
- **Remote Control**: Sends IRCC codes to control the TV (power, volume, navigation, etc.)
- **Web UI**: Dynamic interface with tabs for Applications and Remote Keys
- **Version Display**: Shows current version in the UI header
- **Authentication**: Uses PSK (Pre-Shared Key) for TV API authentication

## API Endpoints

### REST API (for iOS Shortcuts - unchanged)
- **GET/POST `/app`** - Launch an application
  - Query: `?name={app_name}` (e.g., `?name=netflix`)
- **GET/POST `/key/{id}`** - Send a remote control key press
  - Path: `/key/pause`, `/key/volumeup`, etc.

### Web UI API (new)
- **GET `/api/actions`** - Get all available actions (apps and keys)
  - Returns: `{"apps": [...], "keys": [...]}`
- **POST `/api/action/execute/{type}/{name}`** - Execute an action
  - Example: `POST /api/action/execute/app/netflix`
  - Example: `POST /api/action/execute/key/pause`

### Web UI Pages
- **GET `/`** - Main control interface with tabs and version display

## Configuration
The server requires:
- `TV_IP`: TV's IP address on the local network (e.g., `192.168.1.198`)
- `TV_PSK`: Pre-Shared Key for TV API authentication
- `AppURIs`: Map of application names to their URIs (set via environment variables)

### Example `.env` or environment variables:
```env
TV_IP=192.168.1.198
TV_PSK=your_pre_shared_key
NETFLIX_URI=com.sony.dtv.com.netflix.ninja.com.netflix.ninja.MainActivity
DISNEY_URI=com.sony.dtv.com.disney.disneyplus.com.bamtechmedia.dominguez.main.MainActivity
YOUTUBE_URI=com.sony.dtv.com.google.android.youtube.tv.com.google.android.apps.youtube.tv.activity.ShellActivity
DAZN_URI=com.sony.dtv.com.dazn.com.dazn.MainActivity
TV_URI=com.sony.dtv.com.sony.dtv.tvx.com.sony.dtv.tvx.MainActivity
```

## Web UI Features

### Design
- **Dark theme** with Slate color palette
- **Responsive grid** - adapts to mobile, tablet, and desktop
- **Tab navigation** - separate Applications and Remote Keys
- **Icon support** - emoji icons for known apps and keys
- **Visual feedback** - buttons show loading state during requests
- **Version display** - shows current version in header

### User Experience
- Single page application (no full page reloads)
- HTMX for dynamic updates
- Instant feedback on button clicks
- Clean, minimalist interface
- Works on mobile and desktop

## Dependencies
- **Go version**: 1.25+
- **External packages**: None (only standard library)
- **Frontend**:
  - HTMX 1.9.10 (loaded from CDN)

## Conventions
- **Code Style**: Follows Go standard formatting (`gofmt`)
- **Error Handling**: Uses `http.Error` for client-facing errors, `log.Printf` for server logging
- **Naming**: Uses camelCase for variables, PascalCase for types
- **HTTP**: All handlers use `*Server` receiver for shared state
- **Versioning**: Semantic versioning (0.1.0)

## TV API Reference
- **Base URL**: `http://{TV_IP}/sony`
- **Authentication Header**: `X-Auth-PSK: {PRE_SHARED_KEY}`
- **Content-Type**: `application/json`

## Development Notes
- All TV communication goes through `server.Client` (http.Client)
- JSON payloads use `utils.BraviaRequest` struct for consistency
- Query parameters are parsed with `url.ParseQuery`
- Actions are loaded once at server startup (from TV via `getRemoteControllerInfo`)
- Version is defined in `internal/version/version.go`

## Security Considerations
- PSK should be kept secret (in `.env` files, which are gitignored)
- Web UI has **no authentication** (designed for local/LAN use only)
- Only GET and POST methods are accepted on API endpoints
- Input validation is performed on all query parameters

## Running the Project

### Development
```bash
# Build and run
go run cmd/main.go

# Or with hot reload
air
```

### Production
```bash
# Build
go build -o home_client cmd/main.go

# Run
./home_client
```

### Docker
```bash
# Build the Docker image
docker build -t home-client .

# Run with environment variables
docker run -d \
  -p 8080:8080 \
  -e TV_IP=192.168.1.198 \
  -e TV_PSK=your_psk \
  -e NETFLIX_URI=... \
  home-client
```

The web UI will be available at `http://localhost:8080`

## Adding New Features

### Recommandations for Contributors

1. **Feature Branches**: Use `feat/{feature-name}` for new features
   ```bash
   git checkout -b feat/web-ui-tabs
   ```

2. **Commit Messages**: Include `feat:` prefix for feature-related commits
   ```bash
   git commit -m "feat: add web UI with HTMX and Templ"
   ```

3. **Pull Requests**: Use `feat/{feature-name}` as PR title prefix
   - Example: `feat/web-ui: Add dynamic control interface`

4. **Version Management**:
   - Update version in `internal/version/version.go`
   - Update CHANGELOG.md with new version
   - Tag releases with `v{X.Y.Z}` format

5. **Testing**:
   - Test with real TV device when possible
   - Verify both API and UI functionality
   - Test on mobile devices

6. **Documentation**:
   - Update AGENTS.md with new features
   - Update Readme.md with usage instructions
   - Add to CHANGELOG.md

### Adding New Applications
1. Add the app URI to your environment variables:
   ```env
   NEW_APP_URI=com.sony.dtv.com.newapp.MainActivity
   ```
2. Add to `appURIs` map in `cmd/main.go`:
   ```go
   appURIs := map[string]string{
       "new_app": os.Getenv("NEW_APP_URI"),
   }
   ```
3. Add icon in `internal/web/handler.go`:
   ```go
   appIcons := map[string]string{
       "new_app": "🎬",
   }
   ```
4. Update version if this is a significant change

### Adding New Endpoints
- Add handler in `internal/api/`
- Register route in `internal/api/server.go`
- Document in AGENTS.md and Readme.md
- Add to CHANGELOG.md

## Version History
- **0.1.0** (2026-08-23): Web UI with HTMX, tabs, responsive design
- **0.0.1** (2025-12-15): Initial release with REST API

## Extending
- Add new device drivers under `internal/devices/` (implement `send(action, value)`)
- Add pairing/authorization flows for devices that require it
- Customize UI by modifying templates in `internal/web/templates/`
- Add CSS styles in `internal/web/static/styles.css`
- Add new API endpoints in `internal/api/`
