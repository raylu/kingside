package kingside

type Magic struct {
	mask  Bitmask
	magic Bitmask
}

var (
	kingMoves [64]Bitmask
	knightMoves [64]Bitmask
	pawnMoves [2][64]Bitmask
	rookMagicMoves [64][4096]Bitmask
	bishopMagicMoves [64][512]Bitmask

	maskPassed [2][64]Bitmask
	maskInFront [2][64]Bitmask

	// Complete file or rank mask if both squares reside on on the same file
	// or rank.
	maskStraight [64][64]Bitmask

	// Complete diagonal mask if both squares reside on on the same diagonal.
	maskDiagonal [64][64]Bitmask

	// If a king on square [x] gets checked from square [y] it can evade the
	// check from all squares except maskEvade[x][y]. For example, if white
	// king on B2 gets checked by black bishop on G7 the king can't step back
	// to A1 (despite not being attacked by black).
	maskEvade [64][64]Bitmask

	// If a king on square [x] gets checked from square [y] the check can be
	// evaded by moving a piece to maskBlock[x][y]. For example, if white
	// king on B2 gets checked by black bishop on G7 the check can be evaded
	// by moving white piece onto C3-G7 diagonal (including capture on G7).
	maskBlock [64][64]Bitmask

	// Bitmask to indicate pawn attacks for a square. For example, C3 is being
	// attacked by white pawns on B2 and D2, and black pawns on B4 and D4.
	maskPawn [2][64]Bitmask

	// Two arrays to simplify incremental polyglot hash computation.
	hashCastle [16]uint64
	hashEnpassant [8]uint64

	// Distance between two squares.
	distance [64][64]int

	// Most-significant bit (MSB) lookup table.
	msbLookup[256]int
)

func init() {
	initMasks()
	initArrays()
}

func initMasks() {
	for sq := A1; sq <= H8; sq++ {
		row, col := coordinate(sq)

		// Distance, Blocks, Evasions, Straight, Diagonals, Knights, and Kings.
		for i := A1; i <= H8; i++ {
			r, c := coordinate(i)

			distance[sq][i] = max(abs(row - r), abs(col - c))
			setupMasks(sq, i, row, col, r, c)

			if i == sq || abs(i-sq) > 17 {
				continue // No king or knight can reach that far.
			}
			if (abs(r-row) == 2 && abs(c-col) == 1) || (abs(r-row) == 1 && abs(c-col) == 2) {
				knightMoves[sq].set(i)
			}
			if abs(r-row) <= 1 && abs(c-col) <= 1 {
				kingMoves[sq].set(i)
			}
		}

		// Rooks.
		mask := createRookMask(sq)
		bits := uint(mask.count())
		for i := 0; i < (1 << bits); i++ {
			bitmask := mask.magicify(i)
			index := (bitmask * rookMagic[sq].magic) >> 52
			rookMagicMoves[sq][index] = createRookAttacks(sq, bitmask)
		}

		// Bishops.
		mask = createBishopMask(sq)
		bits = uint(mask.count())
		for i := 0; i < (1 << bits); i++ {
			bitmask := mask.magicify(i)
			index := (bitmask * bishopMagic[sq].magic) >> 55
			bishopMagicMoves[sq][index] = createBishopAttacks(sq, bitmask)
		}

		// Pawns.
		if row >= A2H2 && row <= A7H7 {
			if col > 0 {
				pawnMoves[White][sq].set(square(row + 1, col - 1))
				pawnMoves[Black][sq].set(square(row - 1, col - 1))
			}
			if col < 7 {
				pawnMoves[White][sq].set(square(row + 1, col + 1))
				pawnMoves[Black][sq].set(square(row - 1, col + 1))
			}
		}

		// Pawn attacks.
		if row > 1 { // White pawns can't attack first two ranks.
			if col != 0 {
				maskPawn[White][sq] |= bit[sq-9]
			}
			if col != 7 {
				maskPawn[White][sq] |= bit[sq-7]
			}
		}
		if row < 6 { // Black pawns can attack 7th and 8th ranks.
			if col != 0 {
				maskPawn[Black][sq] |= bit[sq+7]
			}
			if col != 7 {
				maskPawn[Black][sq] |= bit[sq+9]
			}
		}

		// Vertical sqs in front of a pawn.
		maskInFront[White][sq] = (maskBlock[sq][A8+col] | bit[A8+col]) & ^bit[sq]
		maskInFront[Black][sq] = (maskBlock[A1+col][sq] | bit[A1+col]) & ^bit[sq]

		// Masks to check for passed pawns.
		if col > 0 {
			maskPassed[White][sq] |= maskInFront[White][sq-1]
			maskPassed[Black][sq] |= maskInFront[Black][sq-1]
			maskPassed[White][sq-1] |= maskInFront[White][sq]
			maskPassed[Black][sq-1] |= maskInFront[Black][sq]
		}
		maskPassed[White][sq] |= maskInFront[White][sq]
		maskPassed[Black][sq] |= maskInFront[Black][sq]
	}
}

