package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
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
	s := ""
	red := "\033[31m"
	reset := "\033[0m"

	if t&Fixed != 0 {
		s += red
	}

	t &= ^Fixed
	switch t {
	case Empty:
		s += " "
	case Left:
		s += "╴"
	case Right:
		s += "╶"
	case Up:
		s += "╵"
	case Down:
		s += "╷"
	case Left | Right:
		s += "─"
	case Up | Down:
		s += "│"
	case Left | Up:
		s += "┘"
	case Left | Down:
		s += "┐"
	case Right | Up:
		s += "└"
	case Right | Down:
		s += "┌"
	case Left | Right | Up:
		s += "┴"
	case Left | Right | Down:
		s += "┬"
	case Left | Up | Down:
		s += "┤"
	case Right | Up | Down:
		s += "├"
	case Left | Right | Up | Down:
		s += "┼"
	}

	return s + reset
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
	// TODO possibbly buggy

	// left
	if x > 0 && m[x-1][y]&Fixed != 0 && rand.Intn(2) == 0 {
		m[x][y] &= ^Left
	}
	// right
	if x < len(m)-1 && m[x+1][y]&Fixed != 0 && rand.Intn(2) == 0 {
		m[x][y] &= ^Right
	}
	// up
	if y > 0 && m[x][y-1]&Fixed != 0 && rand.Intn(2) == 0 {
		m[x][y] &= ^Up
	}
	// down
	if y < len(m[x])-1 && m[x][y+1]&Fixed != 0 && rand.Intn(2) == 0 {
		m[x][y] &= ^Down
	}

	// fix the tile
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
			entropy := 4
			if i > 0 && m[i-1][j]&Fixed != 0 {
				entropy--
			}
			if i < len(m)-1 && m[i+1][j]&Fixed != 0 {
				entropy--
			}
			if j > 0 && m[i][j-1]&Fixed != 0 {
				entropy--
			}
			if j < len(m[i])-1 && m[i][j+1]&Fixed != 0 {
				entropy--
			}
			if t&Fixed == 0 {
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
	w := 20
	h := 10
	if len(os.Args) > 2 {
		wp, err := strconv.Atoi(os.Args[1])
		if err == nil {
			w = wp
		}
		hp, err := strconv.Atoi(os.Args[2])
		if err == nil {
			h = hp
		}
	}
	tmap := initMap(w, h)
	for !tmap.isFixed() {
		tmap = propagate(tmap)
		fmt.Println(tmap)
		tmap = collapse(tmap)
		fmt.Println(tmap)
	}
}
