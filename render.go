package main

import (
	"fmt"
	"strings"
)

// world is displayed as a []Cell array
type Cell struct {
	Char  byte
	Color string
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
			if w.Lakes != nil && w.Lakes[y][x] != nil {
				grid[y][x] = Cell{
					Char:  w.Lakes[y][x].Icon(),
					Color: w.Lakes[y][x].Color(),
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

	// initialize output
	fmt.Printf("Day: %d\n", w.Day)
	var sb strings.Builder
	sb.Grow(w.Width * w.Height * 20)
	events := w.Logger.Events
	eventCount := w.Logger.Num

	// joint everything
	for y := 0; y < w.Height; y++ {
		for x := 0; x < w.Width; x++ {
			c := grid[y][x]
			sb.WriteString(c.Color)
			sb.WriteByte(c.Char)
			sb.WriteString("\033[0m")
		}
		// print log msgs
		idx := len(events) - eventCount + y
		if 0 <= idx && y < eventCount {
			//	output += "  " + events[idx].Msg
			sb.WriteString("  ")
			sb.WriteString(events[idx].Msg)
		}
		sb.WriteByte('\n')
	}
	fmt.Print(sb.String())
}