func initArrays() {
	// MSB lookup table.
	for i := 0; i < len(msbLookup); i++ {
		if i > 127 {
			msbLookup[i] = 7
		} else if i > 63 {
			msbLookup[i] = 6
		} else if i > 31 {
			msbLookup[i] = 5
		} else if i > 15 {
			msbLookup[i] = 4
		} else if i > 7 {
			msbLookup[i] = 3
		} else if i > 3 {
			msbLookup[i] = 2
		} else if i > 1 {
			msbLookup[i] = 1
		}
	}

	// Castle hash values.
	for mask := uint8(0); mask < 16; mask++ {
		if mask & castleKingside[White] != 0 {
			hashCastle[mask] ^= polyglotRandomCastle[0]
		}
		if mask & castleQueenside[White] != 0 {
			hashCastle[mask] ^= polyglotRandomCastle[1]
		}
		if mask & castleKingside[Black] != 0 {
			hashCastle[mask] ^= polyglotRandomCastle[2]
		}
		if mask & castleQueenside[Black] != 0 {
			hashCastle[mask] ^= polyglotRandomCastle[3]
		}
	}

	// Enpassant hash values.
	for col := A1; col <= H1; col++ {
		hashEnpassant[col] = polyglotRandomEnpassant[col]
	}
}

func createRookMask(square int) Bitmask {
	r, c := coordinate(square)
	bitmask := (maskRank[r] | maskFile[c]) ^ bit[square]

	return *bitmask.trim(r, c)
}

func createBishopMask(square int) Bitmask {
	r, c := coordinate(square)
	bitmask := Bitmask(0)

	if sq := square + 7; sq <= H8 && col(sq) == c - 1 {
		bitmask = maskDiagonal[square][sq]
	} else if sq := square - 7; sq >= A1 && col(sq) == c + 1 {
		bitmask = maskDiagonal[square][sq]
	}

	if sq := square + 9; sq <= H8 && col(sq) == c + 1 {
		bitmask |= maskDiagonal[square][sq]
	} else if sq := square - 9; sq >= A1 && col(sq) == c - 1 {
		bitmask |= maskDiagonal[square][sq]
	}
	bitmask ^= bit[square]

	return *bitmask.trim(r, c)
}

func createRookAttacks(sq int, mask Bitmask) (bitmask Bitmask) {
	row, col := coordinate(sq)

	// North
	for c, r := col, row + 1; r <= 7; r++ {
		b := bit[square(r, c)]
		bitmask |= b
		if mask & b != 0 {
			break
		}
	}
	// East
	for c, r := col + 1, row; c <= 7; c++ {
		b := bit[square(r, c)]
		bitmask |= b
		if mask & b != 0 {
			break
		}
	}
	// South
	for c, r := col, row - 1; r >= 0; r-- {
		b := bit[square(r, c)]
		bitmask |= b
		if mask & b != 0 {
			break
		}
	}
	// West
	for c, r := col - 1, row; c >= 0; c-- {
		b := bit[square(r, c)]
		bitmask |= b
		if mask & b != 0 {
			break
		}
	}
	return
}

