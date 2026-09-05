package main

import (
	"fmt"
	"slices"
)

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

// a river starts from one of the boundaries
// should start at flat ground
// find longest path
func generateRiver(w *World) {
	// a river source is a random boundary cell
	find_endpoints := func() []pair {
		starts := []pair{}
		for i := 0; i < w.Height; i++ {
			delta := 1
			if 1 < i && i+1 < w.Height {
				delta = w.Width - 1
			}
			for j := 0; j < w.Width; j += delta {
				if w.Map[i][j].Height <= 4 {
					starts = append(starts, pair{y: i, x: j})
				}
			}
		}
		return starts
	}
	endpoints := find_endpoints()
	if len(endpoints) < 2 {
		return
	}
	w.Rng.Shuffle(len(endpoints), func(i, j int) {
		endpoints[i], endpoints[j] = endpoints[j], endpoints[i]
	})
	// sort by distance to corners
	// prefer corners
	cornerDist := func(p pair) int {
		dx := min(p.x, w.Width-1-p.x)
		dy := min(p.y, w.Height-1-p.y)
		return dx + dy
	}
	slices.SortFunc(endpoints, func(a, b pair) int {
		return cornerDist(a) - cornerDist(b)
	})

	start := endpoints[0]
	end := endpoints[1]
	Dist := func(a, b pair) int {
		return abs(a.x-b.x) + abs(b.y-a.y)
	}
	for i := range endpoints[2:] {
		if Dist(endpoints[2+i], start) > Dist(end, start) {
			end = endpoints[2+i]
		}
	}

	generateRiverBetween(w, start, end)
}

// use a slightly drunken walker
// penalize:
// - drifting away from the end terminal point
// - going from a lower cell to a higher cell
func _river_cost(w *World, from, to, end pair) float64 {
	height := w.Map[from.y][from.x].Height - w.Map[to.y][to.x].Height
	drift := abs(to.x-end.x) + abs(to.y-end.y)
	return float64(height*70 + drift*30)
}

func generateRiverBetween(w *World, start, end pair) {
	inMap := func(y, x int) bool {
		return y >= 0 && x >= 0 && y < w.Height && x < w.Width
	}
	// why does Go not have a usable heap?
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
	dist[id(start)] = 0
	dirs4 := [4][2]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}

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
			c := dist[cur] + _river_cost(w, pair{y: cy, x: cx}, pair{y: ny, x: nx}, end)
			if dist[nxt] < 0 || c < dist[nxt] {
				dist[nxt] = c
				steps[nxt] = steps[cur] + 1
				prev[nxt] = cur
			}
		}
	}

	// walk back from the end to the start
	route := []pair{}
	for i := id(end); i >= 0; i = prev[i] {
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
			banks = 5
		case t >= 0.35:
			banks = 4
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
