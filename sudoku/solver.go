package sudoku

import (
	"fmt"
	"slices"
	"strconv"

	"github.com/DylanSp/sudoku-toolkit/utils"
	"github.com/samber/lo"
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
				utils.Assertf(cellPossibilities.Size() > 0, "Cell %v has no possibilities remaining", i)

				if cellPossibilities.Size() < 2 {
					cellValue := puzzle.underlyingGrid.cells[i].value
					utils.Assertf(cellValue != nil, "Cell %v with only 1 possibility should have a value assigned", i)

					continue
				}

				return i
			}

			fmt.Println(puzzle.underlyingGrid.String())

			cellPtrs := []string{}
			for i, cellPtr := range puzzle.underlyingGrid.cells {
				// fmt.Printf("Cell %v: %p\n", i, cellPtr)
				str := fmt.Sprintf("%p", cellPtr)
				cellPtrs = append(cellPtrs, str)

				peers := cellPtr.AllPeers()
				peersSlice := peers.Elements()

				peerIndexes := lo.Map(peersSlice, func(peer *Cell, _ int) string {
					peersStr := ""
					if peer == nil {
						fmt.Println("Nil peer?!?!")
						// peersStr = append(peersStr, ".")
						peersStr += "."
					} else {
						peersStr += strconv.Itoa(peer.index)
					}
					return peersStr
				})

				fmt.Printf("Cell %v:\n", i)
				fmt.Printf("Value: %v\n", *cellPtr.value)
				fmt.Printf("Possibilities: %v\n", puzzle.possibleValues[i].Elements())
				fmt.Printf("Indexes of peers: %v\n", peerIndexes)
				fmt.Println()

				utils.Assert(len(peerIndexes) == 20, "all cells should have exactly 20 peers")
			}
			slices.Sort(cellPtrs)
			uniqdCellPtrs := slices.Compact(cellPtrs)
			fmt.Println(len(uniqdCellPtrs))

			utils.Assert(len(uniqdCellPtrs) == 81, "Some cell pointers were non-unique")

			utils.Assert(!puzzle.underlyingGrid.IsCompletelyFilled() || puzzle.underlyingGrid.IsValidSolution(), "Invalid solution")

			panic("couldn't find a cell with at least 2 possibilities, even though there should be one")
		}

		// find the first (empty) cell with at least 2 possibilities,
		// pick one possibility and set it, remembering other possibilities in case search fails,
		// then recursively search
		// TODO - use a heuristic to find a good search candidate?
		// Norvig mentions searching from a cell with the smallest set of remaining values
		searchCandidateIndex := findFirstSearchCandidate(puzzle)
		possibilitiesForSearchCell := &puzzle.possibleValues[searchCandidateIndex] // take a reference so we can mutate this set

		var valueChosenForSearch int
		remainingPossibilities := []int{}
		for i, possibility := range possibilitiesForSearchCell.Elements() {
			if i == 0 {
				valueChosenForSearch = possibility
			} else {
				remainingPossibilities = append(remainingPossibilities, possibility)
			}
		}

		// put a deep copy of Puzzle here and make changes on it
		// otherwise, values assigned in a search branch will stay assigned even after backtracking, putting the puzzle in an inconsistent state
		puzzleWithSearchBranch := puzzle.deepClone()
		puzzleWithSearchBranch.possibleValues[searchCandidateIndex].DeleteAll()
		puzzleWithSearchBranch.possibleValues[searchCandidateIndex].Add(valueChosenForSearch)
		puzzleWithSearchBranch.underlyingGrid.cells[searchCandidateIndex].value = &valueChosenForSearch

		possibleSolution, ok := attemptBacktrackingSolve(puzzleWithSearchBranch)
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
				// if there's a single possibility, set that cell's value
				// from the if statement's condition, possibilitiesForCell already is restricted to just that value
				possibility := possibilitiesForCell.Elements()[0]

				peers := cell.AllPeers()
				for _, peer := range peers.Elements() {
					peerValue := peer.value
					utils.Assertf(peerValue == nil || *peerValue == possibility, "Peer %v has the same value", peer.index)
				}

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

	// debugging
	// deletionIndex := -1
	// deletedValue := -1

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

						// debugging
						deletionIndex := i
						deletedValue := *peer.value
						fmt.Printf("Deleted %v from cell %v\n", deletedValue, deletionIndex)
						fmt.Printf("Remaining possibilities: %v\n", puzzle.possibleValues[deletionIndex].Elements())
						utils.Assertf(!puzzle.possibleValues[deletionIndex].Has(deletedValue), "Possibilities for cell %v still include %v", deletionIndex, deletedValue)
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

func (puzzle *Puzzle) deepClone() Puzzle {
	newPuzzle := Puzzle{
		underlyingGrid: puzzle.underlyingGrid.DeepClone(),
	}
	// fmt.Printf("Original puzzle: %p\n", &puzzle.underlyingGrid)
	// fmt.Printf("New puzzle: %p\n", &newPuzzle.underlyingGrid)
	utils.Assert(fmt.Sprintf("%p", &puzzle.underlyingGrid) != fmt.Sprintf("%p", &newPuzzle.underlyingGrid), "puzzle clone's grid has the same memory address")

	newPossibleValues := []utils.Set[int]{}
	for _, cellPossibilities := range puzzle.possibleValues {
		// hypothesis - cloning the Set isn't working properly?
		// doesn't seem to be the case
		newCellPossibilities := cellPossibilities.Clone()

		// newCellPossibilities := utils.Set[int]{}
		// for _, possibility := range cellPossibilities.Elements() {
		// 	newCellPossibilities.Add(possibility)
		// }

		newPossibleValues = append(newPossibleValues, newCellPossibilities)
	}

	newPuzzle.possibleValues = newPossibleValues

	return newPuzzle
}
