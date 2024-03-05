package sudoku_test

import (
	"testing"

	"github.com/DylanSp/sudoku-toolkit/sudoku"
	"github.com/stretchr/testify/assert"
)

func TestSolveWithBasicStrategies(t *testing.T) {
	t.Run("Solving a 4x4 puzzle with only one cell initially empty", func(t *testing.T) {
		puzzle := "143232144123234."
		expectedSolution := "1432321441232341"

		initialGrid := sudoku.ParseSingleGrid(puzzle)
		computedSolution := sudoku.SolveWithBasicStrategies(initialGrid)

		assert.EqualValues(t, expectedSolution, computedSolution.String())
	})
}
