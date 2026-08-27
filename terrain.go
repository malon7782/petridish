package main

import "fmt"

type Mountain struct {
	Height   int
	ColorStr string
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

func generateMountain(w *World, numPeaks int) {
	w.Map = make([][]*Mountain, w.Height)
	for y := 0; y < w.Height; y++ {
		w.Map[y] = make([]*Mountain, w.Width)
	}

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

			w.Map[y][x] = &Mountain{
				Height: maxHeight,
				// Calculate the color
				ColorStr: fmt.Sprintf("\033[38;5;%dm", ColorCode),
			}
		}
	}
}
