package main

import "fmt"

type Lake struct {
	Height float64
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

	ColorStr := ""

	if l.Height >= 1 {
		ColorStr = fmt.Sprintf("\033[38;5;%dm", 33)
	} else if l.Height >= 0.6 {
		ColorStr = fmt.Sprintf("\033[38;5;%dm", 39)
	} else if l.Height >= 0.2 {
		ColorStr = fmt.Sprintf("\033[38;5;%dm", 45)
	}
	return ColorStr
}

func isFlatDry(w *World, y, x int) bool {
	return w.Map[y][x].Height <= 1 && w.Lakes[y][x].Height == 0.0
}
// a river starts from one of the boundaries
// should start at flat ground
// find longest path
func generateRiver(w *World) {
	// a river source is a random flat boundary cell that is not water yet.
	findsrc := func(w *World) pair {
		starts := []pair{}
		for i := 0; i < w.Height; i++ {
			delta := 1
			if 1 < i && i+1 < w.Height { delta = w.Width-1 }
			for j := 0; j < w.Width; j += delta {
				if isFlatDry(w, i, j) { starts = append(starts, pair{y: i, x: j}) }
			}
		}
		if len(starts) == 0 {
			return pair{y: -1, x: -1}
		}
		return starts[w.Rng.Intn(len(starts))]
	}
	src := findsrc(w)
	generateRiverAt(w, src.x, src.y)
}

func _whichside(w *World, p pair) int {
	switch {
	case p.y == 0:
		return 0
	case p.y == w.Height-1:
		return 1
	case p.x == 0:
		return 2
	case p.x == w.Width-1:
		return 3
	}
	return -1
}

func _rivercost(w *World, y, x int) float64 {
	if w.Lakes[y][x].Height > 0.0 {
		return 0.5
	}
	h := w.Map[y][x].Height
	cost := 1.0
	switch {
	case h <= 1:
	case h == 2:
		cost += 6
	case h == 3:
		cost += 20
	default:
		cost += 60 * float64(h-3)
	}
	// running along the map edge looks odd, so nudge the river inland
	if _whichside(w, pair{y: y, x: x}) >= 0 {
		cost += 2
	}
	// noise makes the river wander on flats
	return cost + 1.5*w.Rng.Float64()
}

