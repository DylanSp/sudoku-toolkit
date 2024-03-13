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

	utils.Assert(solution.underlyingGrid.IsValidSolution(), "Backtracking solver reached invalid solution")

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

	// justEnteredInnerLoop := false

	for {
		// justEnteredInnerLoop = true

		// apply basic rules and assignments as long as possible, until either the grid is completed or no progress can be made
		for {
			if puzzle.underlyingGrid.IsCompletelyFilled() && !puzzle.underlyingGrid.IsValidSolution() {
				// fmt.Printf("justEnteredInnerLoop: %v\n", justEnteredInnerLoop)

				// search value creates an invalid state - backtrack
				return puzzle, false

				// TODO - does this work? why?
				// when I was debugging this and had this check cause a panic, it triggered when the inner for-loop has been entered for the first time, so I *think* it's due to the search value causing an invalid solution
				// but the possibility chosen as the search value should have been eliminated before it was chosen as a search value
			}

			utils.Assert(!puzzle.underlyingGrid.IsCompletelyFilled() || puzzle.underlyingGrid.IsValidSolution(), "Grid is invalid at start of inner for-loop")

			anyValuesEliminated := puzzle.eliminatePossibilitiesByRules()

			utils.Assert(!puzzle.underlyingGrid.IsCompletelyFilled() || puzzle.underlyingGrid.IsValidSolution(), "Grid is invalid after eliminating possibilities")

			if puzzle.hasReachedContradiction() {
				fmt.Println("Puzzle has reached contradiction")
				// the puzzle can no longer be solved, some cell doesn't have any possibilities left
				// therefore, we made an incorrect choice when searching - return false so we can backtrack
				return puzzle, false
			}

			anyValuesAssigned := puzzle.assignValuesForSinglePossibilities()

			// utils.Assert(!puzzle.underlyingGrid.IsCompletelyFilled() || puzzle.underlyingGrid.IsValidSolution(), "Grid is invalid after assigning values for single possibilities")

			if puzzle.underlyingGrid.IsValidSolution() {
				return puzzle, true
			}

			// TODO - does this work? if so, why?
			if puzzle.underlyingGrid.IsCompletelyFilled() {
				// grid is completely filled, but is invalid
				// therefore, we made an incorrect choice when searching - return false so we can backtrack
				return puzzle, false
			}

			// no progress made and puzzle is still incomplete - break out of inner loop and start searching
			if !anyValuesEliminated && !anyValuesAssigned {
				break
			}

			// justEnteredInnerLoop = false
		}

		utils.Assert(!puzzle.underlyingGrid.IsCompletelyFilled() || puzzle.underlyingGrid.IsValidSolution(), "Grid is invalid when beginning search")

		// puzzle isn't complete, no progress can be made with simple strategies
		// all remaining empty cells have at least 2 possibilities

		// declared locally because this isn't useful outside this function
		// takes a puzzle as a parameter instead of closing over the existing puzzle to avoid any weirdness with captured values in recursive calls
		// TODO - move this into separate function?
		// findFirstSearchCandidate := func(puzzle Puzzle) int {
		// 	for i, cellPossibilities := range puzzle.possibleValues {
		// 		utils.Assertf(cellPossibilities.Size() > 0, "Cell %v has no possibilities remaining", i)

		// 		if cellPossibilities.Size() < 2 {
		// 			cellValue := puzzle.underlyingGrid.cells[i].value
		// 			utils.Assertf(cellValue != nil, "Cell %v with only 1 possibility should have a value assigned", i)

		// 			continue
		// 		}

		// 		return i
		// 	}

		// 	fmt.Println(puzzle.underlyingGrid.String())

		// 	cellPtrs := []string{}
		// 	for i, cellPtr := range puzzle.underlyingGrid.cells {
		// 		// fmt.Printf("Cell %v: %p\n", i, cellPtr)
		// 		str := fmt.Sprintf("%p", cellPtr)
		// 		cellPtrs = append(cellPtrs, str)

		// 		peers := cellPtr.AllPeers()
		// 		peersSlice := peers.Elements()

		// 		peerIndexes := lo.Map(peersSlice, func(peer *Cell, _ int) string {
		// 			peersStr := ""
		// 			if peer == nil {
		// 				fmt.Println("Nil peer?!?!")
		// 				// peersStr = append(peersStr, ".")
		// 				peersStr += "."
		// 			} else {
		// 				peersStr += strconv.Itoa(peer.index)
		// 			}
		// 			return peersStr
		// 		})

		// 		fmt.Printf("Cell %v:\n", i)
		// 		fmt.Printf("Value: %v\n", *cellPtr.value)
		// 		fmt.Printf("Possibilities: %v\n", puzzle.possibleValues[i].Elements())
		// 		fmt.Printf("Indexes of peers: %v\n", peerIndexes)
		// 		fmt.Println()

		// 		utils.Assert(len(peerIndexes) == 20, "all cells should have exactly 20 peers")
		// 	}
		// 	slices.Sort(cellPtrs)
		// 	uniqdCellPtrs := slices.Compact(cellPtrs)
		// 	fmt.Println(len(uniqdCellPtrs))

		// 	utils.Assert(len(uniqdCellPtrs) == 81, "Some cell pointers were non-unique")

		// 	utils.Assert(!puzzle.underlyingGrid.IsCompletelyFilled() || puzzle.underlyingGrid.IsValidSolution(), "Invalid solution")

		// 	panic("couldn't find a cell with at least 2 possibilities, even though there should be one")
		// }

		gridBeforeChoosingSearchCandidate := puzzle.underlyingGrid.cells

		// find the first (empty) cell with at least 2 possibilities,
		// pick one possibility and set it, remembering other possibilities in case search fails,
		// then recursively search
		// TODO - use a heuristic to find a good search candidate?
		// Norvig mentions searching from a cell with the smallest set of remaining values
		// searchCandidateIndex := findFirstSearchCandidate(puzzle)
		searchCandidateIndex := puzzle.findFirstSearchCandidate()
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

		fmt.Println("Recursing")
		fmt.Printf("Searching with value %v in cell %v\n", valueChosenForSearch, searchCandidateIndex)
		possibleSolution, ok := attemptBacktrackingSolve(puzzleWithSearchBranch)
		if ok {
			fmt.Println("Returning")
			// valid solution found, return it
			return possibleSolution, true
		}
		fmt.Println("Backtracking")

		// search based on the chosen value didn't find a solution; eliminate it, then return to start of loop
		possibilitiesForSearchCell.DeleteAll()
		for _, possibility := range remainingPossibilities {
			possibilitiesForSearchCell.Add(possibility)
		}
		puzzle.underlyingGrid.cells[searchCandidateIndex].value = nil

		gridAfterBacktracking := puzzle.underlyingGrid.cells

		for i, cellFromBefore := range gridBeforeChoosingSearchCandidate {
			cellFromAfter := gridAfterBacktracking[i]

			if cellFromBefore == nil || cellFromAfter == nil {
				panic("something is funky, cell is nil")
			}

			if cellFromBefore.value == nil && cellFromAfter.value == nil {
				continue
			}

			if cellFromBefore.value == nil && cellFromAfter.value != nil {
				fmt.Printf("Cell %v retains value from search\n", i)
				panic("backtracking failure")
			}

			if cellFromBefore.value != nil && cellFromAfter.value == nil {
				fmt.Printf("Value of cell %v was erased during search\n", i)
				panic("backtracking failure")
			}

			if *(cellFromBefore.value) != *(cellFromAfter.value) {
				fmt.Printf("Value of cell %v was changed from %v to %v during search\n", i, *cellFromBefore.value, *cellFromAfter.value)
				panic("backtracking failure")
			}
		}

		utils.Assert(!puzzle.underlyingGrid.IsCompletelyFilled() || puzzle.underlyingGrid.IsValidSolution(), "Grid is invalid after resetting from backtrack")
	}
}

