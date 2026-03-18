import SwiftUI

@main
struct PrismApp: App {
    @StateObject private var appState = AppState()

    var body: some Scene {
        WindowGroup("Prism") {
            ContentView()
                .environmentObject(appState)
                .preferredColorScheme(.dark)
        }
        .windowStyle(.hiddenTitleBar)
        .windowToolbarStyle(.unified(showsTitle: false))
        .defaultSize(width: 1280, height: 800)
        .commands {
            CommandGroup(replacing: .newItem) {
                Button("New Tab") { appState.addTab() }
                    .keyboardShortcut("t", modifiers: .command)

                Button("Close Tab") { appState.closeActiveTab() }
                    .keyboardShortcut("w", modifiers: .command)

                Divider()

                Button("Add Pane") { appState.activeTab?.addSession() }
                    .keyboardShortcut("d", modifiers: .command)
            }

            CommandGroup(replacing: .help) {
                Button("Prism Help") {}
            }
        }
    }
}
