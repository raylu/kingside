package kingside

import (
	`strings`
	`time`
)

type Game struct {
	initial     string  // Initial position (FEN or algebraic).
}

// Use single statically allocated variable.
var game Game

// We have two ways to initialize the game: 1) pass FEN string, and 2) specify
// white and black pieces using regular chess notation.
// In latter case we need to tell who gets to move first when starting the game.
// The second option is a bit less pricise (ex. no en-passant square) but it is
// much more useful when writing tests from memory.
func NewGame(args ...string) *Game {
	game = Game{}
	switch len(args) {
	case 0: // Initial position.
		game.initial = `rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1`
	case 1: // Genuine FEN.
		game.initial = args[0]
	case 2: // Donna chess format (white and black).
		game.initial = args[0] + ` : ` + args[1]
	}
	return &game
}

func (game *Game) start() *Position {
	engine.clock.halt = false
	tree, node, rootNode = [1024]Position{}, 0, 0

	// Was the game started with FEN or algebraic notation?
	sides := strings.Split(game.initial, ` : `)
	if len(sides) == 2 {
		return NewPosition(game, sides[White], sides[Black])
	}
	return NewPositionFromFEN(game, game.initial)
}

func (game *Game) position() *Position {
	return &tree[node]
}

// "The question of whether machines can think is about as relevant as the
// question of whether submarines can swim." -- Edsger W. Dijkstra
func (game *Game) Think() Move {
	start := time.Now()
	position := game.position()
	rootNode = node

	_, move := position.search(-Checkmate, Checkmate, 1)
	game.printBestMove(move, since(start))
	return move
}

func (game *Game) printBestMove(move Move, duration int64) {
	if engine.uci {
		engine.uciBestMove(move, duration)
	} else {
		engine.replBestMove(move)
	}
}

func (game *Game) String() string {
	return game.position().String()
}
