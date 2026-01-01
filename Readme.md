# MothyHome 🚀

[![Go](https://img.shields.io/badge/go-1.25%2B-blue)](https://golang.org) [![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

Simple Go service to expose quick actions (for iOS Shortcuts) to control a Sony Bravia TV and other devices in the future.

## ✨ Features
- Run a small HTTP/CLI service that sends commands to networked devices (Sony Bravia initially).
- Trigger actions from iOS Shortcuts using "Get Contents of URL".
- Minimal, extensible configuration for adding more devices.

## ⚙️ Requirements
- Go 1.18+ (or latest stable)
- Local network access to your Sony Bravia
- (Optional) TV pairing/PSK depending on model

## ⚡ Install / Build
Clone and build:
```
git clone https://github.com/Moth13/MothyHome.git
cd MothyHome
go build -o mothyhome ./...
```

Run (default serves HTTP on :8080):
```
./mothyhome
```

Or run directly:
```
make server
```

## 🛠️ Configuration
Example `app.env`:
```env
NETFLIX_URI=com.sony.dtv.com.netflix.ninja.com.netflix.ninja.MainActivity
DISNEY_URI=com.sony.dtv.com.disney.disneyplus.com.bamtechmedia.dominguez.main.MainActivity
YOUTUBE_URI=com.sony.dtv.com.google.android.youtube.tv.com.google.android.apps.youtube.tv.activity.ShellActivity
DAZN_URI=com.sony.dtv.com.dazn.com.dazn.MainActivity
TV_URI=com.sony.dtv.com.sony.dtv.tvx.com.sony.dtv.tvx.MainActivity
TV_IP=192.168.1.198
TV_PSK=#TO GET IN YOUR TV
```

## 📡 HTTP API (for iOS Shortcuts)
- POST /app/{appname}

Example curl to open netflix:
```
curl -X POST http://<host>:8080/app/netflix
```
- POST /key/{keyname}

Example curl to set pause:
```
curl -X POST http://<host>:8080/key/pause
```

Use this URL in an iOS Shortcut with "Get Contents of URL" (POST, JSON body).

## 🧩 Extending
- Add device drivers under internal/devices (implement send(action, value)).
- Add pairing/authorization flows for devices that require it.

## 🤝 Contributing
PRs and issues welcome. Keep changes focused and documented.

## 📄 License
MIT — see LICENSE file.