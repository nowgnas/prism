import SwiftUI
import AppKit

// MARK: - SwiftUI bridge

/// Hosts all tab grid views simultaneously at the AppKit level.
/// Tab switching = isHidden toggle (O(1), no view recreation).
struct TabHostView: NSViewRepresentable {
    @EnvironmentObject var appState: AppState

    func makeNSView(context: Context) -> TabHostNSView {
        TabHostNSView()
    }

    func updateNSView(_ host: TabHostNSView, context: Context) {
        host.update(tabs: appState.tabs, activeIndex: appState.activeTabIndex)
    }
}

// MARK: - AppKit host

final class TabHostNSView: NSView {
    /// One NSHostingView per tab, keyed by tab ID. Never recreated.
    private var tabViews: [UUID: NSHostingView<TabGridView>] = [:]

    func update(tabs: [TabModel], activeIndex: Int) {
        let activeID = activeIndex < tabs.count ? tabs[activeIndex].id : nil

        // 1. Create hosting views for newly added tabs
        for tab in tabs where tabViews[tab.id] == nil {
            let hv = NSHostingView(rootView: TabGridView(tab: tab))
            hv.translatesAutoresizingMaskIntoConstraints = false
            addSubview(hv)
            NSLayoutConstraint.activate([
                hv.leadingAnchor.constraint(equalTo: leadingAnchor),
                hv.trailingAnchor.constraint(equalTo: trailingAnchor),
                hv.topAnchor.constraint(equalTo: topAnchor),
                hv.bottomAnchor.constraint(equalTo: bottomAnchor),
            ])
            tabViews[tab.id] = hv
        }

        // 2. Remove hosting views for deleted tabs
        let liveIDs = Set(tabs.map(\.id))
        for (id, view) in tabViews where !liveIDs.contains(id) {
            view.removeFromSuperview()
            tabViews.removeValue(forKey: id)
        }

        // 3. Show only the active tab — everything else is hidden
        //    isHidden = true stops Metal rendering for that subtree (GPU-free)
        for (id, view) in tabViews {
            view.isHidden = id != activeID
        }
    }
}
