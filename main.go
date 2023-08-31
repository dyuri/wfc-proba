package main

import (
	"fmt"
	"math/rand"
)

type Tile int
type Map [][]Tile
type Point struct {
	x int
	y int
}

const (
	Empty Tile = 0
	Left  Tile = 1 << iota
	Right
	Up
	Down
	Fixed Tile = 1 << 10
)

const InitialTile = Empty | Left | Right | Up | Down

func (t Tile) String() string {
	t &= ^Fixed
	switch t {
	case Empty:
		return " "
	case Left:
		return "╴"
	case Right:
		return "╶"
	case Up:
		return "╵"
	case Down:
		return "╷"
	case Left | Right:
		return "─"
	case Up | Down:
		return "│"
	case Left | Up:
		return "┘"
	case Left | Down:
		return "┐"
	case Right | Up:
		return "└"
	case Right | Down:
		return "┌"
	case Left | Right | Up:
		return "┴"
	case Left | Right | Down:
		return "┬"
	case Left | Up | Down:
		return "┤"
	case Right | Up | Down:
		return "├"
	case Left | Right | Up | Down:
		return "┼"
	}
	return fmt.Sprintf("?%d?", t)
}

func (m Map) String() string {
	w := len(m[0])
	h := len(m)
	s := ""
	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			s += m[i][j].String()
		}
		s += "\n"
	}
	return s
}

func initMap(w, h int) Map {
	m := make([][]Tile, h)
	for i := range m {
		m[i] = make([]Tile, w)
		for j := range m[i] {
			m[i][j] = InitialTile
		}
	}
	return m
}

func propagate(m Map) Map {
	for i := range m {
		for j := range m[i] {
			t := m[i][j]
			if t&Fixed == 0 {
				top := Empty
				bottom := Empty
				left := Empty
				right := Empty
				if i > 0 {
					top = m[i-1][j]
				}
				if i < len(m)-1 {
					bottom = m[i+1][j]
				}
				if j > 0 {
					left = m[i][j-1]
				}
				if j < len(m[i])-1 {
					right = m[i][j+1]
				}

				if t&Up != 0 && top&Down == 0 {
					m[i][j] &= ^Up
				}
				if t&Down != 0 && bottom&Up == 0 {
					m[i][j] &= ^Down
				}
				if t&Left != 0 && left&Right == 0 {
					m[i][j] &= ^Left
				}
				if t&Right != 0 && right&Left == 0 {
					m[i][j] &= ^Right
				}

				// Epmty is Fixed
				if m[i][j] == Empty {
					m[i][j] |= Fixed
				}
			}
		}
	}
	return m
}

func fixTile(m Map, x, y int) Map {
	// choose directions randomly
	randomTile := Empty
	// left
	if rand.Intn(2) == 0 {
		randomTile |= Left
	}
	// right
	if rand.Intn(2) == 0 {
		randomTile |= Right
	}
	// up
	if rand.Intn(2) == 0 {
		randomTile |= Up
	}
	// down
	if rand.Intn(2) == 0 {
		randomTile |= Down
	}

	// fix the tile
	m[x][y] &= randomTile
	m[x][y] |= Fixed

	return m
}

func fixRandomTile(m Map) Map {
	x := rand.Intn(len(m))
	y := rand.Intn(len(m[x]))
	m = fixTile(m, x, y)

	return m
}

func collapse(m Map) Map {

	minEntropy := 10
	minTiles := []Point{}

	// calculate entropy
	for i := range m {
		for j := range m[i] {
			t := m[i][j]
			// only if has fixed neighbors
			hasFixedNeighbor := false
			if i > 0 && m[i-1][j]&Fixed != 0 {
				hasFixedNeighbor = true
			}
			if i < len(m)-1 && m[i+1][j]&Fixed != 0 {
				hasFixedNeighbor = true
			}
			if j > 0 && m[i][j-1]&Fixed != 0 {
				hasFixedNeighbor = true
			}
			if j < len(m[i])-1 && m[i][j+1]&Fixed != 0 {
				hasFixedNeighbor = true
			}
			if hasFixedNeighbor && t&Fixed == 0 {
				entropy := 0
				if t&Up != 0 {
					entropy++
				}
				if t&Down != 0 {
					entropy++
				}
				if t&Left != 0 {
					entropy++
				}
				if t&Right != 0 {
					entropy++
				}

				if entropy < minEntropy {
					minEntropy = entropy
					minTiles = []Point{}
				}

				if entropy == minEntropy {
					minTiles = append(minTiles, Point{i, j})
				}
			}
		}
	}

	if len(minTiles) == 0 {
		m = fixRandomTile(m)
		return m
	}

	// choose one tile randomly
	p := minTiles[rand.Intn(len(minTiles))]
	m = fixTile(m, p.x, p.y)

	return m
}

func (m Map) isFixed() bool {
	for i := range m {
		for j := range m[i] {
			if m[i][j]&Fixed == 0 {
				return false
			}
		}
	}
	return true
}

func main() {
	tmap := initMap(20, 10)
	for !tmap.isFixed() {
		tmap = propagate(tmap)
		tmap = collapse(tmap)
	}
	fmt.Println(tmap)
}
