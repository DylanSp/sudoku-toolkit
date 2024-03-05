# Example puzzles

All 9x9 puzzles are taken from [Peter Norvig's post "Solving Every Sudoku Puzzle"](https://norvig.com/sudoku.html). The original URLs are https://norvig.com/easy50.txt, https://norvig.com/top95.txt (which I renamed to hard95.txt), and https://norvig.com/hardest.txt.

easy50-puzzle1-pretty.txt and easy50-puzzle2-pretty.txt are prettified versions of the first two puzzles from 9x9/easy50.txt.

## Notes on difficulty

easy50-puzzle1 can be solved with basic strategies only; easy50-puzzle2 cannot.

## Processing to match my format

### easy50.txt

1. Replace `0`s with `.`s - `sed -i 's/0/./g' easy50.txt`
2. Concatenate all lines for a single puzzle onto a single line - `awk '/([0-9\.])$/ { printf("%s", $0); next } 1' easy50.txt > awk_50.txt` (taken from https://stackoverflow.com/a/8519651/5847190)
