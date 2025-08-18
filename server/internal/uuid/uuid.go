// Copyright 2025 Alan Barbosa Lima
// Licensed under the Apache License, Version 2.0

// This package implements a UUID generator, accourding to UUID
// Version 7 (based on the UNIX Epoch). This description can be found
// in RFC9562 section 5.3.
package uuid

import (
	"fmt"
	"math/rand/v2"
	"sync"
	"time"
)

// This represents a uint128 UUID type. Elements of this type can,
// and should, be compared using the == operator. The zero value of
// UUID is classified as a Nil UUID in RFC9562 section 5.9.
type UUID struct {
	lo, hi uint64
}

const (
	_62BitMask = (1 << 62) - 1
	_48BitMask = (1 << 48) - 1
	_12BitMask = (1 << 12) - 1
)

var (
	source rand.Source // The source for all pseudo-random number needed.
	mulock sync.Mutex  // A mutex for safe concurrent UUID generation.
)

func init() {
	lo, hi := uint64(time.Now().UnixNano()), uint64(0)
	hi *= lo * 0xCAFE_D0CE

	source = rand.NewPCG(lo, hi)
}

// Generates a new UUID accourding to version 7. It's safe to call
// this function from multiple goroutines.
func NewUUIDv7() UUID {
	const (
		version = 0b0111
		variant = 0b01
	)

	unixTimestamp := uint64(time.Now().UnixMilli() & _48BitMask)

	mulock.Lock()
	randA := source.Uint64() & _12BitMask
	randB := source.Uint64() & _62BitMask
	mulock.Unlock()

	return UUID{
		lo: unixTimestamp<<0 | version<<48 | randA<<52,
		hi: variant | randB<<2,
	}
}

// Implements the interface [fmt.Stringer] on the UUID type.
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
