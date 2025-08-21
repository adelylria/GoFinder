# GoFinder

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**GoFinder** — a fast, lightweight application launcher written in Go using the **Fyne** GUI toolkit. It focuses on quick keyboard-driven discovery and launch of installed applications, system tray presence, and a global hotkey to toggle visibility.

---

# Status & Supported Platforms

* ## **Supported platforms:**

    - Windows is **partially** supported (extraction of icons from executables and shortcuts, global hotkey, tray integration). 

    - Linux implements discovery via `.desktop` files and returns `Application` objects but lacks some integrations (tray behavior and advanced icon extraction). 

    - macOS is scaffolded and ready for implementation.

---

# Repository layout (high level)

```
cmd/                # main entry & build resources
core/               # cross-cutting core modules (hotkey, resource, ui, logger)
logic/              # platform-specific discovery, execution and icon logic
models/             # data models (Application, AppState, ...)
```

---

# What GoFinder does

* Discover installed applications per-platform (`.lnk` on Windows, `.desktop` on Linux).
* Present a small keyboard-focused UI to search and launch apps.
* Show application icons when available (no reserved space if icon is missing).
* Register a global hotkey to toggle the launcher (native implementation for Windows via `cgo` bridge).
* Keep a system tray icon while the launcher runs in the background.

---

# Key design & architecture

* **Pluggable AppFinders**: `logic.AppFinder` is an interface implemented per OS (`windowsAppFinder`, `linuxAppFinder`, `darwinAppFinder`). `FindApplications()` picks the correct implementation at runtime.
* **Separation of concerns**: `core` holds UI and platform-agnostic code, `logic` contains OS-specific implementations (discovery, icon extraction, runner).
* **Icon pipeline (Windows)**: try image files (`.png`, `.ico`), `ExtractIconEx` with index, `SHGetFileInfo`, or as a last resort the executable icon. `HICON` objects are converted to Go `image.Image`, encoded as PNG and wrapped in `fyne.Resource`.
* **Caching**: icons are cached in-memory (`map[string]fyne.Resource`) protected by `sync.RWMutex` to avoid repeated extraction.
* **Hotkey**: a native (C) bridge registers a global hotkey on Windows; the Go side receives toggle/exit events.
* **UI**: `core/ui` exposes a `Launcher` and a `ThemeConfig` to centralize visual metrics and behavior (search entry, styled list, selection handling).

---

# Building

Requirements:

* Go 1.20+ (or compatible)
* A C toolchain for building code that uses `cgo` (Windows native hotkey registration) when compiling on or for Windows.

Build is handled via a **Makefile**. Common commands:

```bash
make run             # Run GoFinder locally
make build           # Build for the current platform (binary in ./build)
make build-windows   # Cross-compile for Windows
make build-linux-darwin # Cross-compile for Linux and macOS
make clean           # Clean build artifacts
```

---

# Dependencies

Add or verify the following in `go.mod`:

* `fyne.io/fyne/v2` — GUI toolkit.
* `github.com/lxn/win` — Win32 wrappers used for icon extraction and GDI.
* `golang.org/x/sys/windows` — low-level Windows syscalls.
* `github.com/parsiya/golnk` — parse Windows `.lnk` shortcuts.
* `github.com/fyne-io/image/ico` — decode `.ico` files.

---

# Usage

1. Build with `make build` (or the appropriate target).
2. Run the executable from `./build` (e.g., `./build/goFinder.exe` on Windows).
3. The app will place an icon in the system tray and remain running.
4. Use the configured global hotkey to toggle the launcher UI. Type to filter results, navigate with the arrow keys, and press Enter to launch.

---

# Troubleshooting & notes

* `.lnk` parsing may yield paths containing environment variables — the code expands and normalizes them and warns if targets are missing or inaccessible.
* Icon extraction uses Win32 APIs and may fail on restricted accounts; the code falls back gracefully to other icon sources.
* If the hotkey registration fails, verify that another application does not already hold the same combination and that the process has permission to register global hotkeys on the platform.

---

# Development notes (quick)

* Entry point: `cmd/main.go` → `logic.FindApplications()` → `ui.RunLauncher(apps)`.
* Platform-specific files use build tags (e.g. `//go:build windows`) — they will be excluded on unsupported platforms.
* To inspect icon extraction behavior, review `logic/windows_icon.go` and the conversion function `IconHandleToImage`.

--- 

# Roadmap

* \[\~] Linux — `.desktop` discovery implemented; add tray and enhanced icon resolution.
* \[\~] Darwin — Investigate
* \[\~] Tests
* \[\~] Logging
* \[\~] Themes 
* \[\~] Configuration for keybinding/themes ...

---

## 📄 License

GoFinder is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
