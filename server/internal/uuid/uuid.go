// Copyright 2025 Alan Lima. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// This package implements a UUID generator,
// accourding to UUID Version 7 (based on the
// UNIX Epoch). This description can be found
// in RFC9562 section 5.3.

package uuid

import (
	"fmt"
	"math/rand/v2"
	"time"
)

// This represents a uint128 UUID type. Elements
// of this type can be compared using the ==
// operator.
type UUID struct {
	lo, hi uint64
}

const (
	_62BitMask = (1 << 62) - 1
	_48BitMask = (1 << 48) - 1
	_12BitMask = (1 << 12) - 1
)

const (
	version = 0b0111
	variant = 0b01
)

// Generates a new UUID accourding to version 7.
func NewUUID() UUID {
	unixTimestamp := uint64(time.Now().UnixMilli() & _48BitMask)
	randA := rand.Uint64() & _12BitMask
	randB := rand.Uint64() & _62BitMask

	return UUID{
		lo: unixTimestamp<<0 | version<<48 | randA<<52,
		hi: variant | randB<<2,
	}
}

// Implements the interface [fmt.Stringer] on the
// UUID type.
func (uuid UUID) String() string {
	fst4HexByte := (uuid.lo & 0x0000_0000_FFFF_FFFF) >> 0
	snd2HexByte := (uuid.lo & 0x0000_FFFF_0000_0000) >> 32
	trd2HexByte := (uuid.lo & 0xFFFF_0000_0000_0000) >> 48
	for2HexByte := (uuid.hi & 0x0000_0000_0000_FFFF) >> 0
	fif6HexByte := (uuid.hi & 0xFFFF_FFFF_FFFF_0000) >> 16

	return fmt.Sprintf(
		"%08x-%04x-%04x-%04x-%012x",
		fst4HexByte, snd2HexByte, trd2HexByte,
		for2HexByte, fif6HexByte,
	)
}