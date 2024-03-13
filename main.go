package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/DylanSp/sudoku-toolkit/sudoku"
)

func main() {
	// apparently valid
	// solution := "462831957795426183381795426173984265659312748248567319926178534834259671517643892"

	// solution := "462813957795426183381795426173984265659231748248567319926178534834659271517342698"

	// parsedSolution := sudoku.ParseSingleGrid(solution)
	// fmt.Println(parsedSolution.IsValidSolution())
	// return

	currentWorkingDir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	// filename := filepath.Join(currentWorkingDir, "examples", "4x4", "filledGrid.txt")
	// filename := filepath.Join(currentWorkingDir, "examples", "4x4", "oneCellEmpty.txt")
	// filename := filepath.Join(currentWorkingDir, "examples", "4x4", "simpleChallenge.txt")

	filename := filepath.Join(currentWorkingDir, "examples", "9x9", "easy50.txt")
	// filename := filepath.Join(currentWorkingDir, "examples", "9x9", "hard95.txt")
	// filename := filepath.Join(currentWorkingDir, "examples", "9x9", "hardest.txt")

	grids, err := sudoku.LoadGridsFromFile(filename)
	if err != nil {
		fmt.Printf("Unable to load grids from %v\n", filename)
		fmt.Println(err)
		panic(err)
	}

	solvedGrid := sudoku.SolveWithBacktracking(grids[2])
	fmt.Println(solvedGrid.String())

	// for _, grid := range grids {
	// 	// solvedGrid := sudoku.SolveWithBasicStrategies(grid)
	// 	solvedGrid := sudoku.SolveWithBacktracking(grid)
	// 	fmt.Println(solvedGrid.String())
	// }
}
