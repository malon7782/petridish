package main

func (w *World) simulateWorld() {
	for _, e := range w.Entities {
		e.Simulate(w)
	}
	w.simulateWeather()
	w.updateMoisture()
	w.simulateGrass()
	w.simulateRiver()

	w.Day += 1
}
