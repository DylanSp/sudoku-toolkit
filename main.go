package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/DylanSp/sudoku-toolkit/sudoku"
)

func main() {
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

	solvedGrid := sudoku.SolveWithBacktracking(grids[1])
	fmt.Println(solvedGrid.String())

	// for _, grid := range grids {
	// 	// solvedGrid := sudoku.SolveWithBasicStrategies(grid)
	// 	solvedGrid := sudoku.SolveWithBacktracking(grid)
	// 	fmt.Println(solvedGrid.String())
	// }
}
