package sudoku

import (
	"github.com/DylanSp/sudoku-toolkit/utils"
)

// TODO - not sure if I want to export the Puzzle type
// It's currently exported because gopls won't allow renaming it to "puzzle" due to potentially shadowing parameter names
// this might be a gopls bug - see https://github.com/golang/go/issues/66150
type Puzzle struct {
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

	puzzle := newPuzzle(grid)

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

// public entrypoint, validating input and wrapping attemptBacktrackingSolve()
func SolveWithBacktracking(grid Grid) Grid {
	if grid.IsValidSolution() {
		return grid
	}

	puzzle := newPuzzle(grid)

	solution, ok := attemptBacktrackingSolve(puzzle)
	if !ok {
		// TODO - more graceful error handling
		panic("Unable to solve puzzle with backtracking")
	}

	return solution.underlyingGrid
}

// attempts to solve a puzzle with recursive, backtracking search
// if it finds a valid solution, will return (solution, true)
// if it can't find a valid solution, will return (<unfinished puzzle>, false)
// TODO - better name?
func attemptBacktrackingSolve(puzzle Puzzle) (Puzzle, bool) {

	// TODO - we need a way to detect when we've encountered an invalid search branch (no candidates remain for at least one empty cell)
	// then return false in that case, so solver can backtract
	// not sure where that should go in this control flow

	for {

		// apply basic rules and assignments as long as possible, until either the grid is completed or no progress can be made
		for {
			anyValuesEliminated := puzzle.eliminatePossibilitiesByRules()

			if puzzle.hasReachedContradiction() {
				// the puzzle can no longer be solved, some cell doesn't have any possibilities left
				// therefore, we made an incorrect choice when searching - return false so we can backtrack
				return puzzle, false
			}

			anyValuesAssigned := puzzle.assignValuesForSinglePossibilities()

			if puzzle.underlyingGrid.IsValidSolution() {
				return puzzle, true
			}

			// no progress made and puzzle is still incomplete - break out of inner loop and start searching
			if !anyValuesEliminated && !anyValuesAssigned {
				break
			}
		}

		// puzzle isn't complete, no progress can be made with simple strategies
		// all remaining empty cells have at least 2 possibilities

		// declared locally because this isn't useful outside this function
		// takes a puzzle as a parameter instead of closing over the existing puzzle to avoid any weirdness with captured values in recursive calls
		// TODO - move this into separate function?
		findFirstSearchCandidate := func(puzzle Puzzle) int {
			for i, cellPossibilities := range puzzle.possibleValues {
				if cellPossibilities.Size() < 2 {
					continue
				}

				return i
			}

			panic("couldn't find a cell with at least 2 possibilities, even though there should be one!")
		}

		// find the first (empty) cell with at least 2 possibilities,
		// pick one possibility and set it, remembering other possibilities in case search fails,
		// then recursively search
		// TODO - use a heuristic to find a good search candidate?
		// Norvig mentions searching from a cell with the smallest set of remaining values
		searchCandidateIndex := findFirstSearchCandidate(puzzle)
		possibilitiesForSearchCell := &puzzle.possibleValues[searchCandidateIndex]

		var valueChosenForSearch int
		remainingPossibilities := []int{}
		for i, possibility := range possibilitiesForSearchCell.Elements() {
			if i == 0 {
				valueChosenForSearch = possibility
			} else {
				remainingPossibilities = append(remainingPossibilities, possibility)
			}
		}

		possibilitiesForSearchCell.DeleteAll()
		possibilitiesForSearchCell.Add(valueChosenForSearch)
		puzzle.underlyingGrid.cells[searchCandidateIndex].value = &valueChosenForSearch

		// TODO - do we need to do some sort of deep clone on `puzzle` before calling this,
		// so if we need to backtrack, the original `puzzle` is still in the state it was before searching?
		possibleSolution, ok := attemptBacktrackingSolve(puzzle)
		if ok {
			// valid solution found, return it
			return possibleSolution, true
		}

		// search based on the chosen value didn't find a solution; eliminate it, then return to start of loop
		possibilitiesForSearchCell.DeleteAll()
		for _, possibility := range remainingPossibilities {
			possibilitiesForSearchCell.Add(possibility)
		}
		puzzle.underlyingGrid.cells[searchCandidateIndex].value = nil
	}
}

// detects whether a search has reached a contradiction by making an incorrect assumption - at least one Cell doesn't have any possible valid values
func (puzzle *Puzzle) hasReachedContradiction() bool {
	for _, cellPossibilities := range puzzle.possibleValues {
		if cellPossibilities.Size() == 0 {
			return true
		}
	}

	return false
}

func newPuzzle(grid Grid) Puzzle {
	puzzle := Puzzle{
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

// go through all empty cells; if there's only one possible value, set that cell's value to that possibility
// returns true iff at least one value was assigned
func (puzzle *Puzzle) assignValuesForSinglePossibilities() bool {
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
func (puzzle *Puzzle) eliminatePossibilitiesByRules() bool {
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
					deletionMade := possibilitiesForCell.Delete(*peer.value)
					if deletionMade {
						eliminationsMadeInLoop = true
						eliminationsMadeInMethod = true
					}
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
