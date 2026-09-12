package main

func (w *World) updateMaxHeightMap() {
	// this shit doesn't deserve a standalone generation logic
	if w.MaxHeightMap == nil {
		w.MaxHeightMap = make([][]float64, w.Height)
		for y := 0; y < w.Height; y++ {
			w.MaxHeightMap[y] = make([]float64, w.Width)
		}
	}

	for y := range w.Height {
		for x := range w.Width {
			w.MaxHeightMap[y][x] = max(-42.42, w.Map[y][x].Height)
		}
	}
}
