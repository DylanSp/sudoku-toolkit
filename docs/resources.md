# Resources

## Programming

Links to resources for programming with Sudoku.

- [Solving and generating puzzles in Go](https://eli.thegreenplace.net/2022/sudoku-go-and-webassembly/). Also contains information on generating puzzles and displaying them in the browser; the Go code generates a puzzle as an SVG, and is compiled to WebAssembly to run in the browser.
- [Classic Peter Norvig post on solving Sudoku puzzles](https://norvig.com/sudoku.html). Also contains some lists of hard puzzles.
- Determining difficulty of Sudoku puzzles:
  - [Puzzling.SE discussion](https://puzzling.stackexchange.com/questions/29/what-are-the-criteria-for-determining-the-difficulty-of-sudoku-puzzle)
  - [Academic paper](https://arxiv.org/abs/1403.7373)
- [Hacker News post on a browser Sudoku implementation](https://news.ycombinator.com/item?id=38913220) - potentially useful for inspiration, especially looking at the various critiques brought up on HN.

## Sudoku

- https://www.sudoku9x9.com/blankgrid/ - allows submitting puzzles and trying to solve them, can check (manually-entered) solutions for validity, playing at https://www.sudoku9x9.com/ can reproduce SolveWithBasicStrategies() by filling all empty cells with candidates and converting single marks into regular entries
