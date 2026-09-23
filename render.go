package main

import (
	"fmt"
	"os"
	"strings"
)

// world is displayed as a []Cell array
type Cell struct {
	Char  byte
	Color string
}

func (w *World) renderLog() {
	const (
		HORIZ = '─'
		VERT = '│'
		TL = '┌'
		TR = '┐'
		BL = '└'
		BR = '┘'
	)

	W := w.Width
	if W < 3 {
		panic(fmt.Sprintf("Map width is %d. Too small.", W))
	}

	events := w.Logger.Events
	events = events[max(0, len(events) - 5) : len(events)]
	var sb strings.Builder
	lines := 2                              // top and bottom border
	for _, ev := range events {
		lines += (len(ev.Msg) + W-3) / (W-2) // same ceiling division as nlines
	}
	sb.Grow(lines * (W + 5))

	sb.WriteRune(TL)
	for i := 1; i < W-1; i++ { sb.WriteRune(HORIZ) }
	sb.WriteRune(TR)
	sb.WriteByte('\n')

	for _, ev := range(events) {
		msg := ev.Msg
		nlines := (len(msg) + (W-2) - 1) / (W-2)
		for i := 0; i < nlines; i++ {
			sb.WriteRune(VERT)
			sb.WriteString(msg[i*(W-2) : min(len(msg), (i+1)*(W-2))])
			diff := (i+1)*(W-2) - len(msg)
			for diff > 0 { sb.WriteByte(' '); diff-- }
			sb.WriteRune(VERT)
			sb.WriteByte('\n')
		}
	}
	sb.WriteRune(BL)
	for i := 1; i < W-1; i++ { sb.WriteRune(HORIZ) }
	sb.WriteRune(BR)
	sb.WriteByte('\n')

	fmt.Print(sb.String())
}

func (w *World) renderWorld() {
	fmt.Print("\033[H")

	grid := make([][]Cell, w.Height)

	// iterate the terrains
	for y := 0; y < w.Height; y++ {
		grid[y] = make([]Cell, w.Width)
		for x := 0; x < w.Width; x++ {
			// mountain
			grid[y][x] = Cell{
				Char:  w.Map[y][x].Icon(),
				Color: w.Map[y][x].Color(),
			}
			// lake
			if w.Lakes != nil {
				w.drawLakeFlow()
			}
			if w.Lakes[y][x].Height > 0.2 { // Needswork: this check violates the architecture principle,
				// but now we're too lazy to refactor it.
				grid[y][x] = Cell{
					Char:  w.Lakes[y][x].Icon(),
					Color: w.Lakes[y][x].Color(),
				}
			}
			// grass
			if w.GrassMap != nil && w.GrassMap[y][x].IsAlive() {
				grid[y][x] = Cell{
					Char:  w.GrassMap[y][x].Icon(),
					Color: w.GrassMap[y][x].Color(),
				}
			}
		}
	}

	// iterate entity layers.
	// Needswork: to be honest, this part looks damn ugly
	for layer := 1; layer <= 1; layer++ {
		for _, e := range w.Entities {
			if e.Layer() == layer {
				y, x := e.Pos()
				if x >= 0 && x < w.Width && y >= 0 && y < w.Height {
					grid[y][x] = Cell{
						Char:  e.Icon(),
						Color: e.Color(),
					}
				}
			}
		}
	}

	// decoration
	w.decorateGrid(grid)

	// initialize output
	fmt.Printf("Day: %d Temperature: %.2f Rain: %d\n", w.Day, w.Weathers.Temperature, w.Weathers.RainLeft)
	var sb strings.Builder
	sb.Grow(w.Width * w.Height * 20)

	// joint everything
	for y := 0; y < w.Height; y++ {
		for x := 0; x < w.Width; x++ {
			c := grid[y][x]
			sb.WriteString(c.Color)
			sb.WriteByte(c.Char)
			sb.WriteString("\033[0m")
		}
		sb.WriteByte('\n')
	}
	fmt.Print(sb.String())

	if os.Getenv("PETRIDISH_DEBUG") != "" {
		for y := 0; y < w.Height; y++ {
			for x := 0; x < w.Width; x++ {
				fmt.Fprintf(os.Stderr, "%.0f ", w.Map[y][x].Height+w.Lakes[y][x].Height)
			}
			fmt.Fprintf(os.Stderr, "\n")
		}
		fmt.Fprintf(os.Stderr, "\n")
	}

	w.renderLog()
}
