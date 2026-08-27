package main

import (
	"bytes"
	"fmt"
)

func (w *World) renderWorld() {
	fmt.Print("\033[H")

	grid := make([][]byte, w.Height)
	for i := range grid {
		grid[i] = bytes.Repeat([]byte{'.'}, w.Width)
	}

	// draw all entities

	for _, e := range w.Entities {
		x, y := e.Pos()
		if x >= 0 && x < w.Width && y >= 0 && y < w.Height {
			grid[y][x] = e.Icon()
		}
	}

	// stats

	fmt.Printf("Day: %d\n", w.Day)
	for _, row := range grid {
		fmt.Println(string(row))
	}
}
