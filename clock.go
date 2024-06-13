package quicktick

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

type QuickTick struct {
	multiplier      float64
	tickerInterval  time.Duration
	startTime       time.Time
	startRealTime   time.Time
	currentDuration int64
	done            chan struct{}
	mu              sync.Mutex
	once            sync.Once
}

// Create a new QuickTick clock with the given multiplier.
//
// Uses the current time and updates each Millisecond. For more customizability use the NewCustom function.
func New(multiplier float64) *QuickTick {
	startTime := time.Now()

	ac := &QuickTick{
		multiplier:      multiplier,
		tickerInterval:  time.Millisecond,
		startTime:       startTime,
		startRealTime:   startTime,
		currentDuration: 0,
		done:            make(chan struct{}),
	}
	go ac.run()
	return ac
}

// NewCustom creates a new QuickTick clock starting at the given startTime, with the specified multiplier and updateInterval.
//
// Parameters:
//   - startTime: The initial time from which the accelerated time will be calculated.
//   - multiplier: The rate at which the accelerated time progresses relative to real time. For example, a multiplier of 2.0 means the accelerated time runs twice as fast as real time.
//   - updateInterval: The interval at which the clock will update. This allows for customization of how frequently the clock recalculates the accelerated time.
//
// Example Usage:
//
//	startTime := time.Now()
//	multiplier := 1.5 // Time runs 1.5 times faster
//	updateInterval := 500 * time.Millisecond // Update the clock every 500 milliseconds
//	clock := quicktick.NewCustom(startTime, multiplier, updateInterval)
//
// Use this function if you need to customize the clock's update frequency and starting time.
func NewCustom(startTime time.Time, multiplier float64, updateInterval time.Duration) *QuickTick {
	ac := &QuickTick{
		multiplier:      multiplier,
		tickerInterval:  updateInterval,
		startTime:       startTime,
		startRealTime:   time.Now(),
		currentDuration: 0,
		done:            make(chan struct{}),
	}
	go ac.run()
	return ac
}

// Create a new QuickTick clock with the given context and multiplier.
//
// Uses the current time and updates each Millisecond. For more customizability use the NewCustomCtx function.
func NewCtx(ctx context.Context, multiplier float64) *QuickTick {
	startTime := time.Now()

	ac := &QuickTick{
		multiplier:      multiplier,
		tickerInterval:  time.Millisecond,
		startTime:       startTime,
		startRealTime:   startTime,
		currentDuration: 0,
		done:            make(chan struct{}),
	}
	go ac.runWithContext(ctx)
	return ac
}

// NewCustomCtx creates a new QuickTick clock starting at the given startTime, with the specified context, multiplier, and updateInterval.
//
// Parameters:
//   - ctx: The context to control the lifecycle of the clock. The clock will stop updating when the context is done.
//   - startTime: The initial time from which the accelerated time will be calculated.
//   - multiplier: The rate at which the accelerated time progresses relative to real time. For example, a multiplier of 2.0 means the accelerated time runs twice as fast as real time.
//   - updateInterval: The interval at which the clock will update. This allows for customization of how frequently the clock recalculates the accelerated time.
//
// Example Usage:
//
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//	startTime := time.Now()
//	multiplier := 1.5 // Time runs 1.5 times faster
//	updateInterval := 500 * time.Millisecond // Update the clock every 500 milliseconds
//	clock := quicktick.NewCustomCtx(ctx, startTime, multiplier, updateInterval)
//
// Use this function if you need to customize the clock's update frequency and starting time, and also want to control the clock's lifecycle with a context.
func NewCustomCtx(ctx context.Context, startTime time.Time, multiplier float64, updateInterval time.Duration) *QuickTick {
	ac := &QuickTick{
		multiplier:      multiplier,
		tickerInterval:  updateInterval,
		startTime:       startTime,
		startRealTime:   time.Now(),
		currentDuration: 0,
		done:            make(chan struct{}),
	}
	go ac.runWithContext(ctx)
	return ac
}

// Start clock
func (ac *QuickTick) run() {
	ticker := time.NewTicker(time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ac.updateClock()
		case <-ac.done:
			return
		}
	}
}

// Start clock with context
func (ac *QuickTick) runWithContext(ctx context.Context) {
	ticker := time.NewTicker(time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ac.updateClock()
		case <-ctx.Done():
			ac.Stop()
			return
		case <-ac.done:
			return
		}
	}
}

func (ac *QuickTick) updateClock() {
	ac.mu.Lock()
	elapsedRealTime := time.Since(ac.startRealTime)
	ac.mu.Unlock()
	acceleratedElapsedTime := elapsedRealTime.Seconds() * ac.multiplier
	atomic.StoreInt64(&ac.currentDuration, int64(acceleratedElapsedTime*float64(time.Second)))
}

// Now returns the current accelerated time.
func (ac *QuickTick) Now() time.Time {
	currentDuration := atomic.LoadInt64(&ac.currentDuration)
	accumulatedDuration := time.Duration(currentDuration)
	acceleratedTime := ac.startTime.Add(accumulatedDuration)
	return acceleratedTime
}

// Stop stops the QuickTick clock.
func (ac *QuickTick) Stop() {
	ac.once.Do(func() {
		close(ac.done)
	})
}

// Reset resets the QuickTick clock.
func (ac *QuickTick) Reset() {
	ac.mu.Lock()
	defer ac.mu.Unlock()
	ac.startRealTime = time.Now()
	atomic.StoreInt64(&ac.currentDuration, 0)
}
