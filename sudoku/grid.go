package sudoku

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/DylanSp/sudoku-toolkit/utils"
	"github.com/samber/lo"
)

type Cell struct {
	index          int   // index of this cell in containingGrid.cells
	containingGrid *Grid // pointer to act as a reference, avoiding infinite loops; should never be nil
	value          *int  // nil == no value
}

// utility method for clarity
func (c *Cell) isEmpty() bool {
	return c.value == nil
}

func (c *Cell) String() string {
	if c.isEmpty() {
		return "."
	}

	// TODO - generalize for different base sizes/ranges of possible values
	if *c.value >= 1 && *c.value <= 9 {
		return strconv.Itoa(*c.value)
	}

	panic(fmt.Sprintf("don't know to print cell value %v", *c.value))
}

// TODO - potential performance optimization - precalculate this for all cells when grid is created
func (c *Cell) AllPeers() utils.Set[*Cell] {
	peers := utils.Set[*Cell]{}

	indexes := c.allPeerIndexes()
	for _, peerIndex := range indexes.Elements() {
		peers.Add(c.containingGrid.cells[peerIndex])
	}

	return peers
}

func (c *Cell) allPeerIndexes() utils.Set[int] {
	allPeerIndexes := utils.Set[int]{}

	// we can't inline these into "range c.rowPeerIndexes().Elements()" because .Elements() has a pointer receiver,
	// so we need an intermediate variable for each group of indexes

	// TODO - potentially combine the loops in c.rowPeerIndexes() and c.colPeerIndexes()?
	// we can use a divmod-like function: div, mod := i / sideLength, i % sideLength

	rowPeerIndexes := c.rowPeerIndexes()
	for _, peerIndex := range rowPeerIndexes.Elements() {
		allPeerIndexes.Add(peerIndex)
	}

	colPeerIndexes := c.colPeerIndexes()
	for _, peerIndex := range colPeerIndexes.Elements() {
		allPeerIndexes.Add(peerIndex)
	}

	boxPeerIndexes := c.boxPeerIndexes()
	for _, peerIndex := range boxPeerIndexes.Elements() {
		allPeerIndexes.Add(peerIndex)
	}

	return allPeerIndexes
}

// indexes of all cells that are in the same row as c
func (c *Cell) rowPeerIndexes() utils.Set[int] {
	peerIndexes := utils.Set[int]{}

	for i := range c.containingGrid.cells {
		// c isn't its own peer
		if i == c.index {
			continue
		}

		// integer division, rounding down
		if i/c.containingGrid.sideLength() == c.index/c.containingGrid.sideLength() {
			peerIndexes.Add(i)
		}
	}

	return peerIndexes
}

// indexes of all cells that are in the same column as c
func (c *Cell) colPeerIndexes() utils.Set[int] {
	peerIndexes := utils.Set[int]{}

	for i := range c.containingGrid.cells {
		// c isn't its own peer
		if i == c.index {
			continue
		}

		if i%c.containingGrid.sideLength() == c.index%c.containingGrid.sideLength() {
			peerIndexes.Add(i)
		}
	}

	return peerIndexes
}

// indexes of all cells that are in the same box as c
func (c *Cell) boxPeerIndexes() utils.Set[int] {
	// unsure if there's an easy arithmetic way to find this
	// instead, loop through all boxes in the grid, find which one contains c, then loop through its contents

	peerIndexes := utils.Set[int]{}

	var boxContainingC []*Cell

	// should always find and assign a value for boxContainingC
	for _, box := range c.containingGrid.boxes() {
		if slices.ContainsFunc(box, func(cell *Cell) bool {
			return cell.index == c.index
		}) {
			boxContainingC = box
			break
		}
	}

	for _, cellInBox := range boxContainingC {
		// c isn't its own peer
		if cellInBox.index == c.index {
			continue
		}

		peerIndexes.Add(cellInBox.index)
	}

	return peerIndexes
}

type Grid struct {
	baseSize int

	// flat slice of cells, row-by-row; may change representation later
	// invariant: len(cells) == baseSize ^ 4
	cells []*Cell // references to cells so they can be mutated; none of these should ever be nil
}

