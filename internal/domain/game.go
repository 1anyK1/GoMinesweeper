package domain

import (
	"errors"
	"math/rand"
	"strings"
)

const (
	Hidden = rune(CellHidden)
	Flag   = rune(CellFlag)
	Mine   = rune(CellMine)
	Zero   = rune(CellZero)
)

type GameStatus int

const (
	StatusPlaying GameStatus = iota
	StatusWin
	StatusLose
)

type Game struct {
	Size       int
	Mines      int
	Field      []rune
	Visible    []rune
	FirstClick bool
	Status     GameStatus
}

func NewGame(size, mines int) (*Game, error) {
	if size <= 0 {
		return nil, errors.New("size must be positive")
	}
	if mines < 0 || mines >= size*size {
		return nil, errors.New("invalid mines count")
	}

	g := &Game{
		Size:       size,
		Mines:      mines,
		Field:      make([]rune, size*size),
		Visible:    make([]rune, size*size),
		FirstClick: true,
		Status:     StatusPlaying,
	}

	for i := range g.Visible {
		g.Visible[i] = Hidden
	}

	return g, nil
}

func (g *Game) Open(x, y int) error {
	if !g.inBounds(x, y) {
		return errors.New("coordinates out of bounds")
	}
	if g.Status != StatusPlaying {
		return errors.New("game is already finished")
	}

	if g.FirstClick {
		g.generateField(x, y)
		g.FirstClick = false
	}

	idx := g.index(x, y)

	if g.Visible[idx] == Flag {
		return nil
	}

	if g.Field[idx] == Mine {
		g.openMines()
		g.Status = StatusLose
		return nil
	}

	g.openCells(x, y)

	return nil
}

func (g *Game) ToggleFlag(x, y int) error {
	if !g.inBounds(x, y) {
		return errors.New("coordinates out of bounds")
	}
	if g.Status != StatusPlaying {
		return errors.New("game is already finished")
	}

	idx := g.index(x, y)

	if g.Visible[idx] != Hidden && g.Visible[idx] != Flag {
		return nil
	}

	if g.Visible[idx] == Flag {
		g.Visible[idx] = Hidden
	} else {
		g.Visible[idx] = Flag
	}

	if g.checkWinByFlags() {
		g.openMines()
		g.Status = StatusWin
	}

	return nil
}

func (g *Game) Reset() {
	for i := range g.Visible {
		g.Visible[i] = Hidden
	}

	for i := range g.Field {
		g.Field[i] = Zero
	}

	g.FirstClick = true
	g.Status = StatusPlaying
}

func (g *Game) Render() string {
	var b strings.Builder

	for x := 0; x < g.Size; x++ {
		for y := 0; y < g.Size; y++ {
			b.WriteRune(g.Visible[g.index(x, y)])
			b.WriteRune(' ')
		}
		b.WriteRune('\n')
	}

	return b.String()
}

func (g *Game) generateField(firstX, firstY int) {
	for {
		g.initField()
		g.placeMines()
		g.calculateNumbers()

		if g.Field[g.index(firstX, firstY)] == Zero {
			return
		}
	}
}

func (g *Game) initField() {
	for i := range g.Field {
		g.Field[i] = Zero
	}
}

func (g *Game) placeMines() {
	placed := 0

	for placed < g.Mines {
		x := rand.Intn(g.Size)
		y := rand.Intn(g.Size)
		idx := g.index(x, y)

		if g.Field[idx] != Mine {
			g.Field[idx] = Mine
			placed++
		}
	}
}

func (g *Game) calculateNumbers() {
	dx := []int{-1, -1, -1, 0, 0, 1, 1, 1}
	dy := []int{-1, 0, 1, -1, 1, -1, 0, 1}

	for x := 0; x < g.Size; x++ {
		for y := 0; y < g.Size; y++ {
			idx := g.index(x, y)

			if g.Field[idx] == Mine {
				continue
			}

			count := 0

			for k := 0; k < 8; k++ {
				nx := x + dx[k]
				ny := y + dy[k]

				if g.inBounds(nx, ny) && g.Field[g.index(nx, ny)] == Mine {
					count++
				}
			}

			g.Field[idx] = rune('0' + count)
		}
	}
}

func (g *Game) openCells(x, y int) {
	if !g.inBounds(x, y) {
		return
	}

	idx := g.index(x, y)

	if g.Visible[idx] != Hidden {
		return
	}

	g.Visible[idx] = g.Field[idx]

	if g.Field[idx] != Zero {
		return
	}

	dx := []int{-1, -1, -1, 0, 0, 1, 1, 1}
	dy := []int{-1, 0, 1, -1, 1, -1, 0, 1}

	for k := 0; k < 8; k++ {
		g.openCells(x+dx[k], y+dy[k])
	}
}

func (g *Game) openMines() {
	for i := range g.Field {
		if g.Field[i] == Mine {
			g.Visible[i] = Mine
		}
	}
}

func (g *Game) checkWinByFlags() bool {
	correctFlags := 0

	for i := range g.Field {
		if g.Field[i] == Mine && g.Visible[i] == Flag {
			correctFlags++
		}
	}

	return correctFlags == g.Mines
}

func (g *Game) inBounds(x, y int) bool {
	return x >= 0 && x < g.Size && y >= 0 && y < g.Size
}

func (g *Game) index(x, y int) int {
	return x*g.Size + y
}
