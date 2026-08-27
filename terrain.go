package main

import "fmt"

type Mountain struct {
	X, Y     int
	Height   int
	ColorStr string
}

func (m *Mountain) Simulate(w *World) {}

func (m *Mountain) Pos() (int, int) {
	return m.X, m.Y
}

func (m *Mountain) Icon() byte {
	if m.Height <= 0 {
		return '0'
	}
	if m.Height > 9 {
		return '9'
	}
	return byte('0' + m.Height)
}

func (m *Mountain) Layer() int {
	return 1
}

func (m *Mountain) Color() string {
	return m.ColorStr
}

// -----------------------------------------------
// world init related

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

//Needswork: this part is ai-generated. Perhaps we need a better algo?
//这个函数必须保证颜色与高度对应，因为render端只负责分别绘制字符与颜色

func generateTerrain(w *World, numPeaks int) {
	type peak struct {
		x, y, h int
	}

	var peaks []peak
	for i := 0; i < numPeaks; i++ {
		peaks = append(peaks, peak{
			x: w.Rng.Intn(w.Width),
			y: w.Rng.Intn(w.Height),
			h: w.Rng.Intn(5) + 5,
		})
	}

	for y := 0; y < w.Height; y++ {
		for x := 0; x < w.Width; x++ {
			maxHeight := -999
			for _, p := range peaks {
				distance := abs(x-p.x) + abs(y-p.y)
				h := p.h - distance

				if h > maxHeight {
					maxHeight = h
				}
			}
			if maxHeight < 0 {
				maxHeight = 0
			}
			ColorCode := 240
			if maxHeight > 0 {
				ColorCode = ColorCode + (maxHeight-1)*2
			}

			w.Entities = append(w.Entities, &Mountain{
				X:      x,
				Y:      y,
				Height: maxHeight,
				// Calculate the color
				ColorStr: fmt.Sprintf("\033[38;5;%dm", ColorCode),
			})
		}
	}
}
