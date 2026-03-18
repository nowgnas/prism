import SwiftUI

struct TabGridView: View {
    @ObservedObject var tab: TabModel

    var body: some View {
        GeometryReader { geo in
            let rects = gridLayout(count: tab.sessions.count, in: geo.size)
            ZStack(alignment: .topLeading) {
                ForEach(Array(tab.sessions.enumerated()), id: \.element.id) { idx, session in
                    if idx < rects.count {
                        let rect = rects[idx]
                        TerminalPaneView(
                            session: session,
                            tabColor: tab.color,
                            isActive: idx == tab.activeSessionIndex
                        )
                        .frame(width: rect.width, height: rect.height)
                        .offset(x: rect.minX, y: rect.minY)
                        .onTapGesture { tab.activeSessionIndex = idx }
                    }
                }
            }
        }
        .background(NeonTheme.background)
        .padding(4)
    }

    // MARK: - Layout algorithm (N = 1–6)

    private func gridLayout(count: Int, in size: CGSize) -> [CGRect] {
        let w = size.width - 8  // account for outer padding
        let h = size.height - 8
        let g: CGFloat = 4      // gap between panes

        switch count {
        case 1:
            return [CGRect(x: 0, y: 0, width: w, height: h)]

        case 2:
            let hw = (w - g) / 2
            return [
                CGRect(x: 0,      y: 0, width: hw,      height: h),
                CGRect(x: hw + g, y: 0, width: w - hw - g, height: h),
            ]

        case 3:
            let hw = (w - g) / 2
            let hh = (h - g) / 2
            return [
                CGRect(x: 0,      y: 0,      width: hw,          height: h),
                CGRect(x: hw + g, y: 0,      width: w - hw - g,  height: hh),
                CGRect(x: hw + g, y: hh + g, width: w - hw - g,  height: h - hh - g),
            ]

        case 4:
            let hw = (w - g) / 2
            let hh = (h - g) / 2
            return [
                CGRect(x: 0,      y: 0,      width: hw,          height: hh),
                CGRect(x: hw + g, y: 0,      width: w - hw - g,  height: hh),
                CGRect(x: 0,      y: hh + g, width: hw,          height: h - hh - g),
                CGRect(x: hw + g, y: hh + g, width: w - hw - g,  height: h - hh - g),
            ]

        case 5:
            let hh  = (h - g) / 2
            let tw  = (w - g) / 2
            let bw  = (w - 2 * g) / 3
            return [
                CGRect(x: 0,           y: 0,      width: tw,          height: hh),
                CGRect(x: tw + g,      y: 0,      width: w - tw - g,  height: hh),
                CGRect(x: 0,           y: hh + g, width: bw,          height: h - hh - g),
                CGRect(x: bw + g,      y: hh + g, width: bw,          height: h - hh - g),
                CGRect(x: bw * 2 + g * 2, y: hh + g, width: w - bw * 2 - g * 2, height: h - hh - g),
            ]

        case 6:
            let hh = (h - g) / 2
            let cw = (w - 2 * g) / 3
            return [
                CGRect(x: 0,            y: 0,      width: cw,              height: hh),
                CGRect(x: cw + g,       y: 0,      width: cw,              height: hh),
                CGRect(x: cw * 2 + g*2, y: 0,      width: w - cw*2 - g*2,  height: hh),
                CGRect(x: 0,            y: hh + g, width: cw,              height: h - hh - g),
                CGRect(x: cw + g,       y: hh + g, width: cw,              height: h - hh - g),
                CGRect(x: cw * 2 + g*2, y: hh + g, width: w - cw*2 - g*2,  height: h - hh - g),
            ]

        default:
            return [CGRect(x: 0, y: 0, width: w, height: h)]
        }
    }
}
