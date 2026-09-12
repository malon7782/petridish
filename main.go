package main

import (
	"flag"
	"fmt"
	"time"
)

func main() {
	fmt.Print("\033[2J")
	fmt.Print("\033[?25l")

	seed := flag.Int64("seed", time.Now().UnixNano(), "Seed of the world")
	maxDays := flag.Int("days", 1000, "Number of days")
	tickMs := flag.Int("tick", 50, "Rate of time")
	flag.Parse()

	fmt.Print("\033[2J")
	w := randomWorld(*seed, 40, 20)

	for w.Day < *maxDays {
		w.renderWorld()
		time.Sleep(time.Duration(*tickMs) * time.Millisecond)
		w.simulateWorld()
	}

	w.renderWorld()
	fmt.Println("Done.")
}
