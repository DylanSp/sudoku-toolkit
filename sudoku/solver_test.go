package sudoku_test

import (
	"testing"

	"github.com/DylanSp/sudoku-toolkit/sudoku"
	"github.com/stretchr/testify/assert"
)

func TestSolveWithBasicStrategies(t *testing.T) {
	t.Run("Solving a 4x4 challenge with only one cell initially empty", func(t *testing.T) {
		challenge := "143232144123234."
		expectedSolution := "1432321441232341"

		initialGrid := sudoku.ParseSingleGrid(challenge)
		computedSolution := sudoku.SolveWithBasicStrategies(initialGrid)

		assert.EqualValues(t, expectedSolution, computedSolution.String())
	})

	t.Run("Solving a 4x4 challenge with only 4 givens to start with", func(t *testing.T) {
		challenge := "1......4..2..3.."
		expectedSolution := "1432321441232341"

		initialGrid := sudoku.ParseSingleGrid(challenge)
		computedSolution := sudoku.SolveWithBasicStrategies(initialGrid)

		assert.EqualValues(t, expectedSolution, computedSolution.String())
	})
}