func EmptyGrid(baseSize int) Grid {
	grid := Grid{
		baseSize: baseSize,
	}

	// example: a 9x9 Sudoku has baseSize 3: each side of the grid is 3*3 = 9 cells
	// the total grid is 9*9 = 81 = 3^4 cells
	// use repeated multiplication instead of math.Pow() to avoid converting to/from float64
	grid.cells = make([]*Cell, baseSize*baseSize*baseSize*baseSize)

	return grid
}

func (g *Grid) sideLength() int {
	return g.baseSize * g.baseSize
}

// each house in the grid must have all elements in the range from 1 to maxElement(), inclusive
// this is the same calculation as sideLength(), but split into a separate method for clarity
func (g *Grid) maxElement() int {
	return g.baseSize * g.baseSize
}

func (g *Grid) cellAt(row int, col int) *Cell {
	// could inline these into a single line, but gofmt's arithmetic formatting makes the inlined version less clear
	rowBaseIndex := row * g.sideLength()
	index := rowBaseIndex + col
	return g.cells[index]
}

func (g *Grid) rows() [][]*Cell {
	rows := [][]*Cell{}

	for r := 0; r < g.sideLength(); r++ {
		row := []*Cell{}
		// rowBaseIndex := r * g.sideLength() // could inline this, but gofmt's arithmetic formatting makes the inlined version less clear
		for c := 0; c < g.sideLength(); c++ {
			// idx := rowBaseIndex + c // could inline this, but gofmt's arithmetic formatting makes the inlined version less clear
			// row[c] = g.cells[idx]
			// row[c] = g.cellAt(r, c)
			row = append(row, g.cellAt(r, c))
		}
		rows = append(rows, row)
	}

	return rows
}

func (g *Grid) cols() [][]*Cell {
	cols := [][]*Cell{}

	for c := 0; c < g.sideLength(); c++ {
		col := []*Cell{}

		for r := 0; r < g.sideLength(); r++ {
			// idx := r*g.sideLength() + c
			// col[r] = g.cells[idx]
			// col[r] = g.cellAt(r, c)
			col = append(col, g.cellAt(r, c))
		}

		cols = append(cols, col)
	}

	return cols
}

func (g *Grid) boxes() [][]*Cell {
	boxes := [][]*Cell{}

	// boxRow and boxCol are the row/column of the boxes within the overall grid;
	// for example, for a 9x9 sudoku (baseSize == 3), the grid has boxes in a 3x3 arrangement
	// so the upper-left box has boxRow 0 and boxCol 0
	// the center-left box has boxRow 1 and boxCol 0
	// the middle box has boxRow 1 and boxCol 1
	// (and so on)

	// TODO - explain calculations
	// TODO - is there a way to calculate each cell's row and column independently, then use g.cellAt(), to make this simpler?

	for boxRow := 0; boxRow < g.baseSize; boxRow++ {
		for boxCol := 0; boxCol < g.baseSize; boxCol++ {
			box := []*Cell{}

			boxBaseIndex := (boxRow * g.baseSize * g.sideLength()) + (boxCol * g.baseSize)

			// r and c are coordinates relative to the inside of the box
			for r := 0; r < g.baseSize; r++ {
				for c := 0; c < g.baseSize; c++ {
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
		if cell.isEmpty() {
			return false
		}
	}

	return true
}

// checks if the house contains all elements in the range from 1 to maxElement, inclusive
func isHouseValid(house []*Cell, maxElement int) bool {
	if len(house) != maxElement {
		return false
	}

	for _, cell := range house {
		if cell.isEmpty() {
			return false
		}
	}

	// all cells in the house are non-nil

	elementsInHouse := lo.Map(house, func(cell *Cell, _ int) int {
		return *cell.value
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

func (g *Grid) String() string {
	var b strings.Builder

	for _, cell := range g.cells {
		b.WriteString(cell.String())
	}

	return b.String()
}

// func (g *Grid) allPeersOfCell()
