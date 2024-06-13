package quicktick

import (
	"context"
	"sync"
	"testing"
	"time"
)

const (
	multiplier      = 1.5
	startTimeOffset = -1 * time.Hour
	updateInterval  = 500 * time.Millisecond
	tolerance       = 0.1
)

func TestNew(t *testing.T) {
	clock := New(multiplier)
	defer clock.Stop()

	time.Sleep(2000 * time.Millisecond)

	acceleratedTime := clock.Now()
	elapsedRealTime := time.Since(clock.startRealTime)

	if acceleratedTime.Before(clock.startTime) {
		t.Errorf("Accelerated time should not be before start time")
	}

	acceleratedElapsed := acceleratedTime.Sub(clock.startTime).Seconds()
	expectedElapsed := elapsedRealTime.Seconds() * multiplier

	if acceleratedElapsed < expectedElapsed-tolerance || acceleratedElapsed > expectedElapsed+tolerance {
		t.Errorf("Expected accelerated time to be approximately %v, but got %v", expectedElapsed, acceleratedElapsed)
	}
}

func TestNewCustom(t *testing.T) {
	startTime := time.Now().Add(startTimeOffset)
	clock := NewCustom(startTime, multiplier, updateInterval)
	defer clock.Stop()

	time.Sleep(2000 * time.Millisecond)

	acceleratedTime := clock.Now()
	if acceleratedTime.Before(startTime) {
		t.Errorf("Accelerated time should not be before start time")
	}

	elapsedRealTime := time.Since(clock.startRealTime).Seconds()
	expectedElapsed := elapsedRealTime * multiplier
	acceleratedElapsed := acceleratedTime.Sub(startTime).Seconds()

	if acceleratedElapsed < expectedElapsed-tolerance || acceleratedElapsed > expectedElapsed+tolerance {
		t.Errorf("Expected accelerated time to be approximately %v, but got %v", expectedElapsed, acceleratedElapsed)
	}
}

func TestNewCtx(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	clock := NewCtx(ctx, multiplier)
	defer clock.Stop()

	time.Sleep(2000 * time.Millisecond)

	acceleratedTime := clock.Now()
	elapsedRealTime := time.Since(clock.startRealTime)

	if acceleratedTime.Before(clock.startTime) {
		t.Errorf("Accelerated time should not be before start time")
	}

	acceleratedElapsed := acceleratedTime.Sub(clock.startTime).Seconds()
	expectedElapsed := elapsedRealTime.Seconds() * multiplier

	if acceleratedElapsed < expectedElapsed-tolerance || acceleratedElapsed > expectedElapsed+tolerance {
		t.Errorf("Expected accelerated time to be approximately %v, but got %v", expectedElapsed, acceleratedElapsed)
	}
}

func TestNewCustomCtx(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	startTime := time.Now().Add(startTimeOffset)
	clock := NewCustomCtx(ctx, startTime, multiplier, updateInterval)
	defer clock.Stop()

	time.Sleep(2000 * time.Millisecond)

	acceleratedTime := clock.Now()
	if acceleratedTime.Before(startTime) {
		t.Errorf("Accelerated time should not be before start time")
	}

	elapsedRealTime := time.Since(clock.startRealTime).Seconds()
	expectedElapsed := elapsedRealTime * multiplier
	acceleratedElapsed := acceleratedTime.Sub(startTime).Seconds()

	if acceleratedElapsed < expectedElapsed-tolerance || acceleratedElapsed > expectedElapsed+tolerance {
		t.Errorf("Expected accelerated time to be approximately %v, but got %v", expectedElapsed, acceleratedElapsed)
	}
}

func TestConcurrency(t *testing.T) {
	clock := New(multiplier)
	defer clock.Stop()

	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 1000; i++ {
				clock.Now()
				time.Sleep(1 * time.Millisecond)
			}
		}()
	}

	wg.Wait()
}

func TestConcurrency2(t *testing.T) {
	clock := New(multiplier)
	defer clock.Stop()

	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer clock.Stop()
			for i := 0; i < 1000; i++ {
				clock.Now()
				clock.Reset()
				time.Sleep(1 * time.Millisecond)
			}
		}()
	}

	wg.Wait()
}
