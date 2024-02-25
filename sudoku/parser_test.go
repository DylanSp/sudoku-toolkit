package sudoku

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadGridsFromFile(t *testing.T) {
	t.Run("Parsing 9x9 grids", func(t *testing.T) {
		// TODO - figure out a less brittle way get the path to the example files? maybe copy them to a subfolder in this directory?
		currentWorkingDir, err := os.Getwd()
		assert.NoError(t, err)

		folderWithExamples := filepath.Join(currentWorkingDir, "..", "examples", "9x9")
		exampleFiles := []string{
			"easy50.txt",
			"hard95.txt",
			"hardest.txt",
		}
		for _, exampleFile := range exampleFiles {
			filename := filepath.Join(folderWithExamples, exampleFile)
			testLoading9x9GridsFromFile(t, filename)
		}
	})
}

// tests that parsing succeeds and returns a Grid with the right size
func testLoading9x9GridsFromFile(t *testing.T, filename string) {
	t.Helper()

	grids, err := LoadGridsFromFile(filename)
	assert.NoError(t, err)

	for _, grid := range grids {
		assert.EqualValues(t, 3, grid.baseSize)
	}
}
