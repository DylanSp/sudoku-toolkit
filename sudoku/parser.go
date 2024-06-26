package sudoku

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func LoadGridsFromFile(filename string) ([]Grid, error) {
	lines, err := readFileLines(filename)
	if err != nil {
		return nil, err
	}

	grids := []Grid{}

	for _, line := range lines {
		grid := ParseSingleGrid(line)
		grids = append(grids, grid)
	}

	return grids, nil
}

// TODO - move to utils folder? isn't specifically related to anything Sudoku-related
func readFileLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	lines := []string{}

	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	err = scanner.Err()
	if err != nil {
		return nil, err
	}

	return lines, nil
}

// public for testing purposes, and because a pure parsing function is useful to have
// TODO - potentially move functions for loading from a file into a separate package?
func ParseSingleGrid(str string) Grid {
	switch len(str) {
	case 16:
		return parseBaseSize2Grid(str)
	case 81:
		return parseBaseSize3Grid(str)
	case 256:
		panic("Parsing not yet implemented for 16x16 puzzles")
	case 625:
		panic("Parsing not yet implemented for 25x25 puzzles")
	case 1296:
		panic("Parsing not yet implemented for 36x36 puzzles")
	default:
		panic(fmt.Sprintf("Unrecognized grid size: %v", len(str)))
	}
}

func parseBaseSize2Grid(str string) Grid {
	grid := EmptyGrid(2)

	for pos, ch := range str {
		switch ch {
		case '1', '2', '3', '4':
			intValue, _ := strconv.Atoi(string(ch)) // ignore error, conversion should always be valid
			grid.cells[pos] = &Cell{
				index:          pos,
				containingGrid: &grid,
				value:          &intValue,
			}
		case '.':
			grid.cells[pos] = &Cell{
				index:          pos,
				containingGrid: &grid,
				value:          nil,
			}
		}
		// no default case - ignore all other runes
	}

	return grid
}

func parseBaseSize3Grid(str string) Grid {
	grid := EmptyGrid(3)

	for pos, ch := range str {
		switch ch {
		case '1', '2', '3', '4', '5', '6', '7', '8', '9':
			intValue, _ := strconv.Atoi(string(ch)) // ignore error, conversion should always be valid
			grid.cells[pos] = &Cell{
				index:          pos,
				containingGrid: &grid,
				value:          &intValue,
			}
		case '.':
			grid.cells[pos] = &Cell{
				index:          pos,
				containingGrid: &grid,
				value:          nil,
			}
		}
		// no default case - ignore all other runes
	}

	return grid
}
