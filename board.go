// board.go - Board object encapsulating state and state transitions for one Tic Tac Toe game.
package main

import (
	"math/rand"
	"time"
)

// Board and CellState types
type CellState byte
type Board []CellState

// Constants
const (
	X     CellState = 'X'
	O               = 'O'
	BLANK           = ' '
)

const (
	XWin         = "x-win"
	OWin         = "o-win"
	Tie          = "tie"
	StillPlaying = "still-playing"
)

// Vars used when computing game state

// Magic square where the sum of any row, column, or diagonal is 15.
var magicSquare = []int{
	8, 1, 6,
	3, 5, 7,
	4, 9, 2,
}

var boardLines = [][]int{
	// rows
	{0, 1, 2},
	{3, 4, 5},
	{6, 7, 8},

	// columns
	{0, 3, 6},
	{1, 4, 7},
	{2, 5, 8},

	// diagonals
	{0, 4, 8},
	{2, 4, 6},
}

// Private helpers
func lineContains(line []int, sum int) bool {
	for _, n := range line {
		if sum == n {
			return true
		}
	}
	return false
}

func (b Board) contains(cellState CellState) bool {
	for _, c := range b {
		if cellState == c {
			return true
		}
	}
	return false
}

func randomCell(n int) int {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	return r.Intn(n)
}

// Compute magic square sum for given line (which could be a row, column, or diagonal).
//
// Multiply each square value by 1 if the cell has an X, two if an O. Blank cells are ignored. If
// the sum is 15, X has won. If 30 (15 * 2) O has won.
func (b Board) sumLine(line []int) int {
	sum := 0
	for _, index := range line {
		cellState := b[index]
		squareValue := magicSquare[index]
		if cellState == X {
			sum += squareValue
		}
		if cellState == O {
			sum += squareValue * 2
		}
	}
	return sum
}

// Private Board method
func (b Board) getBlankCells() []int {
	blankCells := make([]int, 0, 9)
	for i, cellState := range b {
		if cellState == BLANK {
			blankCells = append(blankCells, i)
		}
	}

	return blankCells
}

// Public Board methods and constructor
func NewBoard() Board {
	return Board{BLANK, BLANK, BLANK, BLANK, BLANK, BLANK, BLANK, BLANK, BLANK}
}

func (b Board) Move(cell int, player CellState) {
	b[cell] = player

	blankCells := b.getBlankCells()

	if len(blankCells) > 0 {
		n := randomCell(len(blankCells))
		b[blankCells[n]] = O
	}
}

func (b Board) IsDisabled(cell int) bool {
	return b[cell] != BLANK || b.IsGameOver()
}

func (b Board) GetGameState() string {
	lineSums := make([]int, 8)
	for i, line := range boardLines {
		lineSums[i] = b.sumLine(line)
	}

	// If the sum of any line is 15, X has won.
	if lineContains(lineSums, 15) {
		return XWin
	}
	// If the sum of any line is 30, O has won.
	if lineContains(lineSums, 30) {
		return OWin
	}
	// Else if there is still a blank square on the board, we're still playing.
	if b.contains(BLANK) {
		return StillPlaying
	}
	// Else game is over and it's a tie.
	return Tie
}

func (b Board) IsGameOver() bool {
	gs := b.GetGameState()
	return gs != StillPlaying
}
