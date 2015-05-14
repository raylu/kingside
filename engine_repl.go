package kingside

import(
	`fmt`
	`runtime`
)

var (
	ansiRed   = "\033[0;31m"
	ansiGreen = "\033[0;32m"
	ansiTeal  = "\033[0;36m"
	ansiNone  = "\033[0m"
)

func (e *Engine) replBestMove(move Move) *Engine {
	fmt.Printf(ansiTeal + "kingside's move: %s", move)
	fmt.Println(ansiNone + "\n")
	return e
}

// "There are two types of command interfaces in the world of computing: good
// interfaces and user interfaces." -- Daniel J. Bernstein
func (e *Engine) Repl() *Engine {
	var game *Game
	var position *Position

	// Suppress ANSI colors when running Windows.
	if runtime.GOOS == `windows` {
		ansiRed, ansiGreen, ansiTeal, ansiNone = ``, ``, ``, ``
	}

	setup := func() {
		if game == nil || position == nil {
			game = NewGame()
			position = game.start()
			fmt.Printf("%s\n", position)
		}
	}

	think := func() {
		if move := game.Think(); move != 0 {
			position = position.makeMove(move)
			fmt.Printf("%s\n", position)
		}
	}

	fmt.Printf("kingside\nType ? for help.\n\n")
	for command, parameter := ``, ``; ; command, parameter = ``, `` {
		fmt.Print(`kingside> `)
		fmt.Scanln(&command, &parameter)

		switch command {
		case ``:
		case `exit`, `quit`:
			return e
		case `go`:
			setup()
			think()
		case `help`, `?`:
			fmt.Println("The commands are:\n\n" +
				"  exit           Exit the program\n" +
				"  go             Take side and make a move\n" +
				"  help           Display this help\n" +
				"  new            Start new game\n" +
				"  undo           Undo last move\n\n" +
				"To make a move use algebraic notation, for example e2e4, Ng1f3, or e7e8Q\n")
		case `new`:
			game, position = nil, nil
			setup()
		case `undo`:
			if position != nil {
				position = position.undoLastMove()
				fmt.Printf("%s\n", position)
			}
		default:
			setup()
			if move, validMoves := NewMoveFromString(position, command); move != 0 {
				position = position.makeMove(move)
				think()
			} else { // Invalid move or non-evasion on check.
				fancy := e.fancy; e.fancy = false
				fmt.Printf("%s appears to be an invalid move; valid moves are %v\n", command, validMoves)
				e.fancy = fancy
			}
		}
	}
	return e
}
