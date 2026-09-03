package main

import "fmt"

type Weather struct {
	Temperature         float64
	RainIntensity       float64 // defined as "change in lake height per unit time per unit area"
	RainTotal, RainLeft int
}

func generateWeather(w *World) {
	w.Weathers = &Weather{
		Temperature: 20.0, //baseline

		RainIntensity: 0.0,
		RainTotal:     0,
		RainLeft:      0}
}

func (w *World) simulateWeather() {
	// Temperature
	// w.Weathers.Temperature += w.Rng.Float64()*3.0 - 1.49 // temperature change ~ [-1.49, 1.51)
	// w.Logger.Add(w.Day, fmt.Sprintf("T: %d", int(w.Weathers.Temperature)))

	// Rain
	// - roll the duration
	if w.Weathers.RainLeft <= 0 && w.Rng.Intn(100) < 5 {
		w.Weathers.RainTotal = w.Rng.Intn(20) + 10 // last for 20 ~ 29 days
		w.Weathers.RainLeft = w.Weathers.RainTotal

		w.Logger.Add(w.Day, fmt.Sprintf("It's raining!"))
	}

	// - rain intensity algo
	if w.Weathers.RainLeft > 0 {
		w.Weathers.RainLeft--
		p := float64(w.Weathers.RainTotal-w.Weathers.RainLeft) / float64(w.Weathers.RainTotal)
		peak := w.Rng.Float64() * w.Rng.Float64() * 4
		w.Weathers.RainIntensity = peak * p * (1.0 - p)

		w.Logger.Add(w.Day, fmt.Sprintf("Days left: %d", w.Weathers.RainLeft))
	} else {
		w.Weathers.RainIntensity = 0.0
	}
}
