# MothyHome 🚀

[![Go](https://img.shields.io/badge/go-1.25%2B-blue)](https://golang.org) [![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

Simple Go service to expose quick actions to control a Sony Bravia TV. Now features a **web UI** for easy control from any browser!

## ✨ Features

### API Service
- Run a small HTTP service that sends commands to networked devices (Sony Bravia initially)
- Trigger actions from iOS Shortcuts using "Get Contents of URL"
- Minimal, extensible configuration for adding more devices

### Web UI (NEW!)
- **Responsive design** - works on mobile, tablet, and desktop
- **Tab navigation** - separate Applications and Remote Control sections
- **Dynamic buttons** - automatically generated from configured apps and TV remote keys
- **HTMX-powered** - instant feedback without page reloads
- **Dark theme** - with Slate color palette
- **Icon support** - emoji icons for known apps and keys

## ⚙️ Requirements
- Go 1.25+ (or latest stable)
- Local network access to your Sony Bravia
- (Optional) TV pairing/PSK depending on model

## ⚡ Install / Build

### Development
```bash
# Clone the repository
git clone https://github.com/Moth13/MothyHome.git
cd MothyHome

# Build and run
go run cmd/main.go

# Or with hot reload (if air is configured)
air
```

### Production
```bash
# Build the binary
go build -o mothyhome cmd/main.go

# Run
./mothyhome
```

### Docker
```bash
# Build the Docker image
docker build -t mothyhome .

# Run with environment variables
docker run -d \
  -p 8080:8080 \
  -e TV_IP=192.168.1.198 \
  -e TV_PSK=your_psk \
  -e NETFLIX_URI=... \
  mothyhome
```

## 🛠️ Configuration

### Environment Variables

Create an `.env` file or set environment variables:

```env
# Required
TV_IP=192.168.1.198
TV_PSK=your_pre_shared_key

# Application URIs (optional - configure which apps you want to control)
NETFLIX_URI=com.sony.dtv.com.netflix.ninja.com.netflix.ninja.MainActivity
DISNEY_URI=com.sony.dtv.com.disney.disneyplus.com.bamtechmedia.dominguez.main.MainActivity
YOUTUBE_URI=com.sony.dtv.com.google.android.youtube.tv.com.google.android.apps.youtube.tv.activity.ShellActivity
DAZN_URI=com.sony.dtv.com.dazn.com.dazn.MainActivity
TV_URI=com.sony.dtv.com.sony.dtv.tvx.com.sony.dtv.tvx.MainActivity
```

### Finding URIs
To find the URI for an app on your TV:
1. Open the app on your Sony Bravia
2. Check the TV's app list or use the Sony Bravia API documentation
3. Common URIs are listed above

## 📡 HTTP API

### REST API (for iOS Shortcuts - unchanged)
- **GET `/app`** - Launch an application
  - Query: `?name={app_name}` (e.g., `?name=netflix`)
  
  Example curl to open netflix:
  ```bash
  curl -X GET "http://<host>:8080/app?name=netflix"
  ```

- **GET `/key/{keyname}`** - Send a remote control key press
  - Path: `/key/pause`, `/key/volumeup`, etc.
  
  Example curl to set pause:
  ```bash
  curl -X GET "http://<host>:8080/key/pause"
  ```

### Web UI API (NEW!)
- **GET `/api/actions`** - Get all available actions (apps and keys)
  - Returns JSON: `{"apps": [...], "keys": [...]}`
  
  Example:
  ```bash
  curl http://localhost:8080/api/actions
  ```

- **POST `/api/action/execute/{type}/{name}`** - Execute an action
  - Example: `POST /api/action/execute/app/netflix`
  - Example: `POST /api/action/execute/key/pause`

### Web UI
- **GET `/`** - Main control interface with tabs
  - Open in browser: `http://localhost:8080`
  - Displays all configured apps and available remote keys
  - Responsive grid layout

## 🌐 Web Interface

### Access
After starting the server, open your browser to:
```
http://localhost:8080
```

### Features
- **Applications Tab**: Shows all configured streaming apps
- **Remote Keys Tab**: Shows all available TV remote control buttons
- **Responsive Grid**: Adapts to your screen size
- **Visual Feedback**: Buttons show loading state during requests
- **No Authentication**: Designed for local/LAN use only

### Screenshot (Concept)
```
┌─────────────────────────────────────┐
│     ⌂ Home Client                     │
│     Contrôle Sony Bravia              │
│                                     │
│  [Applications] [Touches]            │
├─────────────────────────────────────┤
│  📺 Netflix    🏰 Disney+    📹 YouTube │
│  ⚽ DAZN      📺 TV                   │
├─────────────────────────────────────┤
│  ⏸️ Pause    ▶️ Play    ⏹️ Stop      │
│  🔊 Volume+   🔇 Volume-               │
└─────────────────────────────────────┘
```

## 🧩 Extending

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
3. Restart the server

### Adding Custom Icons
Edit the icon mappings in `internal/web/handler.go`:
```go
appIcons := map[string]string{
    "new_app": "🎬",
}
keyIcons := map[string]string{
    "my_key": "🔘",
}
```

### Adding New Device Types
- Add device drivers under `internal/devices/` (implement `send(action, value)`)
- Add pairing/authorization flows for devices that require it

## 🎯 Project Structure
```
.
├── cmd/
│   └── main.go              # Entry point
├── internal/
│   ├── api/
│   │   ├── app.go           # App launching endpoint
│   │   ├── key.go           # Key/Remote control endpoint
│   │   ├── server.go        # Server configuration
│   │   └── actions.go       # Web UI API endpoints
│   ├── utils/
│   │   └── const.go         # BraviaRequest struct
│   └── web/
│       ├── handler.go        # Web UI handler
│       ├── assets.go         # Embedded files
│       ├── templates/
│       │   └── index.html    # Main UI template
│       └── static/
│           └── styles.css    # CSS styles
├── Dockerfile
├── Readme.md
└── go.mod
```

## 🤝 Contributing
PRs and issues welcome. Keep changes focused and documented.

## 📄 License
MIT — see LICENSE file.
