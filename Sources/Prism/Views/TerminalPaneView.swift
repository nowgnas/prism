import SwiftUI
import SwiftTerm
import AppKit

// MARK: - NSViewRepresentable bridge

/// Wraps a cached NeonPane inside a thin TerminalContainer.
/// The NeonPane (and its LocalProcessTerminalView) is stored on SessionModel
/// so it persists across tab switches — tabs never lose their terminal state.
struct TerminalPaneView: NSViewRepresentable {
    @ObservedObject var session: SessionModel
    let tabColor: SwiftUI.Color
    let isActive: Bool

    func makeCoordinator() -> TerminalCoordinator {
        // Reuse existing coordinator if already created for this session
        if let existing = session.terminalCoordinator {
            return existing
        }
        let coordinator = TerminalCoordinator(session: session)
        session.terminalCoordinator = coordinator
        return coordinator
    }

    func makeNSView(context: Context) -> TerminalContainer {
        let container = TerminalContainer()
        let pane = demandPane(coordinator: context.coordinator)
        container.attach(pane: pane)
        return container
    }

    func updateNSView(_ container: TerminalContainer, context: Context) {
        let pane = demandPane(coordinator: context.coordinator)
        container.attach(pane: pane)
        pane.update(color: NSColor(tabColor), isActive: isActive)
    }

    // MARK: - Pane lifecycle (created once, reused forever)

    private func demandPane(coordinator: TerminalCoordinator) -> NeonPane {
        if let existing = session.neonPane {
            return existing
        }

        let tv = LocalProcessTerminalView(frame: .zero)

        // Appearance
        tv.nativeForegroundColor = NeonTheme.termFG
        tv.nativeBackgroundColor = NeonTheme.termBG
        tv.font = NSFont(name: "MesloLGS NF", size: 13)
            ?? NSFont(name: "JetBrains Mono", size: 13)
            ?? NSFont.monospacedSystemFont(ofSize: 13, weight: .regular)

        // Ghostty-like behaviour
        tv.optionAsMetaKey = true

        // Delegate
        tv.processDelegate = coordinator
        coordinator.terminalView = tv

        // Start shell
        let shell = ProcessInfo.processInfo.environment["SHELL"] ?? "/bin/zsh"
        tv.startProcess(executable: shell, args: [], environment: nil, execName: nil)

        let pane = NeonPane()
        pane.install(terminalView: tv, color: NSColor(tabColor), isActive: isActive)

        session.neonPane = pane
        session.startPolling()

        return pane
    }
}

// MARK: - TerminalCoordinator

final class TerminalCoordinator: NSObject, LocalProcessTerminalViewDelegate {
    let session: SessionModel
    weak var terminalView: LocalProcessTerminalView?

    init(session: SessionModel) {
        self.session = session
    }

    func sizeChanged(source: LocalProcessTerminalView, newCols: Int, newRows: Int) {}

    func setTerminalTitle(source: LocalProcessTerminalView, title: String) {
        DispatchQueue.main.async { self.session.title = title }
    }

    func processTerminated(source: TerminalView, exitCode: Int32?) {
        session.markExited()
    }

    func hostCurrentDirectoryUpdate(source: TerminalView, directory: String?) {}
    func requestOpenLink(source: TerminalView, link: String, params: [String: String]) {}
    func bell(source: TerminalView) {}
}

// MARK: - TerminalContainer (thin wrapper — recreated by SwiftUI, not preserved)

final class TerminalContainer: NSView {
    private(set) var pane: NeonPane?

    /// Re-parents the NeonPane into this container.
    /// AppKit silently removes from old superview before adding here.
    func attach(pane newPane: NeonPane) {
        guard newPane !== pane else { return }
        pane?.removeFromSuperview()
        pane = newPane
        addSubview(newPane)
        needsLayout = true
        // Restore focus to the terminal after re-parenting
        if let window = window {
            window.makeFirstResponder(newPane)
        }
    }

    override func layout() {
        super.layout()
        pane?.frame = bounds
    }

    override func viewDidMoveToWindow() {
        super.viewDidMoveToWindow()
        if window != nil, let pane = pane {
            window?.makeFirstResponder(pane)
        }
    }
}

// MARK: - NeonPane (neon border + glow container, owned by SessionModel)

final class NeonPane: NSView {
    private(set) var terminalView: LocalProcessTerminalView?
    var shellPID: pid_t { terminalView?.process.shellPid ?? 0 }

    override var acceptsFirstResponder: Bool { true }

    func install(terminalView tv: LocalProcessTerminalView, color: NSColor?, isActive: Bool) {
        self.terminalView = tv
        addSubview(tv)
        wantsLayer = true
        layer?.masksToBounds = false
        applyStyle(color: color, isActive: isActive)
    }

    func update(color: NSColor?, isActive: Bool) {
        CATransaction.begin()
        CATransaction.setDisableActions(true)
        applyStyle(color: color, isActive: isActive)
        CATransaction.commit()
    }

    private func applyStyle(color: NSColor?, isActive: Bool) {
        let borderColor = isActive
            ? (color ?? NeonTheme.inactiveBorder)
            : NeonTheme.inactiveBorder

        layer?.borderColor     = borderColor.cgColor
        layer?.borderWidth     = NeonTheme.borderWidth
        layer?.cornerRadius    = NeonTheme.cornerRadius
        layer?.backgroundColor = NeonTheme.termBG.cgColor

        if isActive {
            layer?.shadowColor   = borderColor.cgColor
            layer?.shadowRadius  = 8
            layer?.shadowOpacity = 0.7
            layer?.shadowOffset  = .zero
        } else {
            layer?.shadowOpacity = 0
        }
    }

    override func layout() {
        super.layout()
        let inset = NeonTheme.borderWidth + 1
        terminalView?.frame = bounds.insetBy(dx: inset, dy: inset)
    }

    override func keyDown(with event: NSEvent) {
        terminalView?.keyDown(with: event)
    }

    override func becomeFirstResponder() -> Bool {
        terminalView?.window?.makeFirstResponder(terminalView)
        return super.becomeFirstResponder()
    }
}
