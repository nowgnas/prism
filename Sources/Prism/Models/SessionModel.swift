import SwiftUI
import Combine

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
final class SessionModel: ObservableObject, Identifiable {
    let id = UUID()
    @Published var status: TerminalStatus = .idle
    @Published var title: String = ""

    /// Shell PID — set by TerminalPaneView after process starts.
    var shellPID: pid_t = 0

    private var pollTimer: Timer?

    init() {}

    /// Start periodic process polling (3-second interval to keep memory/CPU low).
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
        guard shellPID > 0 else { return }
        let children = AIDetector.childNames(of: shellPID)
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
