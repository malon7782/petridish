package main

import (
	"fmt"
)

// world is displayed as a []Cell array
type Cell struct {
	Char  byte
	Color string
}

func (w *World) renderWorld() {
	fmt.Print("\033[H")

	// initialize the ground (gray '.')
	grid := make([][]Cell, w.Height)
	for y := 0; y < w.Height; y++ {
		grid[y] = make([]Cell, w.Width)
		for x := 0; x < w.Width; x++ {
			grid[y][x] = Cell{
				Char:  w.Map[y][x].Icon(),
				Color: w.Map[y][x].Color(),
			}
		}
	}

	// iterate layers. currently, only []Sheeps are in layer 1.
	for layer := 1; layer <= 1; layer++ {
		for _, e := range w.Entities {
			if e.Layer() == layer {
				x, y := e.Pos()
				if x >= 0 && x < w.Width && y >= 0 && y < w.Height {
					grid[y][x] = Cell{
						Char:  e.Icon(),
						Color: e.Color(),
					}
				}
			}
		}
	}

	// print
	fmt.Printf("Day: %d\n", w.Day)
	var output string
	events := w.Logger.Events
	eventCount := w.Logger.Num
	for y := 0; y < w.Height; y++ {
		for x := 0; x < w.Width; x++ {
			c := grid[y][x]
			output += c.Color + string(c.Char) + "\033[0m"
		}
		// print log msgs
		idx := len(events) - eventCount + y
		if 0 <= idx && y < eventCount {
			output += "  " + events[idx].Msg
		}
		output += "\n"
	}
	fmt.Print(output)
}
