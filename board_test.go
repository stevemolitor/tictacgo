package main

import (
	"testing"
)

func TestMove(t *testing.T) {
	board := NewBoard()
	board.Move(0, X)
	cell := board[0]
	if cell != X {
		t.Fatalf(`b.Move(0, "X") = %q, want "X"`, cell)
	}
}

func TestIsDisabled(t *testing.T) {
	board := NewBoard()
	if board.IsDisabled(0) {
		t.Fatalf("b.IsDisabled(0) = false, want true")
	}
}

var xWinsBoards = []Board{
	{
		X, X, X,
		O, BLANK, O,
		BLANK, BLANK, BLANK,
	},
	{
		O, BLANK, O,
		X, X, X,
		BLANK, BLANK, BLANK,
	},
	{
		O, BLANK, O,
		BLANK, BLANK, BLANK,
		X, X, X,
	},
	{
		X, BLANK, O,
		X, BLANK, BLANK,
		X, O, O,
	},
	{
		X, O, O,
		X, X, O,
		X, O, X,
	},
	{
		O, O, X,
		X, X, O,
		X, O, X,
	},
}

var oWinsBoards = []Board{
	{
		O, O, O,
		X, BLANK, X,
		X, BLANK, BLANK,
	},
	{
		X, BLANK, X,
		O, O, O,
		X, BLANK, BLANK,
	},
	{
		X, BLANK, X,
		X, BLANK, BLANK,
		O, O, O,
	},
	{
		O, X, BLANK,
		O, X, X,
		O, BLANK, BLANK,
	},
	{
		O, X, X,
		X, O, X,
		O, X, O,
	},
	{
		X, X, O,
		X, O, X,
		O, X, O,
	},
}

var tieBoards = []Board{
	{
		O, O, X,
		X, X, O,
		O, X, X,
	},
	{
		X, O, X,
		O, O, X,
		X, X, O,
	},
}

var stillPlayingBoards = []Board{
	{
		O, O, X,
		X, BLANK, O,
		O, X, X,
	},
	{
		X, BLANK, BLANK,
		O, O, X,
		X, X, O,
	},
}

func TestGetGameState(t *testing.T) {
	for _, board := range xWinsBoards {
		gameState := board.GetGameState()
		if gameState != XWin {
			t.Fatalf(`board.GetGameState() = %v, want %v`, gameState, XWin)
		}
	}

	for _, board := range oWinsBoards {
		gameState := board.GetGameState()
		if gameState != OWin {
			t.Fatalf(`board.GetGameState() = %v, want %v`, gameState, OWin)
		}
	}

	for _, board := range tieBoards {
		gameState := board.GetGameState()
		if gameState != Tie {
			t.Fatalf(`board.GetGameState() = %v, want %v`, gameState, Tie)
		}
	}

	for _, board := range stillPlayingBoards {
		gameState := board.GetGameState()
		if gameState != StillPlaying {
			t.Fatalf(`board.GetGameState() = %v, want %v`, gameState, Tie)
		}
	}
}
