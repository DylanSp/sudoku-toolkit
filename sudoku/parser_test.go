package sudoku

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadGridsFromFile(t *testing.T) {
	// TODO - figure out a less brittle way get the path to the example files? maybe copy them to a subfolder in this directory?
	currentWorkingDir, err := os.Getwd()
	assert.NoError(t, err)
	examplesFolder := filepath.Join(currentWorkingDir, "..", "examples")

	t.Run("Parsing 4x4 grids", func(t *testing.T) {
		exampleFiles := []string{
			"filledGrid.txt",
			"oneCellEmpty.txt",
			"simpleChallenge.txt",
		}
		for _, exampleFile := range exampleFiles {
			filename := filepath.Join(examplesFolder, "4x4", exampleFile)
			testLoading4x4GridsFromFile(t, filename)
		}
	})

	t.Run("Parsing 9x9 grids", func(t *testing.T) {
		exampleFiles := []string{
			"easy50.txt",
			"hard95.txt",
			"hardest.txt",
		}
		for _, exampleFile := range exampleFiles {
			filename := filepath.Join(examplesFolder, "9x9", exampleFile)
			testLoading9x9GridsFromFile(t, filename)
		}
	})
}

// tests that parsing succeeds and returns a Grid with the right size
func testLoading4x4GridsFromFile(t *testing.T, filename string) {
	t.Helper()

	grids, err := LoadGridsFromFile(filename)
	assert.NoError(t, err)

	for _, grid := range grids {
		assert.EqualValues(t, 2, grid.baseSize)
	}
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

// TODO - uses non-exported functions (readFileLines(), parseSingleGrid()), but testing this round-trip property is easy and useful
// TODO - use fuzz testing for this?
func TestRoundTripFromFile(t *testing.T) {
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
		lines, err := readFileLines(filename)
		assert.NoError(t, err)
		for _, line := range lines {
			grid := ParseSingleGrid(line)
			assert.EqualValues(t, line, grid.String())
		}
	}
}
