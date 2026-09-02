package main

import "fmt"

type Lake struct {
	Height   float64
	ColorStr string
}

// helper

type pair struct {
	y int
	x int
}

func (l *Lake) Icon() byte {
	if l == nil {
		return '%'
	}
	return '~'
}

func (l *Lake) Color() string {
	if l == nil {
		return ""
	}
	return l.ColorStr
}

// this is suitable for simulating rain. maybe it will live in rain.go or weather.go
// in the future.
func cumLake(w *World, numLakes int) {
	if w.Lakes == nil {
		w.Lakes = make([][]*Lake, w.Height)
		for y := 0; y < w.Height; y++ {
			w.Lakes[y] = make([]*Lake, w.Width)
		}
	}
	for i := 0; i < numLakes; i++ {
		randomy := w.Rng.Intn(w.Height)
		randomx := w.Rng.Intn(w.Width)

		w.Lakes[randomy][randomx] = &Lake{
			Height:   float64(1),
			ColorStr: fmt.Sprintf("\033[38;5;%dm", 45),
		}
		// this is important. we need to dig down a unit depth  so that we can
		//* place water. in case different lake gen functions access the same [][]*Mounatin
		// block, make sure that the height of that block is equal or greater than -1,
		w.Map[randomy][randomx].Height -= 1
	}
}

// generateBFSLake() takes a World instance and a `size` integar, then generate a approximately
// circuler lake on w.Lakes. the radius is not guaranteed (actually, intended to not guarantee)
// to reach `size`.

// when called, it searches the w.Map and find a block where height is no greater than 2
// as source. when expanding, it won't cover the mountains which have heights greater than 2
// either.
// unless, by a chance of 1/2, one lake will be located at some peaks.

// The edge of the lake has brighter blue color but no functional difference.

// It is implemented via BFS algo, as the names suggests.

func generateBFSLake(w *World, size int) {
	if w.Lakes == nil {
		w.Lakes = make([][]*Lake, w.Height)
		for y := 0; y < w.Height; y++ {
			w.Lakes[y] = make([]*Lake, w.Width)
		}
	}
	dirs := [4][2]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}

	count := 0
	inity := w.Rng.Intn(w.Height)
	initx := w.Rng.Intn(w.Width)
	isMntTop := w.Rng.Intn(4) == 0
	if isMntTop {
		// choose a peak for the mountaintop lake
		var peak *Mountain
		for i, row := range w.Map {
			for j, elem := range row {
				used := w.Lakes[i][j] != nil
				if used {
					used = w.Lakes[i][j].Height > 0
				}
				if (peak == nil || elem.Height > peak.Height) && !used {
					peak = elem
					inity = i
					initx = j
				}
			}
		}
	} else {
		for w.Map[inity][initx].Height > 2 {
			inity = w.Rng.Intn(w.Height)
			initx = w.Rng.Intn(w.Width)
		}
	}
	queue := []pair{pair{y: inity, x: initx}}

	maxheight := 2
	if isMntTop {
		size /= 2
		maxheight = 9
	}
	for len(queue) > 0 && count < size {
		sz := len(queue)
		for i := 0; i < sz; i++ {
			cur := queue[0]
			queue = queue[1:]
			for _, d := range dirs {
				ny := cur.y + d[0]
				nx := cur.x + d[1]
				if ny >= 0 && nx >= 0 && ny < w.Height && nx < w.Width &&
					w.Lakes[ny][nx] == nil && w.Map[ny][nx].Height < maxheight {
					// ensure at least genarate the first size/3 rounds of lake terrains
					if count < size/3 {
						queue = append(queue, pair{ny, nx})
					} else if roll := w.Rng.Intn(5); roll > 2 {
						queue = append(queue, pair{ny, nx})
					}
				}
				if float64(count) > 0.5*float64(size) {
					w.Lakes[cur.y][cur.x] = &Lake{Height: float64(1.0),
						ColorStr: fmt.Sprintf("\033[38;5;%dm", 33)}
				}
			}
			w.Map[cur.y][cur.x].Height -= 1
		}
		count += 1
	}
}
