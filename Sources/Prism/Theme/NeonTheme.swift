import SwiftUI

enum NeonTheme {
    // Background colors
    static let background      = Color(hex: "0D0D0D")
    static let sidebarBg       = Color(hex: "111116")
    static let divider         = Color(hex: "1E1E2A")

    // Text
    static let textPrimary     = Color(hex: "E2E2F0")
    static let textDim         = Color(hex: "44445A")
    static let textMid         = Color(hex: "7777AA")

    // 6-color neon palette (tab accent colors)
    static let palette: [Color] = [
        Color(hex: "FF3B3B"), // 0 Neon Crimson
        Color(hex: "FF9F0A"), // 1 Neon Amber
        Color(hex: "32D74B"), // 2 Neon Jade
        Color(hex: "0A84FF"), // 3 Neon Cobalt
        Color(hex: "BF5AF2"), // 4 Neon Violet
        Color(hex: "FF375F"), // 5 Neon Rose
    ]

    static func color(for index: Int) -> Color {
        palette[index % palette.count]
    }

    // Terminal colors (NSColor for AppKit interop)
    static let termFG = NSColor(hex: "E2E2F0")!
    static let termBG = NSColor(hex: "0D0D0D")!

    // Border
    static let inactiveBorder = NSColor(white: 0.18, alpha: 1)
    static let borderWidth: CGFloat = 1.5
    static let cornerRadius: CGFloat = 6
}

// MARK: - Color helpers

extension Color {
    init(hex: String) {
        var s = hex.trimmingCharacters(in: CharacterSet.alphanumerics.inverted)
        var v: UInt64 = 0
        Scanner(string: s).scanHexInt64(&v)
        let r, g, b, a: UInt64
        switch s.count {
        case 3:
            (a, r, g, b) = (255, (v >> 8) * 17, (v >> 4 & 0xF) * 17, (v & 0xF) * 17)
        case 6:
            (a, r, g, b) = (255, v >> 16, v >> 8 & 0xFF, v & 0xFF)
        case 8:
            (a, r, g, b) = (v >> 24, v >> 16 & 0xFF, v >> 8 & 0xFF, v & 0xFF)
        default:
            (a, r, g, b) = (255, 0, 0, 0)
        }
        self.init(.sRGB,
                  red: Double(r) / 255,
                  green: Double(g) / 255,
                  blue: Double(b) / 255,
                  opacity: Double(a) / 255)
    }

    var nsColor: NSColor {
        NSColor(self)
    }
}

extension NSColor {
    convenience init?(hex: String) {
        var s = hex.trimmingCharacters(in: CharacterSet.alphanumerics.inverted)
        var v: UInt64 = 0
        guard Scanner(string: s).scanHexInt64(&v), s.count == 6 else { return nil }
        self.init(srgbRed: Double(v >> 16) / 255,
                  green:   Double(v >> 8 & 0xFF) / 255,
                  blue:    Double(v & 0xFF) / 255,
                  alpha:   1)
    }
}
