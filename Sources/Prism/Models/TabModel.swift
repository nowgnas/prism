import SwiftUI
import Combine

final class TabModel: ObservableObject, Identifiable {
    let id = UUID()
    @Published var name: String
    @Published var sessions: [SessionModel] = []
    @Published var activeSessionIndex: Int = 0
    @Published var status: TerminalStatus = .idle

    let color: Color
    let paletteIndex: Int

    private var cancellables = Set<AnyCancellable>()

    init(paletteIndex: Int) {
        self.paletteIndex = paletteIndex
        self.color = NeonTheme.color(for: paletteIndex)
        self.name = "tab \(paletteIndex + 1)"
        addSession()
    }

    // MARK: - Session management

    func addSession() {
        guard sessions.count < 6 else { return }
        let session = SessionModel()
        let colorIndex = (paletteIndex + sessions.count) % NeonTheme.palette.count
        session.paneColor = NeonTheme.color(for: colorIndex)
        sessions.append(session)
        observe(session)
        // Auto-focus the new pane
        activeSessionIndex = sessions.count - 1
    }

    func removeSession(at index: Int) {
        guard sessions.count > 1, index < sessions.count else { return }
        sessions[index].stopPolling()
        sessions.remove(at: index)
        activeSessionIndex = min(activeSessionIndex, sessions.count - 1)
        recomputeStatus()
    }

    var activeSession: SessionModel? {
        guard activeSessionIndex < sessions.count else { return nil }
        return sessions[activeSessionIndex]
    }

    // MARK: - Pane focus navigation

    func nextPane() {
        guard !sessions.isEmpty else { return }
        activeSessionIndex = (activeSessionIndex + 1) % sessions.count
    }

    func prevPane() {
        guard !sessions.isEmpty else { return }
        activeSessionIndex = (activeSessionIndex - 1 + sessions.count) % sessions.count
    }

    // MARK: - Cleanup

    func cleanup() {
        sessions.forEach { $0.stopPolling() }
    }

    // MARK: - Status

    private func observe(_ session: SessionModel) {
        session.$status
            .receive(on: DispatchQueue.main)
            .sink { [weak self] _ in self?.recomputeStatus() }
            .store(in: &cancellables)
    }

    private func recomputeStatus() {
        status = sessions.map(\.status).max() ?? .idle
    }
}
