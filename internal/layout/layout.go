package layout

// Rect defines the position and size of a terminal pane.
type Rect struct {
	X, Y, W, H int
}

const borderSize = 1 // 1-cell border on each side

// ComputeLayout returns the Rect slices for N panes within a W×H area.
// The returned rects include inner dimensions (border is drawn by the renderer).
func ComputeLayout(n, w, h int) []Rect {
	if n <= 0 {
		return nil
	}
	if n > 6 {
		n = 6
	}
	switch n {
	case 1:
		return []Rect{full(w, h)}
	case 2:
		return layout2(w, h)
	case 3:
		return layout3(w, h)
	case 4:
		return layout4(w, h)
	case 5:
		return layout5(w, h)
	case 6:
		return layout6(w, h)
	}
	return nil
}

func full(w, h int) Rect {
	return Rect{0, 0, w, h}
}

func layout2(w, h int) []Rect {
	// Left | Right
	half := w / 2
	return []Rect{
		{0, 0, half, h},
		{half, 0, w - half, h},
	}
}

func layout3(w, h int) []Rect {
	// Left (full height) | Right top / Right bottom
	half := w / 2
	halfH := h / 2
	return []Rect{
		{0, 0, half, h},
		{half, 0, w - half, halfH},
		{half, halfH, w - half, h - halfH},
	}
}

func layout4(w, h int) []Rect {
	// 2×2 grid
	hw := w / 2
	hh := h / 2
	return []Rect{
		{0, 0, hw, hh},
		{hw, 0, w - hw, hh},
		{0, hh, hw, h - hh},
		{hw, hh, w - hw, h - hh},
	}
}

func layout5(w, h int) []Rect {
	// Top row: 2 panes, Bottom row: 3 panes
	topH := h / 2
	botH := h - topH
	tw := w / 2
	bw := w / 3

	return []Rect{
		// top row
		{0, 0, tw, topH},
		{tw, 0, w - tw, topH},
		// bottom row
		{0, topH, bw, botH},
		{bw, topH, bw, botH},
		{bw * 2, topH, w - bw*2, botH},
	}
}

func layout6(w, h int) []Rect {
	// 2 rows × 3 cols
	rowH := h / 2
	colW := w / 3
	return []Rect{
		{0, 0, colW, rowH},
		{colW, 0, colW, rowH},
		{colW * 2, 0, w - colW*2, rowH},
		{0, rowH, colW, h - rowH},
		{colW, rowH, colW, h - rowH},
		{colW * 2, rowH, w - colW*2, h - rowH},
	}
}

// InnerRect returns the inner content area of a rect after accounting for borders.
func InnerRect(r Rect) Rect {
	return Rect{
		X: r.X + borderSize,
		Y: r.Y + borderSize,
		W: r.W - borderSize*2,
		H: r.H - borderSize*2,
	}
}
