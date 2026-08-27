package main

import "fmt"

type Sheep struct {
	X, Y int
}

func (s *Sheep) Simulate(w *World) {
	dx := w.Rng.Intn(3) - 1
	dy := w.Rng.Intn(3) - 1
	newX := s.X + dx
	newY := s.Y + dy

	if newX >= 0 && newX < w.Width {
		s.X = newX
	}
	if newY >= 0 && newY < w.Height {
		s.Y = newY
	}

	// this message is for demo purposes and is indeed redundant.
	// to be replaced with real events like birth and death of sheep
	w.Logger.Add(w.Day, fmt.Sprintf("Day %d: Sheep moved.", w.Day))
}

func (s *Sheep) Pos() (int, int) {
	return s.X, s.Y
}

func (s *Sheep) Icon() byte {
	return 'S'
}

func (s *Sheep) Layer() int {
	return 1
}

func (s *Sheep) Color() string {
	return "\033[38;5;255m"
}

// map gen related

func generateSheep(w *World, num int) {
	for i := 0; i < num; i++ {
		w.Entities = append(w.Entities, &Sheep{X: w.Rng.Intn(w.Width),
			Y: w.Rng.Intn(w.Height)})
	}
}
