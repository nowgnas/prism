import Foundation

enum StartupBanner {
    /// ANSI escape codes
    private static let reset   = "\u{001B}[0m"
    private static let bold    = "\u{001B}[1m"
    private static let clear   = "\u{001B}[2J\u{001B}[H"  // Clear screen & home cursor
    private static let hide    = "\u{001B}[?25l"          // Hide cursor
    private static let show    = "\u{001B}[?25h"          // Show cursor
    private static let up      = "\u{001B}[A"             // Cursor up
    private static let home    = "\u{001B}[H"             // Cursor home

    // Neon rainbow colors (24-bit true color)
    private static let colors = [
        "\u{001B}[38;2;255;59;59m",    // 0 Crimson  #FF3B3B
        "\u{001B}[38;2;255;159;10m",   // 1 Amber    #FF9F0A
        "\u{001B}[38;2;50;215;75m",    // 2 Jade     #32D74B
        "\u{001B}[38;2;10;132;255m",   // 3 Cobalt   #0A84FF
        "\u{001B}[38;2;191;90;242m",   // 4 Violet   #BF5AF2
        "\u{001B}[38;2;255;55;95m",    // 5 Rose     #FF375F
    ]

    private static let white = "\u{001B}[38;2;226;226;240m"
    private static let dim   = "\u{001B}[38;2;68;68;90m"
    private static let cyan  = "\u{001B}[38;2;0;255;255m"

    /// Prism shape lines
    private static let prismLines = [
        "          ▲          ",
        "         ╱ ╲         ",
        "        ╱   ╲        ",
        "       ╱  ◇  ╲       ",
        "      ╱   ╱╲   ╲     ",
        "     ╱___╱__╲___╲    ",
    ]

    private static let divider = "───────────────────────────────────────────"

    /// Animation frames - each frame is an array of (delay_ms, text)
    static var animationFrames: [(delay: Int, text: String)] {
        var frames: [(Int, String)] = []

        // Frame 0: Hide cursor
        frames.append((0, hide))

        // Frame 1: Show divider
        frames.append((50, "\n\(dim)\(divider)\(reset)\n\n"))

        // Frames 2-7: Build prism line by line with color cycling
        for (i, line) in prismLines.enumerated() {
            let color = colors[i % colors.count]
            frames.append((80, "\(color)\(bold)\(line)\(reset)\n"))
        }

        // Frame 8: Empty line
        frames.append((100, "\n"))

        // Frame 9-20: Typing effect for "P R I S M"
        let title = "✦ P R I S M ✦"
        var titleStr = "     "
        for (i, char) in title.enumerated() {
            let colorIdx = i % colors.count
            titleStr += "\(colors[colorIdx])\(bold)\(char)\(reset)"
        }
        frames.append((50, titleStr + "\n"))

        // Frame 21: Subtitle
        frames.append((100, "\(dim)    Terminal Multiplexer\(reset)\n"))

        // Frame 22: Bottom divider
        frames.append((80, "\n\(dim)\(divider)\(reset)\n"))

        // Subtle glow pulse effect (no cursor movement - avoids timing issues)
        frames.append((200, ""))

        // Final frame: Show cursor and add newlines
        frames.append((50, "\n\(show)"))

        return frames
    }

    /// Static banner (fallback)
    static var banner: String {
        var result = "\n\(dim)\(divider)\(reset)\n\n"

        for (i, line) in prismLines.enumerated() {
            result += "\(colors[i])\(bold)\(line)\(reset)\n"
        }

        result += "\n"

        let title = "✦ P R I S M ✦"
        result += "     "
        for (i, char) in title.enumerated() {
            result += "\(colors[i % colors.count])\(bold)\(char)\(reset)"
        }
        result += "\n"

        result += "\(dim)    Terminal Multiplexer\(reset)\n"
        result += "\n\(dim)\(divider)\(reset)\n\n"

        return result
    }
}