// a river has a source. it flows from its source all the way into
// a mountain or the occean.
// the occean for now simply means outside the map
//
// to describe the gen algorithm:
//   - dijkstra from the source. stepping onto a cell costs more the higher it is,
//     plus a bit of noise so the river meanders
//   - among the flat boundary cells on other sides, pick the mouth with the longest route
//   - fill the route with water. depth and width grow towards the mouth
//
// the algo ends when no boundary cell on another side is reachable
func generateRiverAt(w *World, srcx, srcy int) {
	src := pair{y: srcy, x: srcx}
	if src.y < 0 {
		return
	}
	dirs4 := [4][2]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}
	inMap := func(y, x int) bool {
		return y >= 0 && x >= 0 && y < w.Height && x < w.Width
	}

	// dijkstra. the map is tiny, so a plain O(n^2) scan for the next cell is
	// perfectly fine and saves us a heap.
	n := w.Width * w.Height
	id := func(p pair) int { return p.y*w.Width + p.x }
	dist := make([]float64, n)
	steps := make([]int, n)
	prev := make([]int, n)
	done := make([]bool, n)
	for i := range dist {
		dist[i] = -1
		prev[i] = -1
	}
	dist[id(src)] = 0

	for {
		cur := -1
		for i := 0; i < n; i++ {
			if done[i] || dist[i] < 0 {
				continue
			}
			if cur < 0 || dist[i] < dist[cur] {
				cur = i
			}
		}
		if cur < 0 {
			break
		}
		done[cur] = true
		cy, cx := cur/w.Width, cur%w.Width
		for _, d := range dirs4 {
			ny, nx := cy+d[0], cx+d[1]
			if !inMap(ny, nx) {
				continue
			}
			nxt := id(pair{y: ny, x: nx})
			if done[nxt] {
				continue
			}
			c := dist[cur] + _rivercost(w, ny, nx)
			if dist[nxt] < 0 || c < dist[nxt] {
				dist[nxt] = c
				steps[nxt] = steps[cur] + 1
				prev[nxt] = cur
			}
		}
	}

	// the mouth: a flat boundary cell on another side, reached by the longest
	// route. ties are broken randomly.
	srcSide := _whichside(w, src)
	mouth := -1
	for y := 0; y < w.Height; y++ {
		for x := 0; x < w.Width; x++ {
			p := pair{y: y, x: x}
			side := _whichside(w, p)
			if side < 0 || side == srcSide || !isFlatDry(w, y, x) {
				continue
			}
			i := id(p)
			if dist[i] < 0 {
				continue
			}
			if mouth < 0 || steps[i] > steps[mouth] ||
				(steps[i] == steps[mouth] && w.Rng.Intn(2) == 0) {
				mouth = i
			}
		}
	}
	if mouth < 0 {
		return
	}

	// walk back from the mouth to the source
	route := []pair{}
	for i := mouth; i >= 0; i = prev[i] {
		route = append(route, pair{y: i / w.Width, x: i % w.Width})
	}
	for i, j := 0, len(route)-1; i < j; i, j = i+1, j-1 {
		route[i], route[j] = route[j], route[i]
	}

	fill := func(y, x int, h float64) {
		if w.Lakes[y][x].Height >= h {
			return
		}
		if w.Lakes[y][x].Height == 0.0 {
			// same trick as cumLake(): dig down a unit so the water has
			// somewhere to sit. only dig once per cell.
			w.Map[y][x].Height -= 1
		}
		w.Lakes[y][x] = &Lake{Height: h}
	}

	// ai-assisted
	// the height of the water grows along the route, so the river gets a
	// brighter color near the source and a deeper one near the mouth.
	// keep the source above the 0.2 threshold so the whole river renders.
	const srcDepth, mouthDepth = 0.35, 1.0
	for i, p := range route {
		t := 0.0
		if len(route) > 1 {
			t = float64(i) / float64(len(route)-1)
		}
		depth := srcDepth + (mouthDepth-srcDepth)*t
		fill(p.y, p.x, depth)

		// widen the river downstream. banks are a bit shallower than the
		// channel, and the further downstream, the more banks we flood.
		banks := 0
		switch {
		case t >= 0.75:
			banks = 2
		case t >= 0.35:
			banks = 1
		}
		for k := 0; k < banks; k++ {
			d := dirs4[w.Rng.Intn(len(dirs4))]
			ny, nx := p.y+d[0], p.x+d[1]
			if !inMap(ny, nx) || w.Map[ny][nx].Height > 1 || w.Lakes[ny][nx].Height > 0.0 {
				continue
			}
			if w.Rng.Intn(3) == 0 {
				continue
			}
			fill(ny, nx, depth*0.7)
		}
	}
}

func generateLake(w *World) {
	if w.Lakes == nil {
		w.Lakes = make([][]*Lake, w.Height)
		for y := 0; y < w.Height; y++ {
			w.Lakes[y] = make([]*Lake, w.Width)
		}
	}
	for y := range w.Height {
		for x := range w.Width {
			w.Lakes[y][x] = &Lake{Height: 0.0}
		}
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
					w.Lakes[ny][nx].Height == 0.0 && w.Map[ny][nx].Height < maxheight {
					// ensure at least genarate the first size/3 rounds of lake terrains
					if count < size/3 {
						queue = append(queue, pair{ny, nx})
					} else if roll := w.Rng.Intn(5); roll > 2 {
						queue = append(queue, pair{ny, nx})
					}
				}
				w.Lakes[cur.y][cur.x] = &Lake{Height: float64(1.0)}
			}
			w.Map[cur.y][cur.x].Height -= 1
		}
		count += 1
	}
}
