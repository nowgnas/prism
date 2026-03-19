import AppKit
import SwiftUI

@MainActor
final class AppState: ObservableObject {
    @Published var tabs: [TabModel] = []
    @Published var activeTabIndex: Int = 0

    private var nextPaletteIndex = 0
    private var keyMonitor: Any?

    init() {
        addTab()
        installKeyMonitor()
    }

    // MARK: - Tab management

    func addTab() {
        let tab = TabModel(paletteIndex: nextPaletteIndex)
        nextPaletteIndex += 1
        tabs.append(tab)
        activeTabIndex = tabs.count - 1
    }

    func closeActiveTab() {
        guard tabs.count > 1 else { return }
        tabs[activeTabIndex].cleanup()
        tabs.remove(at: activeTabIndex)
        activeTabIndex = min(activeTabIndex, tabs.count - 1)
    }

    func nextTab() {
        guard !tabs.isEmpty else { return }
        activeTabIndex = (activeTabIndex + 1) % tabs.count
    }

    func prevTab() {
        guard !tabs.isEmpty else { return }
        activeTabIndex = (activeTabIndex - 1 + tabs.count) % tabs.count
    }

    func selectTab(at index: Int) {
        guard index < tabs.count else { return }
        activeTabIndex = index
    }

    var activeTab: TabModel? {
        guard activeTabIndex < tabs.count else { return nil }
        return tabs[activeTabIndex]
    }

    // MARK: - Keyboard

    private func installKeyMonitor() {
        keyMonitor = NSEvent.addLocalMonitorForEvents(matching: .keyDown) { [weak self] event in
            guard let self else { return event }

            let mods = event.modifierFlags.intersection(.deviceIndependentFlagsMask)

            // Shift+Tab → next tab (cycle forward)
            if event.specialKey == .tab && mods == .shift {
                Task { @MainActor in self.nextTab() }
                return nil
            }

            // ⌘⇧] → next tab, ⌘⇧[ → prev tab (standard macOS terminal convention)
            if mods == [.command, .shift] {
                if event.characters == "]" { Task { @MainActor in self.nextTab() }; return nil }
                if event.characters == "[" { Task { @MainActor in self.prevTab() }; return nil }
            }

            // ⌘⌥→ → next pane, ⌘⌥← → prev pane, ⌘⌥W → close pane within active tab
            if mods == [.command, .option] {
                if event.specialKey == .rightArrow {
                    Task { @MainActor in self.activeTab?.nextPane() }
                    return nil
                }
                if event.specialKey == .leftArrow {
                    Task { @MainActor in self.activeTab?.prevPane() }
                    return nil
                }
                if event.charactersIgnoringModifiers?.lowercased() == "w" {
                    Task { @MainActor in self.activeTab?.closeActivePane() }
                    return nil
                }
            }

            return event
        }
    }

    deinit {
        if let monitor = keyMonitor {
            NSEvent.removeMonitor(monitor)
        }
    }
}
