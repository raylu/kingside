package kingside

import (
	`fmt`
	`math/rand`
	`time`
)

// Returns row number in 0..7 range for the given square.
func row(square int) int {
	return square >> 3
}

// Returns column number in 0..7 range for the given square.
func col(square int) int {
	return square & 7
}

// Returns both row and column numbers for the given square.
func coordinate(square int) (int, int) {
	return row(square), col(square)
}

// Returns relative rank for the square in 0..7 range. For example E2 is rank 1
// for white and rank 6 for black.
func rank(color uint8, square int) int {
	return row(square) ^ (int(color) * 7)
}

// Returns 0..63 square number for the given row/column coordinate.
func square(row, column int) int {
	return (row << 3) + column
}

// Flips the square verically for white (ex. E2 becomes E7).
func flip(color uint8, square int) int {
	if color == White {
		return square ^ 56
	}
	return square
}

// Returns a bitmask with light or dark squares set matching the color of the
// square.
func same(square int) Bitmask {
	if bit[square] & maskDark != 0 {
		return maskDark
	}
	return ^maskDark
}

// Returns true if the square resides between two other squares on the same line
// or diagonal, including the edge squares. For example, between(A1, H8, C3) is
// true.
func between(from, to, between int) bool {
	return (maskStraight[from][to] | maskDiagonal[from][to]).on(between)
}

// Returns distance between current and root node.
func ply() int {
	return node - rootNode
}

func uncache(score, ply int) int {
	if score > Checkmate - MaxPly && score <= Checkmate {
		return score - ply
	} else if score >= -Checkmate && score < -Checkmate + MaxPly {
		return score + ply
	}
	return score
}

// Integer version of math/abs.
func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func min64(x, y int64) int64 {
	if x < y {
		return x
	}
	return y
}

func max64(x, y int64) int64 {
	if x > y {
		return x
	}
	return y
}

// Returns time in milliseconds elapsed since the given start time.
func since(start time.Time) int64 {
	return time.Since(start).Nanoseconds() / 1000000
}

// Formats time duration in milliseconds in human readable form (MM:SS.XXX).
func ms(duration int64) string {
	mm := duration / 1000 / 60
	ss := duration / 1000 % 60
	xx := duration - mm * 1000 * 60 - ss * 1000
	return fmt.Sprintf(`%02d:%02d.%03d`, mm, ss, xx)
}

// The generation of random numbers is too important to be left to chance.
// Returns pseudo-random integer in [0, limit] range. It panics if limit <= 0.
func Random(limit int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(limit)
}

func C(color uint8) string {
	return [2]string{`white`, `black`}[color]
}

// Logging wrapper around fmt.Printf() that could be turned on as needed. Typical
// usage is Log(); defer Log() in tests.
func Log(args ...interface{}) {
	switch len(args) {
	case 0:
		// Calling Log() with no arguments flips the logging setting.
		engine.log = !engine.log
		engine.fancy = !engine.fancy
	case 1:
		switch args[0].(type) {
		case bool:
			engine.log = args[0].(bool)
			engine.fancy = args[0].(bool)
		default:
			if engine.log {
				fmt.Println(args...)
			}
		}
	default:
		if engine.log {
			fmt.Printf(args[0].(string), args[1:]...)
		}
	}
}
