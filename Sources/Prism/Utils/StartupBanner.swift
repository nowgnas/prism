import Foundation

enum StartupBanner {
    private static let reset = "\u{001B}[0m"
    private static let bold = "\u{001B}[1m"
    private static let clear = "\u{001B}[2J\u{001B}[H"
    private static let hide = "\u{001B}[?25l"
    private static let show = "\u{001B}[?25h"
    private static let home = "\u{001B}[H"

    private static let colors = [
        "\u{001B}[38;2;255;59;59m",
        "\u{001B}[38;2;255;159;10m",
        "\u{001B}[38;2;50;215;75m",
        "\u{001B}[38;2;10;132;255m",
        "\u{001B}[38;2;191;90;242m",
        "\u{001B}[38;2;255;55;95m",
    ]

    private static let dim = "\u{001B}[38;2;68;68;90m"
    private static let divider = "─────────────────────────────"

    private static let miniPrism = [
        "     ▲     ",
        "    ╱ ╲    ",
        "   ╱___╲   ",
    ]

    /// Clears MOTD noise, prints a short static banner, restores cursor — single feed, no animation.
    static var script: String {
        var s = clear + hide + home
        s += "\(dim)\(divider)\(reset)\n\n"
        for (i, line) in miniPrism.enumerated() {
            s += "\(colors[i % colors.count])\(bold)\(line)\(reset)\n"
        }
        s += "\n     \(colors[3])\(bold)PRISM\(reset) \(dim)· multiplexer\(reset)\n\n"
        s += "\(dim)\(divider)\(reset)\n\n"
        s += show
        return s
    }
}