func (puzzle *Puzzle) findFirstSearchCandidate() int {
	for i, cellPossibilities := range puzzle.possibleValues {
		utils.Assertf(cellPossibilities.Size() > 0, "Cell %v has no possibilities remaining", i)

		if cellPossibilities.Size() < 2 {
			cellValue := puzzle.underlyingGrid.cells[i].value
			utils.Assertf(cellValue != nil, "Cell %v with only 1 possibility should have a value assigned", i)

			continue
		}

		return i
	}

	// Everything below this comment in this method is debugging for being unable to find a search candidate

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
					if peerValue != nil && *peerValue == possibility {
						// fmt.Printf("Current cell is %v, possibility is %v\n", i, possibility)
						// fmt.Printf("Peer %v has the same value %v\n", peer.index, *peerValue)
						// panic("same value in peer")
					}

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
	deletionIndexes := []int{}
	deletedValues := []int{}

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

			utils.Assertf(len(peers) == 20, "Cell %v does not have exactly 20 peers", i)

			// TODO - nested loop here - possible source of inefficiency?
			for _, peer := range peers {
				if peer != nil && !peer.isEmpty() {
					deletionMade := possibilitiesForCell.Delete(*peer.value)
					if deletionMade {
						eliminationsMadeInLoop = true
						eliminationsMadeInMethod = true

						// debugging
						// deletionIndex = i
						// deletedValue = *peer.value

						deletionIndexes = append(deletionIndexes, i)
						deletedValues = append(deletedValues, *peer.value)

						// fmt.Println("Inside loop")
						// fmt.Printf("Deleted %v from cell %v\n", deletedValue, deletionIndex)
						// fmt.Printf("Remaining possibilities: %v\n", puzzle.possibleValues[deletionIndex].Elements())
						// fmt.Printf("possibilitiesForCell.Has(deletedValue): %v\n", possibilitiesForCell.Has(deletedValue))
						// fmt.Printf("puzzle.possibleValues[deletionIndex].Has(deletedValue): %v\n", puzzle.possibleValues[deletionIndex].Has(deletedValue))
						// fmt.Printf("Argument to Assertf: %v\n", !puzzle.possibleValues[deletionIndex].Has(deletedValue))
						// utils.Assertf(!puzzle.possibleValues[deletionIndex].Has(deletedValue), "Possibilities for cell %v still include %v", deletionIndex, deletedValue)
					}
				}
			}

			// debugging

			// if eliminationsMadeInLoop {
			// 	fmt.Println("Outside loop")
			// 	fmt.Printf("possibilitiesForCell.Has(deletedValue): %v\n", possibilitiesForCell.Has(deletedValue))
			// 	fmt.Printf("puzzle.possibleValues[deletionIndex].Has(deletedValue): %v\n", puzzle.possibleValues[deletionIndex].Has(deletedValue))
			// 	utils.Assertf(!puzzle.possibleValues[i].Has(deletedValue), "Possibilities for cell %v still contains %v\n", deletionIndex, deletedValue)
			// }

			// for i, deletionIndex := range deletionIndexes {
			// 	deletedValue := deletedValues[i]
			// 	fmt.Printf("possibilitiesForCell.Has(deletedValue): %v\n", possibilitiesForCell.Has(deletedValue))
			// 	fmt.Printf("puzzle.possibleValues[deletionIndex].Has(deletedValue): %v\n", puzzle.possibleValues[deletionIndex].Has(deletedValue))
			// 	utils.Assertf(!puzzle.possibleValues[deletionIndex].Has(deletedValue), "Possibilities for cell %v still contains %v\n", deletionIndex, deletedValue)
			// }

		}

		if !eliminationsMadeInLoop {
			break
		}
		eliminationsMadeInLoop = false // reset for next iteration
		// deletedValue = -1
		// deletionIndex = -1
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
