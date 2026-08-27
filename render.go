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

	for _, s := range w.Sheep {
		grid[s.Y][s.X] = 'S'
	}

	fmt.Printf("Day: %d\n", w.Day)
	for _, row := range grid {
		fmt.Println(string(row))
	}
}
