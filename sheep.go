package main

type Sheep struct {
	X, Y int
}

func (w *World) handleSheep() {
	for i := range w.Sheep {
		dx := w.rng.Intn(3) - 1
		dy := w.rng.Intn(3) - 1
		newX := w.Sheep[i].X + dx
		newY := w.Sheep[i].Y + dy

		if newX >= 0 && newX < w.Width {
			w.Sheep[i].X = newX
		}
		if newY >= 0 && newY < w.Height {
			w.Sheep[i].Y = newY
		}
	}
}
