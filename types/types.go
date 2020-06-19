package types

import "fmt"

type Board struct {
	Pieces    map[int]Piece `json:"pieces"`
	Positions []Position    `json:"positions"`
	Turn      string        `json:"turn"`
	Picked    int           `json:"picked"`
}

func NewBoard(turn string) Board {
	return Board{
		Positions: defaultPositions,
		Pieces:    defaultPieces,
		Turn:      turn,
		Picked:    0,
	}
}

type Position struct {
	X       int `json:"x"`
	Y       int `json:"y"`
	PieceID int `json:"piece-id"`
}

type Piece struct {
	Length Length `json:"length"`
	Shape  Shape  `json:"shape"`
	Color  Color  `json:"color"`
	Hole   Hole   `json:"hole"`
}

func (p Piece) Code() int {
	var c int = 0
	if p.Color == White {
		c += 1
	} else {
		c += 2
	}

	if p.Hole == Hollow {
		c += 4
	} else {
		c += 8
	}

	if p.Length == Short {
		c += 16
	} else {
		c += 32
	}

	if p.Shape == Square {
		c += 64
	} else {
		c += 128
	}

	return c
}

type Length string

const (
	Short Length = "short"
	Tall  Length = "tall"
)

type Shape string

const (
	Round  Shape = "round"
	Square Shape = "square"
)

type Color string

const (
	Black Color = "black"
	White Color = "white"
)

type Hole string

const (
	Hollow Hole = "hollow"
	Solid  Hole = "solid"
)

var defaultPieces = map[int]Piece{
	1:  {Length: Tall, Shape: Round, Color: White, Hole: Hollow},
	2:  {Length: Short, Shape: Round, Color: White, Hole: Hollow},
	3:  {Length: Tall, Shape: Square, Color: White, Hole: Hollow},
	4:  {Length: Short, Shape: Square, Color: White, Hole: Hollow},
	5:  {Length: Tall, Shape: Round, Color: Black, Hole: Hollow},
	6:  {Length: Short, Shape: Round, Color: Black, Hole: Hollow},
	7:  {Length: Tall, Shape: Square, Color: Black, Hole: Hollow},
	8:  {Length: Short, Shape: Square, Color: Black, Hole: Hollow},
	9:  {Length: Tall, Shape: Round, Color: White, Hole: Solid},
	10: {Length: Short, Shape: Round, Color: White, Hole: Solid},
	11: {Length: Tall, Shape: Square, Color: White, Hole: Solid},
	12: {Length: Short, Shape: Square, Color: White, Hole: Solid},
	13: {Length: Tall, Shape: Round, Color: Black, Hole: Solid},
	14: {Length: Short, Shape: Round, Color: Black, Hole: Solid},
	15: {Length: Tall, Shape: Square, Color: Black, Hole: Solid},
	16: {Length: Short, Shape: Square, Color: Black, Hole: Solid},
}

var defaultPositions = []Position{
	{X: 1, Y: 1, PieceID: 0},
	{X: 1, Y: 2, PieceID: 0},
	{X: 1, Y: 3, PieceID: 0},
	{X: 1, Y: 4, PieceID: 0},
	{X: 2, Y: 1, PieceID: 0},
	{X: 2, Y: 2, PieceID: 0},
	{X: 2, Y: 3, PieceID: 0},
	{X: 2, Y: 4, PieceID: 0},
	{X: 3, Y: 1, PieceID: 0},
	{X: 3, Y: 2, PieceID: 0},
	{X: 3, Y: 3, PieceID: 0},
	{X: 3, Y: 4, PieceID: 0},
	{X: 4, Y: 1, PieceID: 0},
	{X: 4, Y: 2, PieceID: 0},
	{X: 4, Y: 3, PieceID: 0},
	{X: 4, Y: 4, PieceID: 0},
}

type Action struct {
	Picked int `json:"picked"`
	X      int `json:"x"`
	Y      int `json:"y"`
}

func (a Action) String() string {
	return fmt.Sprintf("{x: %d, y: %d, picked: %d", a.X, a.Y, a.Picked)
}

type Player struct {
	Name string
	Id   string
}