func createBishopAttacks(sq int, mask Bitmask) (bitmask Bitmask) {
	row, col := coordinate(sq)

	// North East
	for c, r := col + 1, row + 1; c <= 7 && r <= 7; c, r = c+1, r+1 {
		b := bit[square(r, c)]
		bitmask |= b
		if mask & b != 0 {
			break
		}
	}
	// South East
	for c, r := col + 1, row - 1; c <= 7 && r >= 0; c, r = c+1, r-1 {
		b := bit[square(r, c)]
		bitmask |= b
		if mask & b != 0 {
			break
		}
	}
	// South West
	for c, r := col - 1, row - 1; c >= 0 && r >= 0; c, r = c-1, r-1 {
		b := bit[square(r, c)]
		bitmask |= b
		if mask & b != 0 {
			break
		}
	}
	// North West
	for c, r := col - 1, row + 1; c >= 0 && r <= 7; c, r = c-1, r+1 {
		b := bit[square(r, c)]
		bitmask |= b
		if mask & b != 0 {
			break
		}
	}
	return
}

func setupMasks(square, target, row, col, r, c int) {
	if row == r {
		if col < c {
			maskBlock[square][target].fill(square, 1, bit[target], maskFull)
			maskEvade[square][target].spot(square, -1, ^maskFile[0])
		} else if col > c {
			maskBlock[square][target].fill(square, -1, bit[target], maskFull)
			maskEvade[square][target].spot(square, 1, ^maskFile[7])
		}
		if col != c {
			maskStraight[square][target] = maskRank[r]
		}
	} else if col == c {
		if row < r {
			maskBlock[square][target].fill(square, 8, bit[target], maskFull)
			maskEvade[square][target].spot(square, -8, ^maskRank[0])
		} else {
			maskBlock[square][target].fill(square, -8, bit[target], maskFull)
			maskEvade[square][target].spot(square, 8, ^maskRank[7])
		}
		if row != r {
			maskStraight[square][target] = maskFile[c]
		}
	} else if r+col == row+c { // Diagonals (A1->H8)
		if col < c {
			maskBlock[square][target].fill(square, 9, bit[target], maskFull)
			maskEvade[square][target].spot(square, -9,  ^maskRank[0] & ^maskFile[0])
		} else {
			maskBlock[square][target].fill(square, -9, bit[target], maskFull)
			maskEvade[square][target].spot(square, 9, ^maskRank[7] & ^maskFile[7])
		}
		if shift := (r - c) & 15; shift < 8 { // A1-A8-H8
			maskDiagonal[square][target] = maskA1H8 << uint(8*shift)
		} else { // B1-H1-H7
			maskDiagonal[square][target] = maskA1H8 >> uint(8*(16-shift))
		}
	} else if row+col == r+c { // AntiDiagonals (H1->A8)
		if col < c {
			maskBlock[square][target].fill(square, -7, bit[target], maskFull)
			maskEvade[square][target].spot(square, 7, ^maskRank[7] & ^maskFile[0])
		} else {
			maskBlock[square][target].fill(square, 7, bit[target], maskFull)
			maskEvade[square][target].spot(square, -7, ^maskRank[0] & ^maskFile[7])
		}
		if shift := 7 ^ (r + c); shift < 8 { // A8-A1-H1
			maskDiagonal[square][target] = maskH1A8 >> uint(8*shift)
		} else { // B8-H8-H2
			maskDiagonal[square][target] = maskH1A8 << uint(8*(16-shift))
		}
	}

	// Default values are all 0 for maskBlock[square][target] (Go sets it for us)
	// and all 1 for maskEvade[square][target].
	if maskEvade[square][target] == 0 {
		maskEvade[square][target] = maskFull
	}
}
