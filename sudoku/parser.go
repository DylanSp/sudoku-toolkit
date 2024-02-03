package sudoku

import (
	"bufio"
	"os"
)

// TODO - either hardcode these to base size 3 (9x9) and document that, or make them capable of handling different base sizes

func LoadGridsFromFile(filename string) ([]Grid, error) {
	lines, err := readFileLines(filename)
	if err != nil {
		return nil, err
	}

	grids := []Grid{}

	for _, line := range lines {
		grid, err := parseSingleGrid(line)
		if err != nil {
			return nil, err
		}

		grids = append(grids, grid)
	}

	return grids, nil
}

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

func parseSingleGrid(str string) (Grid, error) {
	panic("not yet implemented")
}
