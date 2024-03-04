package sudoku

import (
	"github.com/DylanSp/sudoku-toolkit/utils"
)

type puzzleInProgress struct {
	underlyingGrid Grid

	// possibleValues[i]: possible values for the cell at knownValues.cells[i]
	// invariant: len(possibleValues) == len(knownValues.cells)
	// TODO - uses int to represent values - should it use a different representation?
	possibleValues []utils.Set[int]
}

func SolveWithBasicStrategies(grid Grid) Grid {
	if grid.IsValidSolution() {
		return grid
	}

	puzzle := newPuzzleInProgress(grid)

	// apply basic rules and assignments as long as possible, until either the grid is completed or no progress can be made
	for {
		anyValuesEliminated := puzzle.eliminatePossibilitiesByRules()
		anyValuesAssigned := puzzle.assignValuesForSinglePossibilities()

		if puzzle.underlyingGrid.IsValidSolution() {
			return puzzle.underlyingGrid
		}

		// no progress made and puzzle is still incomplete
		if !anyValuesEliminated && !anyValuesAssigned {
			// TODO - more graceful error handling
			panic("Unable to solve puzzle with basic strategies")
		}
	}
}

func SolveWithBacktracking(grid Grid) Grid {
	// if grid.IsValidSolution() {
	// 	return grid
	// }

	// puzzle := newPuzzleInProgress(grid)

	panic("not yet implemented")
}

func newPuzzleInProgress(grid Grid) puzzleInProgress {
	puzzle := puzzleInProgress{
		underlyingGrid: grid,
		possibleValues: make([]utils.Set[int], len(grid.cells)),
	}

	for i, cell := range grid.cells {
		if cell.isEmpty() {
			puzzle.possibleValues[i] = allPossibilities(grid.baseSize)
		} else {
			puzzle.possibleValues[i] = utils.Set[int]{}
			puzzle.possibleValues[i].Add(*cell.value)
		}
	}

	return puzzle
}

// returns a set with all possible elements for a grid with the given base size
func allPossibilities(baseSize int) utils.Set[int] {
	possibilities := utils.Set[int]{}

	maxElement := baseSize * baseSize
	for i := 1; i <= maxElement; i++ {
		possibilities.Add(i)
	}

	return possibilities
}

// go through all empty cells; if there's only one possibile value, set that cell's value to that possibility
// returns true iff at least one value was assigned
func (puzzle *puzzleInProgress) assignValuesForSinglePossibilities() bool {
	valueAssigned := false

	for i, cell := range puzzle.underlyingGrid.cells {
		if cell.isEmpty() {
			possibilitiesForCell := puzzle.possibleValues[i]
			if possibilitiesForCell.Size() == 1 {
				possibility := possibilitiesForCell.Elements()[0]
				cell.value = &possibility
				valueAssigned = true
			}
		}
	}

	return valueAssigned
}

// applies the basic rules of Sudoku to eliminate all possibilities ruled out by currently known values
// returns true iff at least one possibility was eliminated
func (puzzle *puzzleInProgress) eliminatePossibilitiesByRules() bool {
	eliminationsMadeInMethod := false // did this method as a whole eliminate any possibilities?

	eliminationsMadeInLoop := false // did a specific iteration of the loop eliminate any possibilities?

	// continue looping until we can no longer eliminate any possibilities
	for {
		for i, cell := range puzzle.underlyingGrid.cells {
			// skip cells that already have values
			if !cell.isEmpty() {
				continue
			}

			possibilitiesForCell := &puzzle.possibleValues[i]

			peerSet := cell.AllPeers()
			peers := peerSet.Elements()

			// TODO - nested loop here - possible source of inefficiency?
			for _, peer := range peers {
				if peer != nil && !peer.isEmpty() {
					possibilitiesForCell.Delete(*peer.value)
					eliminationsMadeInLoop = true
					eliminationsMadeInMethod = true
				}
			}
		}

		if !eliminationsMadeInLoop {
			break
		}
		eliminationsMadeInLoop = false // reset for next iteration
	}

	return eliminationsMadeInMethod

}
