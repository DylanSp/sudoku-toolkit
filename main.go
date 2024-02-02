package main

import (
	"fmt"

	"github.com/DylanSp/sudoku-toolkit/sudoku"
)

func main() {
	grid := sudoku.EmptyGrid(2)

	fmt.Printf("Is an empty grid completely filled? (should be false): %v\n", grid.IsCompletelyFilled())

	// valid 4x4 sudoku solution:
	/*
			2 3 | 1 4
			4 1 | 3 2
			----|----
			3 4 | 2 1
		    1 2 | 4 3
	*/

	grid.SetCellToElement(0, 0, 2)
	grid.SetCellToElement(0, 1, 3)
	grid.SetCellToElement(0, 2, 1)
	grid.SetCellToElement(0, 3, 4)

	grid.SetCellToElement(1, 0, 4)
	grid.SetCellToElement(1, 1, 1)
	grid.SetCellToElement(1, 2, 3)
	grid.SetCellToElement(1, 3, 2)

	grid.SetCellToElement(2, 0, 3)
	grid.SetCellToElement(2, 1, 4)
	grid.SetCellToElement(2, 2, 2)
	grid.SetCellToElement(2, 3, 1)

	grid.SetCellToElement(3, 0, 1)
	grid.SetCellToElement(3, 1, 2)
	grid.SetCellToElement(3, 2, 4)
	grid.SetCellToElement(3, 3, 3)

	fmt.Printf("Is the grid completely filled after setting values? (should be true): %v\n", grid.IsCompletelyFilled())
	fmt.Printf("Is the solution valid? (should be true): %v\n", grid.IsValidSolution())

	grid.SetCellToElement(3, 3, 4)
	fmt.Printf("Is the solution valid after changing a value? (should be false): %v\n", grid.IsValidSolution())
}
