package kingside

import (
)

type MoveGen struct {
	p        *Position
	list     [128]Move
	ply      int
	head     int
	tail     int
	pins     Bitmask
}

// Pre-allocate move generator array (one entry per ply) to avoid garbage
// collection overhead. Last entry serves for utility move generation, ex. when
// converting string notations or determining a stalemate.
var moveList [MaxPly+1]MoveGen

// Returns "new" move generator for the given ply. Since move generator array
// has been pre-allocated already we simply return a pointer to the existing
// array element re-initializing all its data.
func NewGen(p *Position, ply int) (gen *MoveGen) {
	gen = &moveList[ply]
	gen.p = p
	gen.list = [128]Move{}
	gen.ply = ply
	gen.head, gen.tail = 0, 0
	gen.pins = p.pinnedMask(p.king[p.color])

	return gen
}

// Convenience method to return move generator for the current ply.
func NewMoveGen(p *Position) *MoveGen {
	return NewGen(p, ply())
}

// Returns new move generator for the initial step of iterative deepening
// (depth == 1) and existing one for subsequent iterations (depth > 1).
func NewRootGen(p *Position, depth int) *MoveGen {
	if depth == 1 {
		return NewGen(p, 0) // Zero ply.
	}
	return &moveList[0]
}

func (gen *MoveGen) reset() *MoveGen {
	gen.head = 0
	return gen
}

func (gen *MoveGen) size() int {
	return gen.tail
}

func (gen *MoveGen) onlyMove() bool {
	return gen.tail == 1
}

func (gen *MoveGen) NextMove() (move Move) {
	if gen.head < gen.tail {
		move = gen.list[gen.head]
		gen.head++
	}
	return
}

// Returns true if the move is valid in current position i.e. it can be played
// without violating chess rules.
func (gen *MoveGen) isValid(move Move) bool {
	return gen.p.isValid(move, gen.pins)
}

// Removes invalid moves from the generated list. We use in iterative deepening
// to avoid filtering out invalid moves on each iteration.
func (gen *MoveGen) validOnly() *MoveGen {
	for move := gen.NextMove(); move != 0; move = gen.NextMove() {
		if !gen.isValid(move) {
			gen.remove()
		}
	}
	return gen.reset()
}

// Probes a list of generated moves and returns true if it contains at least
// one valid move.
func (gen *MoveGen) anyValid() bool {
	for move := gen.NextMove(); move != 0; move = gen.NextMove() {
		if gen.isValid(move) {
			return true
		}
	}
	return false
}

// Probes valid-only list of generated moves and returns true if the given move
// is one of them.
func (gen *MoveGen) amongValid(someMove Move) bool {
	for move := gen.NextMove(); move != 0; move = gen.NextMove() {
		if someMove == move {
			return true
		}
	}
	return false
}

func (gen *MoveGen) add(move Move) *MoveGen {
	gen.list[gen.tail] = move
	gen.tail++
	return gen
}

// Removes current move from the list by copying over the remaining moves. Head and
// tail pointers get decremented so that calling NextMove() works as expected.
func (gen *MoveGen) remove() *MoveGen {
	copy(gen.list[gen.head-1:], gen.list[gen.head:])
	gen.head--
	gen.tail--
	return gen
}

// Returns an array of generated moves by continuously appending the NextMove()
// until the list is empty.
func (gen *MoveGen) allMoves() (moves []Move) {
	for move := gen.NextMove(); move != 0; move = gen.NextMove() {
		moves = append(moves, move)
	}
	gen.reset()

	return
}
