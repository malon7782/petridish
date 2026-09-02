package main

import "fmt"

type Weather struct {
	RainIntensity       float64 // defined as "change in lake height per unit time per unit area"
	RainTotal, RainLeft int
}

func generateWeather(w *World) {
	w.Weathers = &Weather{RainIntensity: 0.0,
		RainTotal: 0,
		RainLeft:  0}
}

func (w *World) simulateWeather() {
	// Rain
	// - roll the duration
	if w.Weathers.RainLeft <= 0 && w.Rng.Intn(100) < 5 {
		w.Weathers.RainTotal = w.Rng.Intn(10) + 5 // last for 5 ~ 14 days
		w.Weathers.RainLeft = w.Weathers.RainTotal

		w.Logger.Add(w.Day, fmt.Sprintf("It's raining!"))
	}

	// - rain intensity algo
	if w.Weathers.RainLeft > 0 {
		w.Weathers.RainLeft--
		p := float64(w.Weathers.RainTotal-w.Weathers.RainLeft) / float64(w.Weathers.RainTotal)
		peak := w.Rng.Float64()
		w.Weathers.RainIntensity = peak * p * (1.0 - p)

		w.Logger.Add(w.Day, fmt.Sprintf("Days left: %d", w.Weathers.RainLeft))
	} else {
		w.Weathers.RainIntensity = 0.0
	}
}
