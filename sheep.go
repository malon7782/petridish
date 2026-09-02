package main

type Sheep struct {
	Y, X int
}

func (s *Sheep) Simulate(w *World) {
	dx := w.Rng.Intn(3) - 1
	dy := w.Rng.Intn(3) - 1
	newX := s.X + dx
	newY := s.Y + dy

	if newX >= 0 && newX < w.Width && newY >= 0 && newY < w.Height {
		if w.Map[newY][newX].Height >= 2 || w.Lakes[newY][newX] != nil {
			return
		} else {
			s.Y = newY
			s.X = newX
		}
	}
	// this message is for demo purposes and is indeed redundant.
	// to be replaced with real events like birth and death of sheep
	//	w.Logger.Add(w.Day, fmt.Sprintf("Day %d: Sheep moved.", w.Day))
}

func (s *Sheep) Pos() (int, int) {
	return s.Y, s.X
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
		ny := w.Rng.Intn(w.Height)
		nx := w.Rng.Intn(w.Width)
		for w.Map[ny][nx].Height >= 2 || w.Lakes[ny][nx] != nil {
			ny = w.Rng.Intn(w.Height)
			nx = w.Rng.Intn(w.Width)
		}
		w.Entities = append(w.Entities, &Sheep{Y: ny, X: nx})
	}
}
