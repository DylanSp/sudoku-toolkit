# Terminology

Specific terminology used in this codebase.

## Standard Sudoku terminology

- `Grid` - the general term for the full grid of cells, some or all of which may have values.
- `Cell` - a single space in the Grid.
- `Row` - a single row of cells across the entire Grid.
- `Column` - a single columns across the entire Grid.
- `Box` - a square sub-grid - for an NxN Grid, the box has N total Cells.
- `House` - a Row, Column, or Box.
- `Hints`/`Givens` - the Cells that start filled in a Puzzle.
- `Peers` - Cells in the same House as another Cell (usage: the Peers of cell C are all cells that share a Row, Column, or Box with C)

## Nonstandard, codebase-specific terminology

- `Challenge` - The initial "puzzle"; an incomplete grid with some, but not all, Cells filled, posed as a challenge to solve.
- `Puzzle` - A Grid that's in the process of being solved; it has all the Givens from the original Challenge, it may have other Cells filled by a solver that's working on it.
- `Solution` - a completely filled-out Grid.
- `Base size` - an integer representation of a grid's size - a grid with "base size N" has sides that are n^2 long, and uses n^2 different digits/elements as possible values. The base size of the standard 9x9 grid is 3. (Using this as the measurement of size avoids having to use square roots in the code)
