package main

import "math/rand"

type Entity interface {
	Simulate(w *World)
	Pos() (x, y int)
	Icon() byte
	Layer() int
	Color() string
}

type World struct {
	Width  int
	Height int
	Day    int

	// for living spieces
	Entities []Entity
	// for (semi-)stationary elements of the world
	Map      [][]*Mountain
	Lakes    [][]*Lake
	Grass    [][]bool
	Moisture [][]float64

	// for the weather status
	Weathers *Weather

	// at least for now, this array is useful for liquid simulation
	MaxHeightMap [][]float64

	Logger *Logger
	Rng    *rand.Rand
}

func (w *World) updateMaxHeightMap() {
	// this shit doesn't deserve a standalone generation logic
	if w.MaxHeightMap == nil {
		w.MaxHeightMap = make([][]float64, w.Height)
		for y := 0; y < w.Height; y++ {
			w.MaxHeightMap[y] = make([]float64, w.Width)
		}
	}

	for y := range w.Height {
		for x := range w.Width {
			w.MaxHeightMap[y][x] = max(-42.42, w.Map[y][x].Height)
		}
	}
}
