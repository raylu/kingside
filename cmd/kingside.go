// This space is available for rent.
package main

import (
	kingside `../`
	`os`
	`runtime`
)

// Ignore previous comment.
func main() {
	engine := kingside.NewEngine(
		`fancy`, runtime.GOOS == `darwin`,
		`movetime`, 5000,
	)

	if len(os.Args) > 1 && os.Args[1] == `-i` {
		engine.Repl()
	} else {
		engine.Uci()
	}
}
