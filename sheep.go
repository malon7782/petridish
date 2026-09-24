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

// ++

func (s *Sheep) Roam(w *World, y, x int) (int, int) {
	dx := w.Rng.Intn(3) - 1
	dy := w.Rng.Intn(3) - 1
	newX := s.X + dx
	newY := s.Y + dy
	if newX >= 0 && newX < w.Width && newY >= 0 && newY < w.Height {
		if !(w.Map[newY][newX].Height >= 2 || w.Lakes[newY][newX].Height > MinLakeHeightForSheep) {
			return newY, newX
		}
	}
	return s.Y, s.X
}

// implemented via compact BFS search
func (s *Sheep) SeekFood(w *World, y, x int) (int, int) {
	radius := 4
	type node struct{ pos, first pair }
	visited := map[pair]bool{{y, x}: true}
	var queue []node
	for _, d := range dirs4 {
		ny, nx := y+d[0], x+d[1]
		if !w.inMap(ny, nx) || w.Map[ny][nx].Height >= 2 || w.Lakes[ny][nx].Height > MinLakeHeightForSheep {
			continue
		}
		p := pair{ny, nx}
		visited[p] = true
		queue = append(queue, node{p, p})
	}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		if w.GrassMap[cur.pos.y][cur.pos.x].Alive {
			return cur.first.y, cur.first.x
		}
		for _, d := range dirs4 {
			ny, nx := cur.pos.y+d[0], cur.pos.x+d[1]
			if ny < y-radius || ny > y+radius || nx < x-radius || nx > x+radius {
				continue
			}
			if !w.inMap(ny, nx) || w.Map[ny][nx].Height >= 2 || w.Lakes[ny][nx].Height > MinLakeHeightForSheep {
				continue
			}
			p := pair{ny, nx}
			if !visited[p] {
				visited[p] = true
				queue = append(queue, node{p, cur.first})
			}
		}
	}
	return s.Roam(w, y, x)
}

// ++

func (s *Sheep) EatGrass(w *World) {
	s.HP += w.GrassMap[s.Y][s.X].WaterContent
	if s.HP > 100.0 {
		s.HP = 100.0
	}
	w.GrassMap[s.Y][s.X].Alive = false
	w.GrassMap[s.Y][s.X].WaterContent = 0.0
}

// ++

var SheepMovements = map[string]func(*Sheep, *World, int, int) (int, int){
	"roam":     (*Sheep).Roam,
	"seekfood": (*Sheep).SeekFood,
}

var SheepBehaviors = map[string]func(*Sheep, *World){
	"eatgrass": (*Sheep).EatGrass,
}

// ++

func (s *Sheep) Simulate(w *World) {
	s.Y, s.X = SheepMovements["seekfood"](s, w, s.Y, s.X)

	if w.GrassMap[s.Y][s.X].Alive {
		SheepBehaviors["eatgrass"](s, w)
	}
	// Hunger?
	s.HP -= 2.0

	if s.HP < 0.0 {
		s.Alive = false
	}

	w.Logger.Add(w.Day, fmt.Sprintf("Sheep.HP %.2f", s.HP))
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
