# Prism

A neon-styled macOS terminal multiplexer. Dark background, color-coded tabs, up to 6 terminal panes per tab.

---

## Installation

> **Important:** There is another unrelated macOS app called "Prism 11" in the default Homebrew cask registry.
> Always install using the full tap path and `--cask` flag as shown below to avoid installing the wrong app.

### Homebrew (recommended)

```bash
# Step 1 — tap the correct repository
brew tap nowgnas/tap

# Step 2 — install the cask (--cask flag is required)
brew install --cask nowgnas/tap/prism
```

Or in a single command:

```bash
brew install --cask nowgnas/tap/prism
```

### Update

```bash
brew upgrade --cask nowgnas/tap/prism
```

### Uninstall

```bash
brew uninstall --cask nowgnas/tap/prism
```

---

## First Launch — Security Warning

Because Prism is not notarized by Apple, macOS may show:
- **"Prism is damaged and can't be opened"**
- **"Apple could not verify Prism is free of malware"**

This is expected for open-source apps distributed outside the App Store.

**Fix:**

```bash
sudo xattr -cr /Applications/Prism.app
```

Then open Prism normally from `/Applications` or Spotlight.

Alternatively:
1. Open **System Settings → Privacy & Security**
2. Scroll to the bottom
3. Click **"Open Anyway"** next to the Prism warning

---

## Features

- **6 neon-color tabs** — each tab has its own accent color
- **Up to 6 panes per tab** — auto grid layout (1 → 2 → 3 → 4 → 5 → 6)
- **Status indicators** — ○ idle · ● running · ◈ AI · ✕ error
- **AI detection** — detects claude, aider, ollama, llm, etc.
- **Shift+Tab** — cycle between tabs

---

## Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `⌘T` | New tab |
| `⌘W` | Close tab |
| `⌘D` | Add pane |
| `⇧Tab` | Next tab (cycle forward) |
| `⌘⇧]` | Next tab |
| `⌘⇧[` | Previous tab |
| `⌘⌥→` | Focus next pane |
| `⌘⌥←` | Focus previous pane |
| Double-click tab name | Rename tab |
| Right-click tab | Context menu (rename, add pane, close) |

---

## Build from Source

```bash
# Install dependencies
brew install xcodegen

# Generate Xcode project
xcodegen generate

# Open in Xcode
open Prism.xcodeproj
```

---

## Release

```bash
git tag v0.1.0
git push origin v0.1.0
# GitHub Actions builds Prism.app → uploads to Releases → updates Homebrew cask automatically
```
