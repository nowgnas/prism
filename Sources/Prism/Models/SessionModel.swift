import SwiftUI
import Combine
import AppKit

// MARK: - Status

enum TerminalStatus: Int, Comparable {
    case idle    = 0  // ○ dim
    case running = 1  // ● tab color
    case ai      = 2  // ◈ cyan
    case error   = 3  // ✕ red

    static func < (lhs: TerminalStatus, rhs: TerminalStatus) -> Bool {
        lhs.rawValue < rhs.rawValue
    }

    var icon: String {
        switch self {
        case .idle:    return "○"
        case .running: return "●"
        case .ai:      return "◈"
        case .error:   return "✕"
        }
    }

    func iconColor(tabColor: Color) -> Color {
        switch self {
        case .idle:    return NeonTheme.textDim
        case .running: return tabColor
        case .ai:      return Color(hex: "00FFFF")
        case .error:   return Color(hex: "FF3B3B")
        }
    }
}

// MARK: - SessionModel

/// Represents a single PTY terminal session (one pane).
/// Holds strong references to its AppKit views so they survive tab switches.
final class SessionModel: ObservableObject, Identifiable {
    let id = UUID()
    @Published var status: TerminalStatus = .idle
    @Published var title: String = ""

    /// Cached AppKit container — created once, never recreated on tab switch.
    var neonPane: NeonPane?
    /// Cached coordinator — reused across NSViewRepresentable lifecycles.
    var terminalCoordinator: TerminalCoordinator?

    private var pollTimer: Timer?

    init() {}

    func startPolling() {
        stopPolling()
        pollTimer = Timer.scheduledTimer(withTimeInterval: 3.0, repeats: true) { [weak self] _ in
            self?.poll()
        }
    }

    func stopPolling() {
        pollTimer?.invalidate()
        pollTimer = nil
    }

    func markExited() {
        DispatchQueue.main.async { self.status = .error }
        stopPolling()
    }

    private func poll() {
        guard let pid = neonPane?.shellPID, pid > 0 else { return }
        let children = AIDetector.childNames(of: pid)
        DispatchQueue.main.async {
            if children.isEmpty {
                self.status = .idle
            } else if children.contains(where: AIDetector.isAI) {
                self.status = .ai
            } else {
                self.status = .running
            }
        }
    }

    deinit { stopPolling() }
}
