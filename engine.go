package kingside

import (
	`fmt`
	`os`
	`time`
)

type Clock struct {
	halt        bool     // Stop search immediately when set to true.
	softStop    int64    // Target soft time limit to make a move.
	hardStop    int64    // Immediate stop time limit.
	extra       float32  // Extra time factor based on search volatility.
	start       time.Time
	ticker      *time.Ticker
}

type Options struct {
	ponder      bool     // (-) Pondering mode.
	infinite    bool     // (-) Search until the "stop" command.
	maxDepth    int      // Search X plies only.
	maxNodes    int      // (-) Search X nodes only.
	moveTime    int64    // Search exactly X milliseconds per move.
	movesToGo   int64    // Number of moves to make till time control.
	timeLeft    int64    // Time left for all remaining moves.
	timeInc     int64    // Time increment after the move is made.
}

type Engine struct {
	log         bool     // Enable logging.
	uci         bool     // Use UCI protocol.
	fancy       bool     // Represent pieces as UTF-8 characters.
	status      uint8    // Engine status.
	clock       Clock
	options     Options
}

// Use single statically allocated variable.
var engine Engine

func NewEngine(args ...interface{}) *Engine {
	engine = Engine{}
	for i := 0; i < len(args); i += 2 {
		switch value := args[i+1]; args[i] {
		case `uci`:
			engine.uci = value.(bool)
		case `fancy`:
			engine.fancy = value.(bool)
		case `depth`:
			engine.options.maxDepth = value.(int)
		case `movetime`:
			engine.options.moveTime = int64(value.(int))
		}
	}

	return &engine
}

// Dumps the string to standard output.
func (e *Engine) print(arg string) *Engine {
	os.Stdout.WriteString(arg)
	return e
}

func (e *Engine) reply(args ...interface{}) *Engine {
	if len := len(args); len > 1 {
		data := fmt.Sprintf(args[0].(string), args[1:]...)
		e.print(data)
	} else if len == 1 {
		e.print(args[0].(string))
	}
	return e
}

func (e *Engine) limits(options Options) *Engine {
	e.options = options
	return e
}
