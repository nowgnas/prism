import Foundation

enum StartupBanner {
    /// ANSI escape codes
    private static let reset   = "\u{001B}[0m"
    private static let bold    = "\u{001B}[1m"
    private static let clear   = "\u{001B}[2J\u{001B}[H"  // Clear screen & home cursor
    private static let hide    = "\u{001B}[?25l"          // Hide cursor
    private static let show    = "\u{001B}[?25h"          // Show cursor
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

    private static let dim   = "\u{001B}[38;2;68;68;90m"

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

        // Wipe any shell MOTD/prompt noise, hide cursor, draw from top-left
        frames.append((0, clear + hide + home))

        // Top divider
        frames.append((40, "\(dim)\(divider)\(reset)\n\n"))

        // Prism shape line by line
        for (i, line) in prismLines.enumerated() {
            let color = colors[i % colors.count]
            frames.append((70, "\(color)\(bold)\(line)\(reset)\n"))
        }

        frames.append((90, "\n"))

        // Typing effect: indent + one feed per glyph
        let title = "✦ P R I S M ✦"
        frames.append((45, "     "))
        var colorTick = 0
        for char in title {
            if char == " " {
                frames.append((28, " "))
            } else {
                let c = colors[colorTick % colors.count]
                colorTick += 1
                frames.append((45, "\(c)\(bold)\(char)\(reset)"))
            }
        }
        frames.append((60, "\n"))

        frames.append((90, "\(dim)    Terminal Multiplexer\(reset)\n"))
        frames.append((70, "\n\(dim)\(divider)\(reset)\n\n"))

        frames.append((120, show))

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
