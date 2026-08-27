package main

type Sheep struct {
	X, Y int
}

func (s *Sheep) Simulate(w *World) {
	dx := w.rng.Intn(3) - 1
	dy := w.rng.Intn(3) - 1
	newX := s.X + dx
	newY := s.Y + dy

	if newX >= 0 && newX < w.Width {
		s.X = newX
	}
	if newY >= 0 && newY < w.Height {
		s.Y = newY
	}
}

func (s *Sheep) Pos() (int, int) {
	return s.X, s.Y
}

func (s *Sheep) Icon() byte {
	return 'S'
}
