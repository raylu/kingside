package kingside

// Root node search. Basic principle is expressed by Boob's Law: you always find
// something in the last place you look.
func (p *Position) search(alpha, beta, depth int) (score int, move Move) {
	// Root move generator makes sure all generated moves are valid.
	gen := NewRootGen(p, depth)
	gen.generateRootMoves()

	inCheck := p.isInCheck(p.color)
	moveCount := 0
	bestScore := -Checkmate
	bestMove := Move(0)
	for move := gen.NextMove(); move != 0; move = gen.NextMove() {
		moveCount++

		position := p.makeMove(move)
		score = position.Evaluate() // TODO
		position.undoLastMove()
		if score > bestScore {
			bestScore = score
			bestMove = move
		}
	}


	if moveCount == 0 {
		if inCheck {
			score = -Checkmate
		} else {
			score = 0
		}
		if engine.uci {
			engine.uciScore(depth, score, alpha, beta)
		}
	}

	return bestScore, bestMove
}
