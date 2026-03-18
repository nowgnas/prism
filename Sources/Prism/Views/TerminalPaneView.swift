import SwiftUI
import SwiftTerm
import AppKit

// MARK: - SwiftUI bridge

struct TerminalPaneView: NSViewRepresentable {
    @ObservedObject var session: SessionModel
    let tabColor: SwiftUI.Color
    let isActive: Bool

    func makeNSView(context: Context) -> NeonPane {
        let pane = NeonPane()
        let tv = LocalProcessTerminalView(frame: .zero)

        // Appearance
        tv.nativeForegroundColor = NeonTheme.termFG
        tv.nativeBackgroundColor = NeonTheme.termBG
        tv.font = NSFont.monospacedSystemFont(ofSize: 13, weight: .regular)

        // Delegate
        tv.processDelegate = context.coordinator

        // Start shell
        let shell = ProcessInfo.processInfo.environment["SHELL"] ?? "/bin/zsh"
        tv.startProcess(executable: shell,
                        args: [],
                        environment: nil,
                        execName: nil)

        // Store reference for input forwarding
        context.coordinator.terminalView = tv
        session.startPolling()

        pane.install(terminalView: tv, color: NSColor(tabColor), isActive: isActive)
        return pane
    }

    func updateNSView(_ pane: NeonPane, context: Context) {
        pane.update(color: NSColor(tabColor), isActive: isActive)
    }

    func makeCoordinator() -> Coordinator {
        Coordinator(session: session)
    }
}

// MARK: - Coordinator

final class Coordinator: NSObject, LocalProcessTerminalViewDelegate {
    let session: SessionModel
    weak var terminalView: LocalProcessTerminalView?

    init(session: SessionModel) {
        self.session = session
    }

    // Required
    func sizeChanged(source: LocalProcessTerminalView, newCols: Int, newRows: Int) {}

    func setTerminalTitle(source: LocalProcessTerminalView, title: String) {
        DispatchQueue.main.async { self.session.title = title }
    }

    func processTerminated(source: TerminalView, exitCode: Int32?) {
        session.markExited()
    }

    // Optional (provide default implementations)
    func hostCurrentDirectoryUpdate(source: TerminalView, directory: String?) {}
    func requestOpenLink(source: TerminalView, link: String, params: [String: String]) {}
    func bell(source: TerminalView) {}
    func scrolled(source: TerminalView, position: Double) {}
    func rangeChanged(source: TerminalView, startY: Int, endY: Int) {}
}

// MARK: - NeonPane (NSView container with neon border + glow)

final class NeonPane: NSView {
    private var terminalView: LocalProcessTerminalView?

    override var isFlipped: Bool { true }
    override var acceptsFirstResponder: Bool { true }

    func install(terminalView tv: LocalProcessTerminalView, color: NSColor?, isActive: Bool) {
        self.terminalView = tv
        addSubview(tv)
        wantsLayer = true
        layer?.masksToBounds = false
        applyStyle(color: color, isActive: isActive)
    }

    func update(color: NSColor?, isActive: Bool) {
        applyStyle(color: color, isActive: isActive)
    }

    private func applyStyle(color: NSColor?, isActive: Bool) {
        let borderColor = isActive
            ? (color ?? NeonTheme.inactiveBorder)
            : NeonTheme.inactiveBorder

        layer?.borderColor  = borderColor.cgColor
        layer?.borderWidth  = NeonTheme.borderWidth
        layer?.cornerRadius = NeonTheme.cornerRadius
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

    // Forward key events to the embedded terminal view
    override func keyDown(with event: NSEvent) {
        terminalView?.keyDown(with: event)
    }

    override func becomeFirstResponder() -> Bool {
        terminalView?.window?.makeFirstResponder(terminalView)
        return super.becomeFirstResponder()
    }
}
