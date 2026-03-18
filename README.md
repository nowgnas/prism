# prism

A colorful TUI terminal multiplexer — one terminal session split into multiple color-coded tabs, each with up to 6 real PTY panes.

```
┌─────────────────────────────────────────────────────────────────┐
│ PRISM              │  ┌───────────────┐  ┌───────────────┐      │
│ ────────────────── │  │ ~/project     │  │ ~/logs        │      │
│ 1 ● project        │  │               │  │               │      │
│ 2 ○ logs           │  │  $ npm run dev│  │  tail -f ...  │      │
│ 3 ◈ ai-session     │  └───────────────┘  └───────────────┘      │
│                    │  ┌───────────────────────────────────┐      │
│                    │  │ ~/project                         │      │
│                    │  │  $ git log --graph                │      │
│                    │  └───────────────────────────────────┘      │
└─────────────────────────────────────────────────────────────────┘
```

## Install

```bash
# Homebrew (macOS)
brew install nowgnas/tap/prism

# go install
go install github.com/nowgnas/prism@latest

# Direct download — see GitHub Releases
```

## Usage

```bash
prism                          # launch with defaults
prism --config ~/.config/prism/config.toml
prism --version
```

## Key Bindings

All commands are prefixed with **Ctrl+B**.

| Key         | Action                        |
|-------------|-------------------------------|
| `^B n`      | New tab                       |
| `^B w`      | Close current tab             |
| `^B ,`      | Rename tab                    |
| `^B Tab`    | Next tab                      |
| `^B 1-9`    | Jump to tab N                 |
| `^B "`      | Add pane (horizontal split)   |
| `^B %`      | Add pane (vertical split)     |
| `^B x`      | Close current pane            |
| `^B z`      | Zoom / unzoom pane            |
| `^B h/j/k/l`| Move pane focus               |
| `^B ?`      | Toggle help overlay           |
| `^B q`      | Quit                          |

## Tab status icons

| Icon | Meaning         |
|------|-----------------|
| `○`  | Idle            |
| `●`  | Running command |
| `◈`  | AI tool active  |
| `✕`  | Error / exited  |

## Configuration

Config is stored at `~/.config/prism/config.toml` (created on first run):

```toml
shell         = "/bin/zsh"
startup_tabs  = 1
sidebar_width = 20
poll_interval_ms = 500
ai_processes  = ["claude", "aider", "llm", "ollama"]
theme         = "dark"
mouse_enabled = false
```

## Tab color palette

| # | Name    | Hex       |
|---|---------|-----------|
| 0 | Crimson | `#E05C5C` |
| 1 | Amber   | `#E0A05C` |
| 2 | Jade    | `#5CC4A0` |
| 3 | Cobalt  | `#5C8CE0` |
| 4 | Violet  | `#A05CE0` |
| 5 | Rose    | `#E05CA0` |

## Building from source

```bash
go build ./cmd/prism
```

## Release (maintainers)

```bash
git tag v0.1.0
git push origin v0.1.0
# GitHub Actions runs goreleaser automatically
```
