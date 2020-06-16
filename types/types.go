package types

type Playgroud struct {
	Pieces    map[int]Piece
	Positions []Position
	Turn      string
	Picked    int
}

func NewPlayground() *Playgroud {
	return &Playgroud{
		Positions: defaultPositions,
		Pieces:    defaultPieces,
	}
}

type Position struct {
	X     int
	Y     int
	Piece int
}

type Piece struct {
	Length Length
	Shape  Shape
	Color  Color
	Hole   Hole
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
	{X: 1, Y: 1, Piece: 0},
	{X: 1, Y: 2, Piece: 0},
	{X: 1, Y: 3, Piece: 0},
	{X: 1, Y: 4, Piece: 0},
	{X: 1, Y: 1, Piece: 0},
	{X: 1, Y: 2, Piece: 0},
	{X: 1, Y: 3, Piece: 0},
	{X: 1, Y: 4, Piece: 0},
	{X: 3, Y: 1, Piece: 0},
	{X: 3, Y: 2, Piece: 0},
	{X: 3, Y: 3, Piece: 0},
	{X: 3, Y: 4, Piece: 0},
	{X: 4, Y: 1, Piece: 0},
	{X: 4, Y: 2, Piece: 0},
	{X: 4, Y: 3, Piece: 0},
	{X: 4, Y: 4, Piece: 0},
}
