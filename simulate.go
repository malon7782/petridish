package main

func (w *World) simulateWorld() {
	for _, e := range w.Entities {
		e.Simulate(w)
	}
	w.simulateWeather()
	w.updateMoisture()
	w.simulateGrass()

	w.Day += 1
}
