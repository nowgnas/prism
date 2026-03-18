# Prism

A neon-styled macOS terminal multiplexer. Dark background, color-coded tabs, up to 6 terminal panes per tab.

## Install

```bash
brew tap nowgnas/tap
brew install --cask prism
```

> First launch: if macOS blocks the app, run `sudo xattr -cr /Applications/Prism.app`

## Features

- **6 neon-color tabs** — each tab has its own accent color
- **Up to 6 panes per tab** — auto grid layout (1→2→3→4→5→6)
- **Status indicators** — ○ idle · ● running · ◈ AI · ✕ error
- **AI detection** — detects claude, aider, ollama, llm, etc.
- **Shift+Tab** — cycle between tabs

## Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `⌘T` | New tab |
| `⌘W` | Close tab |
| `⌘D` | Add pane |
| `⇧Tab` | Next tab |
| Double-click tab | Rename tab |

## Build from source

```bash
brew install xcodegen
xcodegen generate
open Prism.xcodeproj
```

## Release

```bash
git tag v0.1.0
git push origin v0.1.0
# GitHub Actions builds Prism.app → uploads to Releases → updates Homebrew cask
```
