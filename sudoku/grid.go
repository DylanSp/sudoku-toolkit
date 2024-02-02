package sudoku

import (
	"slices"

	"github.com/samber/lo"
)

// nil == no value
// naive representation; may replace later
type Cell *int

type Grid struct {
	baseSize uint

	// flat slice of cells, row-by-row; may change representation later
	// invariant: len(cells) == baseSize ^ 4
	cells []Cell
}

func EmptyGrid(baseSize uint) Grid {
	grid := Grid{
		baseSize: baseSize,
	}

	// example: a 9x9 Sudoku has baseSize 3: each side of the grid is 3*3 = 9 cells
	// the total grid is 9*9 = 81 = 3^4 cells
	// use repeated multiplication instead of math.Pow() to avoid converting to/from float64
	grid.cells = make([]Cell, baseSize*baseSize*baseSize*baseSize)

	return grid
}

func (g *Grid) sideLength() uint {
	return g.baseSize * g.baseSize
}

// each house in the grid must have all elements in the range from 1 to maxElement(), inclusive
// this is the same calculation as sideLength(), but split into a separate method for clarity
func (g *Grid) maxElement() int {
	return int(g.baseSize * g.baseSize)
}

func (g *Grid) cellAt(row uint, col uint) Cell {
	// could inline these into a single line, but gofmt's arithmetic formatting makes the inlined version less clear
	rowBaseIndex := row * g.sideLength()
	index := rowBaseIndex + col
	return g.cells[index]
}

// TODO - this is currently the only way to set values in a grid; probably want other constructor(s), may not want this direct setter
func (g *Grid) SetCellAt(row uint, col uint, newCell Cell) {
	if newCell != nil && (*newCell < 1 || *newCell > g.maxElement()) {
		// TODO - better error handling?
		panic("tried to set cell value outside the valid range")
	}

	// could inline these into a single line, but gofmt's arithmetic formatting makes the inlined version less clear
	rowBaseIndex := row * g.sideLength()
	index := rowBaseIndex + col
	g.cells[index] = newCell
}

func (g *Grid) rows() [][]Cell {
	rows := [][]Cell{}

	for r := uint(0); r < g.sideLength(); r++ {
		row := []Cell{}
		// rowBaseIndex := r * g.sideLength() // could inline this, but gofmt's arithmetic formatting makes the inlined version less clear
		for c := uint(0); c < g.sideLength(); c++ {
			// idx := rowBaseIndex + c // could inline this, but gofmt's arithmetic formatting makes the inlined version less clear
			// row[c] = g.cells[idx]
			// row[c] = g.cellAt(r, c)
			row = append(row, g.cellAt(r, c))
		}
		rows = append(rows, row)
	}

	return rows
}

func (g *Grid) cols() [][]Cell {
	cols := [][]Cell{}

	for c := uint(0); c < g.sideLength(); c++ {
		col := []Cell{}

		for r := uint(0); r < g.sideLength(); r++ {
			// idx := r*g.sideLength() + c
			// col[r] = g.cells[idx]
			// col[r] = g.cellAt(r, c)
			col = append(col, g.cellAt(r, c))
		}

		cols = append(cols, col)
	}

	return cols
}

func (g *Grid) boxes() [][]Cell {
	boxes := [][]Cell{}

	// boxRow and boxCol are the row/column of the boxes within the overall grid;
	// for example, for a 9x9 sudoku (baseSize == 3), the grid has boxes in a 3x3 arrangement
	// so the upper-left box has boxRow 0 and boxCol 0
	// the center-left box has boxRow 1 and boxCol 0
	// the middle box has boxRow 1 and boxCol 1
	// (and so on)

	// TODO - explain calculations
	// TODO - is there a way to calculate each cell's row and column independently, then use g.cellAt(), to make this simpler?

	for boxRow := uint(0); boxRow < g.baseSize; boxRow++ {
		for boxCol := uint(0); boxCol < g.baseSize; boxCol++ {
			box := []Cell{}

			boxBaseIndex := (boxRow * g.baseSize * g.sideLength()) + (boxCol * g.baseSize)

			// r and c are coordinates relative to the inside of the box
			for r := uint(0); r < g.baseSize; r++ {
				for c := uint(0); c < g.baseSize; c++ {
					index := boxBaseIndex + r*g.sideLength() + c
					box = append(box, g.cells[index])
				}
			}

			boxes = append(boxes, box)
		}
	}

	return boxes
}

func (g *Grid) IsCompletelyFilled() bool {
	for _, cell := range g.cells {
		if cell == nil {
			return false
		}
	}

	return true
}

// checks if the house contains all elements in the range from 1 to maxElement, inclusive
func isHouseValid(house []Cell, maxElement int) bool {
	if len(house) != maxElement {
		return false
	}

	for _, cell := range house {
		if cell == nil {
			return false
		}
	}

	// all cells in the house are non-nil

	elementsInHouse := lo.Map(house, func(cell Cell, _ int) int {
		return *cell
	})
	slices.Sort(elementsInHouse)

	allElements := lo.RangeFrom(1, maxElement)

	return slices.Equal(elementsInHouse, allElements)
}

// checks if each row, column, and box has exactly one of each digit/element
// only checks completely filled-out grids; if a grid has any empty cells, this returns false
// does *not* check if a grid matches a specific puzzle (whether it matches the givens from the puzzle)
func (g *Grid) IsValidSolution() bool {
	if !g.IsCompletelyFilled() {
		return false
	}

	// all cells have values (are non-nil)

	for _, row := range g.rows() {
		if !isHouseValid(row, g.maxElement()) {
			return false
		}
	}

	for _, col := range g.cols() {
		if !isHouseValid(col, g.maxElement()) {
			return false
		}
	}

	for _, box := range g.boxes() {
		if !isHouseValid(box, g.maxElement()) {
			return false
		}
	}

	return true
}
