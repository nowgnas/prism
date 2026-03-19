import SwiftUI

struct ContentView: View {
    @EnvironmentObject var appState: AppState

    var body: some View {
        HSplitView {
            SidebarView()
                .frame(minWidth: 160, idealWidth: 180, maxWidth: 220)
                .background(NeonTheme.sidebarBg)

            TabHostView()
                .frame(minWidth: 400, minHeight: 300)
                .background(NeonTheme.background)
        }
        .background(NeonTheme.background)
        .frame(minWidth: 700, minHeight: 450)
    }
}
