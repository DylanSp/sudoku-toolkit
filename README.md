# sudoku-toolkit

A toolkit in Go for generating, solving, and otherwise analyzing Sudoku puzzles.

This project was inspired by [Eli Bendersky's post on solving and generating puzzles in Go](https://eli.thegreenplace.net/2022/sudoku-go-and-webassembly/). I have several goals that I'd like to accomplish:
1. Create a basic backtracking solver for the traditional 9x9 Sudoku puzzle size. (https://github.com/DylanSp/sudoku-toolkit/issues/1)
2. Generalize the solver to work with larger puzzle sizes. (https://github.com/DylanSp/sudoku-toolkit/issues/2)
3. Explore options for _generating_ puzzles with a given difficulty, ideally basing the difficulty on the techniques required to solve it. (https://github.com/DylanSp/sudoku-toolkit/issues/3)

As of late April 2024, I've been able to get a backtracking solver working for 9x9 puzzles on the [`backtracking-solver-cleaning-up-debugging` branch](https://github.com/DylanSp/sudoku-toolkit/tree/backtracking-solver-cleaning-up-debugging). While the code works, it needs some cleanup and organization, as well as a clearer explanation of how the code implements backtracking search.
