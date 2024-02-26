package sudoku

import (
	"github.com/DylanSp/sudoku-toolkit/utils"
)

type puzzleInProgress struct {
	knownValues Grid

	// possibleValues[i]: possible values for the cell at knownValues.cells[i]
	// invariant: len(possibleValues) == len(knownValues.cells)
	// TODO - uses int to represent values - should it use a different representation?
	possibleValues []utils.Set[int]
}

func SolveWithBacktracking(grid Grid) Grid {
	if grid.IsValidSolution() {
		return grid
	}

	puzzle := createPuzzle(grid)

	panic("not yet implemented")

	// satisfy compiler - will probably not be used
	return puzzle.knownValues
}

func createPuzzle(grid Grid) puzzleInProgress {
	puzzle := puzzleInProgress{
		knownValues:    grid,
		possibleValues: make([]utils.Set[int], len(grid.cells)),
	}

	for i, cell := range grid.cells {
		if cell == nil {
			puzzle.possibleValues[i] = allPossibilities(grid.baseSize)
		} else {
			puzzle.possibleValues[i] = utils.Set[int]{}
			puzzle.possibleValues[i].Add(*cell)
		}
	}

	return puzzle
}

// returns a set with all possible elements for a grid with the given base size
func allPossibilities(baseSize uint) utils.Set[int] {
	possibilities := utils.Set[int]{}

	maxElement := baseSize * baseSize
	for i := 1; i <= int(maxElement); i++ {
		possibilities.Add(i)
	}

	return possibilities
}
