// Copyright 2019-2025 Matthew Wilson and Synesis Information Systems. All
// rights reserved. Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
 * Created: 14th March 2019
 * Updated: 2nd November 2025
 */

// Definition of a number of types that have one-way behaviour.

package sync

import (
	"errors"
	"sync/atomic"
	sync_atomic "sync/atomic"
)

const (
	MaxLatchDistance int64 = 0x7FFF_FFFF_FFFF_FFFF
	latchedFloor     int64 = -1_234_567
)

var (
	errDownLatchInitialValueMustBeGreaterThanThreshold = errors.New("initial value must be greater than the threshold")
	errUpLatchInitialValueMustBeLessThanThreshold      = errors.New("initial value must be less than the threshold")
	errLatchDistanceExceedsMaximum                     = errors.New("latch distance exceeds maximum")
)

// A one-way switch that may be operated safely by multiple concurrent
// goroutines.
type BoolLatch struct {
	value int64
}

// Creates a new BoolLatch.
func NewBoolLatch() BoolLatch {

	return BoolLatch{}
}

// Sets the instance to the latched state if it is not currently latched;
// no effect if already latched.
//
// Returns:
// true if the latch was flipped; false otherwise (meaning it was already
// latched)
func (l *BoolLatch) Set() (flipped bool) {

	if sync_atomic.CompareAndSwapInt64(&l.value, 0, 1) {

		flipped = true
	} else {

		flipped = false
	}

	return
}

// Obtains the current value of the latch, without changing its state.
func (l *BoolLatch) Load() bool {

	return 0 != sync_atomic.LoadInt64(&l.value)
}

// Count-down-to-zero (and beyond, a bit) numeric latch
type _baseLatch struct {
	value int64
}

func (l *_baseLatch) step() (flipped, isLatched bool, newCount int64) {

	flipped = false
	isLatched = false

	newCount = atomic.AddInt64(&l.value, -1)

	if newCount < 1 {

		if newCount == 0 {

			flipped = true
		} else {
			newCount = 0
		}

		isLatched = true

		// once latched, always swap out to the floor, to avoid negative wrap
		atomic.SwapInt64(&l.value, latchedFloor)
	}

	return
}

func (l *_baseLatch) load() (isLatched bool, count int64) {

	count = atomic.LoadInt64(&l.value)

	if count < 1 {

		isLatched = true

		count = 0
	}

	return
}

// A unidirectional latch that counts down from an initial value to a lower
// threshold that may be operated safely by multiple concurrent goroutines.
type DownLatch struct {
	_baseLatch
	addandR int64
}

// Creates a new DownLatch.
//
// Preconditions:
// - initialValue > threshold;
// - initialValue - threshold <= MaxLatchDistance;
func NewDownLatch(initialValue, threshold int64) DownLatch {

	if initialValue <= threshold {

		panic(errDownLatchInitialValueMustBeGreaterThanThreshold)
	}

	if initialValue-threshold > MaxLatchDistance {

		panic(errLatchDistanceExceedsMaximum)
	}

	return DownLatch{
		_baseLatch: _baseLatch{
			value: initialValue - threshold,
		},
		addandR: threshold,
	}
}

func (l *DownLatch) Step() (flipped, isLatched bool, newCount int64) {

	flipped, isLatched, newCount = l._baseLatch.step()

	newCount += l.addandR

	return
}

// Obtains the current value of the latch, without changing its state.
func (l *DownLatch) Load() (isLatched bool, count int64) {

	isLatched, count = l._baseLatch.load()

	count += l.addandR

	return
}

// A unidirectional latch that counts up from an initial value to a higher
// threshold that may be operated safely by multiple concurrent goroutines.
type UpLatch struct {
	_baseLatch
	subandL int64
}

// Creates a new UpLatch.
//
// Preconditions:
// - initialValue < threshold;
// - threshold - initialValue <= MaxLatchDistance;
func NewUpLatch(initialValue, threshold int64) UpLatch {

	if initialValue >= threshold {

		panic(errUpLatchInitialValueMustBeLessThanThreshold)
	}

	if threshold-initialValue > MaxLatchDistance {

		panic(errLatchDistanceExceedsMaximum)
	}

	return UpLatch{
		_baseLatch: _baseLatch{
			value: threshold - initialValue,
		},
		subandL: threshold,
	}
}

func (l *UpLatch) Step() (flipped, isLatched bool, newCount int64) {

	var count int64

	flipped, isLatched, count = l._baseLatch.step()

	newCount = l.subandL - count

	return
}

// Obtains the current value of the latch, without changing its state.
func (l *UpLatch) Load() (isLatched bool, count int64) {

	var _count int64

	isLatched, _count = l._baseLatch.load()

	count = l.subandL - _count

	return
}
