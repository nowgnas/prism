import SwiftUI

struct SidebarView: View {
    @EnvironmentObject var appState: AppState
    @State private var renamingTabID: UUID? = nil
    @State private var renameText: String = ""

    var body: some View {
        VStack(spacing: 0) {
            header
            Divider().background(NeonTheme.divider)
            tabList
            Spacer(minLength: 0)
            Divider().background(NeonTheme.divider)
            bottomBar
        }
        .background(NeonTheme.sidebarBg)
    }

    // MARK: Header

    private var header: some View {
        HStack {
            Text("PRISM")
                .font(.system(size: 10, weight: .bold, design: .monospaced))
                .foregroundColor(NeonTheme.textMid)
                .tracking(3)
            Spacer()
        }
        .padding(.horizontal, 14)
        .padding(.vertical, 10)
    }

    // MARK: Tab list

    private var tabList: some View {
        ScrollView(.vertical, showsIndicators: false) {
            LazyVStack(spacing: 1) {
                ForEach(Array(appState.tabs.enumerated()), id: \.element.id) { idx, tab in
                    tabRow(tab: tab, index: idx)
                }
            }
            .padding(.vertical, 6)
        }
    }

    @ViewBuilder
    private func tabRow(tab: TabModel, index: Int) -> some View {
        let isActive = index == appState.activeTabIndex

        HStack(spacing: 0) {
            // Active bar
            Rectangle()
                .fill(isActive ? tab.color : Color.clear)
                .frame(width: 2)

            HStack(spacing: 7) {
                // Status icon
                Text(tab.status.icon)
                    .font(.system(size: 9, design: .monospaced))
                    .foregroundColor(tab.status.iconColor(tabColor: tab.color))
                    .frame(width: 12)

                // Name (inline rename on double-click)
                if renamingTabID == tab.id {
                    TextField("", text: $renameText)
                        .textFieldStyle(.plain)
                        .font(.system(size: 12, design: .monospaced))
                        .foregroundColor(NeonTheme.textPrimary)
                        .onSubmit { commitRename(tab: tab) }
                        .onExitCommand { renamingTabID = nil }
                } else {
                    Text(tab.name)
                        .font(.system(size: 12, design: .monospaced))
                        .foregroundColor(isActive ? NeonTheme.textPrimary : NeonTheme.textMid)
                        .lineLimit(1)
                        .truncationMode(.tail)
                }

                Spacer(minLength: 0)

                // Pane count badge
                if tab.sessions.count > 1 {
                    Text("\(tab.sessions.count)")
                        .font(.system(size: 9, design: .monospaced))
                        .foregroundColor(NeonTheme.textDim)
                        .padding(.horizontal, 4)
                        .padding(.vertical, 1)
                        .background(
                            RoundedRectangle(cornerRadius: 3)
                                .fill(NeonTheme.divider)
                        )
                }
            }
            .padding(.horizontal, 10)
            .padding(.vertical, 6)
        }
        .background(
            RoundedRectangle(cornerRadius: 0)
                .fill(isActive ? tab.color.opacity(0.08) : Color.clear)
        )
        .contentShape(Rectangle())
        .onTapGesture(count: 2) { beginRename(tab: tab) }
        .onTapGesture { appState.activeTabIndex = index }
        .contextMenu { tabContextMenu(tab: tab, index: index) }
    }

    // MARK: Context menu

    @ViewBuilder
    private func tabContextMenu(tab: TabModel, index: Int) -> some View {
        Button("Rename") { beginRename(tab: tab) }
        Button("Add Pane") { tab.addSession() }
        Divider()
        Button("Close Tab", role: .destructive) {
            appState.tabs[index].cleanup()
            if appState.tabs.count > 1 {
                appState.tabs.remove(at: index)
                appState.activeTabIndex = min(appState.activeTabIndex, appState.tabs.count - 1)
            }
        }
    }

    // MARK: Bottom bar

    private var bottomBar: some View {
        HStack(spacing: 4) {
            Button(action: { appState.addTab() }) {
                Image(systemName: "plus")
                    .font(.system(size: 11))
                    .foregroundColor(NeonTheme.textMid)
            }
            .buttonStyle(.plain)
            .help("New Tab  ⌘T")

            Spacer()

            // Tab index hint
            if appState.tabs.count > 1 {
                Text("⇧Tab")
                    .font(.system(size: 9, design: .monospaced))
                    .foregroundColor(NeonTheme.textDim)
            }
        }
        .padding(.horizontal, 12)
        .padding(.vertical, 7)
    }

    // MARK: Rename helpers

    private func beginRename(tab: TabModel) {
        renameText = tab.name
        renamingTabID = tab.id
    }

    private func commitRename(tab: TabModel) {
        let trimmed = renameText.trimmingCharacters(in: .whitespaces)
        if !trimmed.isEmpty { tab.name = trimmed }
        renamingTabID = nil
    }
}
