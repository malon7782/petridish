package main

var dirs4 = [4][2]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}

type pair struct {
	y int
	x int
}

type RGB struct {
	r, g, b int
}

func (w *World) inMap(y, x int) bool {
	return y >= 0 && x >= 0 && y < w.Height && x < w.Width
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
