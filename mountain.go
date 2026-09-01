package main

import (
	"fmt"
	"math"
	"sort"
)

// Needswork: Mountain should implement Terrain interface.

type Mountain struct {
	Height   int
	ColorStr string
}

func (m *Mountain) Icon() byte {
	if m.Height <= 0 {
		return '+'
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

// Needswork: this part is ai-generated. Perhaps we need a better algo?

func generateMountain(w *World, numPeaks int) {
	w.Map = make([][]*Mountain, w.Height)
	for y := 0; y < w.Height; y++ {
		w.Map[y] = make([]*Mountain, w.Width)
	}

	type peak struct {
		x, y           float64
		height, radius float64
	}
	type point struct {
		x, y float64
	}
	type edge struct {
		i, j int
	}
	type ridge struct {
		points        []point
		height, width float64
	}

	// try to make peaks not too close to each other
	// 贴贴，哒咩
	minDist := math.Min(5, math.Max(3, math.Min(float64(w.Width), float64(w.Height))/2))
	peaks := make([]peak, 0, numPeaks)
	for i := 0; i < numPeaks; i++ {
		best := peak{}
		bestDist := -1.0
		for attempt := 0; attempt < 40; attempt++ {
			candidate := peak{
				x:      float64(w.Rng.Intn(w.Width)),
				y:      float64(w.Rng.Intn(w.Height)),
				height: float64(w.Rng.Intn(5) + 5),
				radius: 2 + w.Rng.Float64()*1.5,
			}
			// peaks appear at boundaries for a chance of 1/4
			// this is to make them look like coming from outside the map
			if w.Rng.Intn(4) == 0 {
				if w.Rng.Intn(2) == 0 {
					candidate.x = float64((w.Width - 1) * w.Rng.Intn(2))
				} else {
					candidate.y = float64((w.Height - 1) * w.Rng.Intn(2))
				}
			}
			nearest := math.Inf(1)
			for _, other := range peaks {
				nearest = math.Min(nearest, math.Hypot(candidate.x-other.x, candidate.y-other.y))
			}
			if nearest > bestDist {
				bestDist = nearest
				best = candidate
			}
			if nearest >= minDist {
				break
			}
		}
		peaks = append(peaks, best)
	}

	edges := make([]edge, 0, len(peaks)*(len(peaks)-1)/2)
	for i := 0; i < len(peaks); i++ {
		for j := i + 1; j < len(peaks); j++ {
			edges = append(edges, edge{i: i, j: j})
		}
	}

	dist := func(i, j int) float64 {
		return math.Hypot(peaks[i].x-peaks[j].x, peaks[i].y-peaks[j].y)
	}
	sort.Slice(edges, func(i, j int) bool {
		return dist(edges[i].i, edges[i].j) < dist(edges[j].i, edges[j].j)
	})

	fa := make([]int, numPeaks)
	for i := range fa {
		fa[i] = i
	}
	find := func(i int) int {
		for i != fa[i] {
			fa[i] = fa[fa[i]]
			i = fa[i]
		}
		return i
	}
	isConnected := func(i, j int) bool {
		i = find(i)
		j = find(j)
		return i == j
	}

	mst := make([]edge, 0, len(peaks)-1)
	for _, e := range edges {
		if isConnected(e.i, e.j) {
			continue
		}
		mst = append(mst, e)
		fa[find(e.i)] = find(e.j)
		if len(mst) == len(peaks)-1 {
			break
		}
	}

	maxRemovable := min(len(mst), 2)
	removeCount := w.Rng.Intn(maxRemovable + 1)
	mst = mst[:len(mst)-removeCount]

	ridges := make([]ridge, 0, len(mst))
	clamp := func(value, limit float64) float64 {
		return math.Max(0, math.Min(value, limit))
	}
	for _, e := range mst {
		a, b := peaks[e.i], peaks[e.j]
		mid := point{
			x: clamp((a.x+b.x)/2+(w.Rng.Float64()*8-4), float64(w.Width-1)),
			y: clamp((a.y+b.y)/2+(w.Rng.Float64()*8-4), float64(w.Height-1)),
		}
		ridges = append(ridges, ridge{
			points: []point{{x: a.x, y: a.y}, mid, {x: b.x, y: b.y}},
			height: math.Min(a.height, b.height) * (0.4 + w.Rng.Float64()*0.15),
			width:  1 + w.Rng.Float64(),
		})
	}

	pointSegDist := func(p, a, b point) float64 {
		dx, dy := b.x-a.x, b.y-a.y
		lengthSquared := dx*dx + dy*dy
		if lengthSquared == 0 {
			return math.Hypot(p.x-a.x, p.y-a.y)
		}
		t := ((p.x-a.x)*dx + (p.y-a.y)*dy) / lengthSquared
		t = math.Max(0, math.Min(1, t))
		return math.Hypot(p.x-(a.x+t*dx), p.y-(a.y+t*dy))
	}

	// height field: h_i(x, y) = H_i * exp(-d^2 / 2r_i^2)
	// where d = hypot(x-x_i, y-y_i)
	for y := 0; y < w.Height; y++ {
		for x := 0; x < w.Width; x++ {
			cell := point{x: float64(x), y: float64(y)}
			maxPeak := 0.0
			peakSum := 0.0
			for _, p := range peaks {
				dx, dy := cell.x-p.x, cell.y-p.y
				height := p.height * math.Exp(-(dx*dx+dy*dy)/(2*p.radius*p.radius))
				peakSum += height
				maxPeak = math.Max(maxPeak, height)
			}

			ridgeHeight := 0.0
			for _, r := range ridges {
				for i := 0; i+1 < len(r.points); i++ {
					distance := pointSegDist(cell, r.points[i], r.points[i+1])
					height := r.height * math.Exp(-(distance*distance)/(2*r.width*r.width))
					ridgeHeight = math.Max(ridgeHeight, height)
				}
			}

			// mix max with sum, adding a bit noise
			heightValue := math.Max(maxPeak+0.12*(peakSum-maxPeak), ridgeHeight)
			heightValue *= 0.85 + 0.3*w.Rng.Float64()
			height := 0
			if heightValue >= 0.75 {
				height = int(math.Round(heightValue))
			}
			if height < 0 {
				height = 0
			} else if height > 9 {
				height = 9
			}

			type RGB struct{ R, G, B uint8 }

			var rockColors = []RGB{
				{57, 58, 63},
				{70, 71, 77},
				{84, 85, 91},
				{98, 99, 105},
				{112, 113, 119},
				{127, 128, 134},
				{142, 143, 149},
				{158, 159, 165},
				{175, 176, 182},
				{192, 193, 199},
			}

			// ground
			ColorCode := RGB{57, 58, 63}
			if height > 0 {
				h := height
				if h > 7 {
					h = 7
				}
				ColorCode = rockColors[h]
			}

			w.Map[y][x] = &Mountain{
				Height:   height,
				ColorStr: fmt.Sprintf("\x1b[38;2;%d;%d;%dm", ColorCode.R, ColorCode.G, ColorCode.B),
			}
		}
	}
}
