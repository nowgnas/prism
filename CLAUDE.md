# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Run

**Xcode is required.** If `xcodebuild` fails with *"requires Xcode, but active developer directory is Command Line Tools"*, point `xcode-select` at the app bundle:

```bash
sudo xcode-select -s /Applications/Xcode.app/Contents/Developer
xcodebuild -version   # should print Xcode / Build version, not an error
```

Optional one-shot checks + Debug build:

```bash
./scripts/check-xcode-env.sh   # verify toolchain
./scripts/build-debug.sh       # xcodegen + xcodebuild Debug
```

```bash
# Generate Xcode project (required before first build or after project.yml changes)
xcodegen generate

# Build (debug)
xcodebuild -project Prism.xcodeproj -scheme Prism -configuration Debug build

# Build release archive
xcodebuild archive \
  -project Prism.xcodeproj \
  -scheme Prism \
  -configuration Release \
  -archivePath build/Prism.xcarchive \
  CODE_SIGNING_REQUIRED=NO CODE_SIGNING_ALLOWED=NO

# Run from Xcode
open Prism.xcodeproj   # then ⌘R
```

**Note:** There is no Swift Package Manager build for the app — it must be built through Xcode/xcodebuild. `Package.swift` is present only for tooling compatibility.

## Release / Distribution

```bash
# Tag triggers GitHub Actions → builds .app → zips → creates GitHub Release → updates homebrew-tap cask
git tag v0.X.Y
git push origin v0.X.Y
```

Homebrew cask install:
```bash
brew install --cask nowgnas/tap/prism   # --cask required; avoids unrelated "Prism 11" app
```

CI workflows: `.github/workflows/release.yml` (tag push → build + deploy), `.github/workflows/ci.yml` (push to main → debug build).

## Architecture

### State hierarchy
```
AppState (ObservableObject, @MainActor)
  └─ [TabModel] (ObservableObject, Identifiable)
       └─ [SessionModel] (ObservableObject, Identifiable)
            ├─ NeonPane?          ← strong ref to AppKit view (persists across tab switches)
            └─ TerminalCoordinator?
```

### View hierarchy
```
ContentView (HSplitView)
  ├─ SidebarView       — tab list, status icons, rename, context menu
  └─ TabHostView       — NSViewRepresentable → TabHostNSView
       └─ [NSHostingView<TabGridView>]  — one per tab, isHidden toggled
            └─ TabGridView             — GeometryReader + ZStack with offset layout
                 └─ [TerminalPaneView] — NSViewRepresentable → TerminalContainer
                      └─ NeonPane      — CALayer border/glow + LocalProcessTerminalView
```

### Key design decisions

**Terminal state preservation across tab switches:**
`SessionModel.neonPane` holds a strong reference to the `NeonPane` (AppKit view + PTY). `TerminalContainer` (thin `NSView`) is recreated by SwiftUI but only calls `attach(pane:)` to re-parent the already-running `NeonPane`. AppKit silently removes the pane from its old superview before adding to the new one — no PTY restart, no state loss.

**O(1) tab switching:**
`TabHostNSView` allocates one `NSHostingView<TabGridView>` per tab on creation and never destroys them. Switching tabs = toggling `isHidden`. Metal rendering pauses automatically for hidden views.

**SwiftUI/SwiftTerm Color conflict:**
Both frameworks export `Color`. Files importing both must qualify as `SwiftUI.Color` for SwiftUI colors and use `NSColor` for AppKit/SwiftTerm contexts.

**Neon border updates:**
`NeonPane.update(color:isActive:)` wraps `CALayer` mutations in `CATransaction.setDisableActions(true)` to prevent implicit animation lag on focus change.

**Rename gesture ordering:**
`SidebarView` declares `.onTapGesture(count: 2)` **before** `.onTapGesture` — SwiftUI processes gestures in declaration order, so the double-tap must be registered first.

### Status detection
`SessionModel` polls every 3 seconds via `AIDetector.childNames(of: shellPID)` (parses `ps -ax -o ppid=,comm=` output). Priority: `error > ai > running > idle`. `TabModel.status` is the `max()` across its sessions.

### Neon color palette (6 tabs, cycles)
Defined in `NeonTheme.swift`: Crimson `#E05C5C`, Amber `#E0A05C`, Jade `#5CC4A0`, Cobalt `#5C8CE0`, Violet `#A05CE0`, Rose `#E05CA0`.

### Key files
| File | Purpose |
|------|---------|
| `AppState.swift` | Root state; tab management; `NSEvent` key monitor (Shift+Tab, ⌘⇧[/], ⌘⌥←→) |
| `Models/TabModel.swift` | Tab state; pane add/remove/navigate; status aggregation |
| `Models/SessionModel.swift` | Per-pane state; owns `NeonPane`; polling timer |
| `Views/TerminalPaneView.swift` | `NSViewRepresentable` bridge; `demandPane()` creates PTY once; `NeonPane` CALayer border |
| `Views/TabHostView.swift` | `TabHostNSView`; isHidden-based tab switching |
| `Views/TabGridView.swift` | Grid layout algorithm for N=1–6 panes |
| `Views/SidebarView.swift` | Tab list UI; rename; context menu |
| `Theme/NeonTheme.swift` | Color palette; `Color(hex:)` / `NSColor(hex:)` extensions |
| `Utils/AIDetector.swift` | `ps`-based child process detection |
| `project.yml` | xcodegen spec (macOS 13+, no code signing, x86_64+arm64) |

## Known Issues / Gotchas

- **Gatekeeper warning on install:** Run `sudo xattr -cr /Applications/Prism.app` to clear quarantine attributes.
- **`isFlipped` must NOT be overridden in NeonPane:** SwiftTerm uses macOS bottom-left coordinates; overriding breaks ANSI rendering.
- **Max 6 panes per tab** is enforced in `TabModel.addSession()`.
- **Font fallback:** Tries MesloLGS NF → JetBrains Mono → system monospaced.
