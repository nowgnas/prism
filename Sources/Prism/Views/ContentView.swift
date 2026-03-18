import SwiftUI

struct ContentView: View {
    @EnvironmentObject var appState: AppState

    var body: some View {
        HSplitView {
            SidebarView()
                .frame(minWidth: 160, idealWidth: 180, maxWidth: 220)
                .background(NeonTheme.sidebarBg)

            Group {
                if let tab = appState.activeTab {
                    TabGridView(tab: tab)
                } else {
                    emptyState
                }
            }
            .frame(minWidth: 400, minHeight: 300)
            .background(NeonTheme.background)
        }
        .background(NeonTheme.background)
        .frame(minWidth: 700, minHeight: 450)
    }

    private var emptyState: some View {
        VStack(spacing: 12) {
            Text("✦")
                .font(.system(size: 32))
                .foregroundColor(NeonTheme.textDim)
            Text("No tabs open")
                .font(.system(.body, design: .monospaced))
                .foregroundColor(NeonTheme.textDim)
            Button("New Tab") { appState.addTab() }
                .keyboardShortcut("t", modifiers: .command)
        }
        .frame(maxWidth: .infinity, maxHeight: .infinity)
        .background(NeonTheme.background)
    }
}
