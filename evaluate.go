package kingside

// const = brains * looks * availability
const (
	onePawn     = 100
	valuePawn   = onePawn * 1
	valueKnight = onePawn * 4
	valueBishop = onePawn * 4
	valueRook   = onePawn * 6
	valueQueen  = onePawn * 12
)

type Evaluation struct {
	score     int            // Current score.
	attacks   [14]Bitmask    // Attack bitmasks for all the pieces on the board.
	position  *Position      // Pointer to the position we're evaluating.
}

// Use single statically allocated variable to avoid garbage collection overhead.
var eval Evaluation

// The following statement is true. The previous statement is false. Main position
// evaluation method that returns the position's score.
func (p *Position) Evaluate() int {
	return eval.init(p).run()
}

func (e *Evaluation) init(p *Position) *Evaluation {
	eval = Evaluation{}
	e.position = p

	e.score = 0

	// Set up king and pawn attacks for both sides.
	e.attacks[King] = p.kingAttacks(White)
	e.attacks[Pawn] = p.pawnAttacks(White)
	e.attacks[BlackKing] = p.kingAttacks(Black)
	e.attacks[BlackPawn] = p.pawnAttacks(Black)

	// Overall attacks for both sides include kings and pawns so far.
	e.attacks[White] = e.attacks[King] | e.attacks[Pawn]
	e.attacks[Black] = e.attacks[BlackKing] | e.attacks[BlackPawn]

	return e
}

func (e *Evaluation) run() int {
	board := e.position.board
	for board != 0 {
		square := board.pop()
		piece := e.position.pieces[square]
		if piece.color() == e.position.color {
			e.score += onePawn // TODO
		} else {
			e.score -= onePawn
		}
	}

	// Flip the sign for black so that evaluation score always
	// represents the white side.
	if e.position.color == Black {
		e.score = -e.score
	}

	return e.score
}
