package domain

type CellState rune

const (
	CellHidden CellState = '#'
	CellFlag   CellState = 'F'
	CellMine   CellState = '*'
	CellZero   CellState = '0'
)
