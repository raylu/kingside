package kingside

// Root node search. Basic principle is expressed by Boob's Law: you always find
// something in the last place you look.
func (p *Position) search(alpha, beta, depth int) (bestScore int, bestMove Move) {
	gen := NewRootGen(p, depth)
	gen.generateRootMoves()

	inCheck := p.isInCheck(p.color)
	moveCount := 0
	bestScore = -Checkmate
	bestMove = Move(0)
	for move := gen.NextMove(); move != 0; move = gen.NextMove() {
		moveCount++

		position := p.makeMove(move)
		score := -position.negamax(depth + 1)
		position.undoLastMove()
		if score > bestScore {
			bestScore = score
			bestMove = move
		}
	}

	if moveCount == 0 {
		if inCheck {
			bestScore = -Checkmate
		} else {
			bestScore = 0
		}
		if engine.uci {
			engine.uciScore(depth, bestScore, alpha, beta)
		}
	}

	return bestScore, bestMove
}

func (p *Position) negamax(depth int) (bestScore int) {
	gen := NewMoveGen(p)
	inCheck := p.isInCheck(p.color)
	if inCheck {
		gen.generateEvasions()
	} else {
		gen.generateMoves()
	}

	moveCount := 0
	bestScore = -Checkmate
	for move := gen.NextMove(); move != 0; move = gen.NextMove() {
		moveCount++

		position := p.makeMove(move)
		score := 0 // TODO
		position.undoLastMove()
		bestScore = max(score, bestScore)
	}

	if moveCount == 0 {
		if inCheck {
			bestScore = -Checkmate
		} else {
			bestScore = 0
		}
	}

	return bestScore
}
