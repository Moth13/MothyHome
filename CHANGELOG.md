# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.0] - 2026-08-23

### Added
- **Web UI**: New dynamic interface built with HTMX for controlling Sony Bravia TV from any browser
- **Tabs Navigation**: Separate tabs for Applications and Remote Control keys
- **Responsive Design**: Grid layout that adapts to mobile, tablet, and desktop screens
- **API Endpoints**: New `/api/actions` and `/api/action/execute/{type}/{name}` endpoints for the web interface
- **Version Display**: Version number now displayed in the UI header
- **Dark Theme**: Modern dark theme with Slate color palette
- **Icon Support**: Emoji icons for known applications and remote keys

### Changed
- **API Methods**: `sendAppRequest` and `sendKeyRequest` now accept both GET and POST methods for compatibility with both iOS Shortcuts and web UI
- **Error Handling**: Improved HTTP status codes and error messages
- **Documentation**: Updated AGENTS.md, Readme.md with new UI information

### Fixed
- **Method Compatibility**: Fixed 405 errors by allowing POST requests on app and key endpoints

## [0.0.1] - 2025-12-15

### Added
- Initial release with REST API for Sony Bravia TV control
- Support for launching applications via `/app` endpoint
- Support for remote control keys via `/key/{id}` endpoint
- Configuration via environment variables (TV_IP, TV_PSK, app URIs)
- iOS Shortcuts integration

