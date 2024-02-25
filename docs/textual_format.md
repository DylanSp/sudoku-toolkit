# Textual format for representing puzzles

1. Each puzzle is represented by a single line.
1. `.` characters represent an unassigned cell.
1. Depending on the size of the puzzle, a set of characters is used to represent given, assigned values.
   1. 4x4 puzzles use the set `[1-4]`.
   1. 9x9 puzzles use the set `[1-9]`.
   1. 16x6 puzzles use the set `[0-9A-G]`. Note that `0` is a valid value, _not_ an unassigned cell.
   1. 25x25 puzzles use the set `[0-9A-P]`.
   1. 36x36 puzzles use the set `[0-9A-Z]`.
1. All characters other than `.` and the characters used to represent givens are ignored.
1. Puzzles larger than 36x36 don't have a specified format.
