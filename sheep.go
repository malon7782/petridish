package main

import "fmt"

const (
	MinLakeHeightForSheep = 0.20 // sheeps are safe to step onto lakes that are
	// shallower than 0.15 unit height
)

type Sheep struct {
	Y, X  int
	HP    float64
	Alive bool
}

func (s *Sheep) Pos() (int, int) { return s.Y, s.X }
func (s *Sheep) Icon() byte      { return 'S' }
func (s *Sheep) Layer() int      { return 1 }
func (s *Sheep) Color() string   { return "\033[38;5;255m" }
func (s *Sheep) IsAlive() bool   { return s.Alive }

func (s *Sheep) Simulate(w *World) {
	dx := w.Rng.Intn(3) - 1
	dy := w.Rng.Intn(3) - 1
	newX := s.X + dx
	newY := s.Y + dy

	if newX >= 0 && newX < w.Width && newY >= 0 && newY < w.Height {
		if w.Map[newY][newX].Height >= 2 || w.Lakes[newY][newX].Height > MinLakeHeightForSheep {
			return
		} else {
			s.Y = newY
			s.X = newX
		}
	}

	if w.GrassMap[s.Y][s.X].Alive {
		s.HP += w.GrassMap[s.X][s.Y].WaterContent / 10.0
		w.GrassMap[s.Y][s.X].Alive = false
		w.GrassMap[s.Y][s.X].WaterContent = 0.0
	}
	// Hunger?
	s.HP -= 5.0

	if s.HP < 0.0 {
		s.Alive = false
	}
	// this message is for demo purposes and is indeed redundant.
	// to be replaced with real events like birth and death of sheep
	w.Logger.Add(w.Day, fmt.Sprintf("Day %d: Sheep moved.", w.Day))
}

// map gen related

func generateSheep(w *World, num int) {
	for i := 0; i < num; i++ {
		ny := w.Rng.Intn(w.Height)
		nx := w.Rng.Intn(w.Width)
		for w.Map[ny][nx].Height >= 2 || w.Lakes[ny][nx].Height > MinLakeHeightForSheep {
			ny = w.Rng.Intn(w.Height)
			nx = w.Rng.Intn(w.Width)
		}
		w.Entities = append(w.Entities, &Sheep{Y: ny, X: nx, HP: 100, Alive: true})
	}
}
