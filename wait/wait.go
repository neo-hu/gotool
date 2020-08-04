package wait

import (
	"context"
	"errors"
	"math/rand"
	"sync"
	"time"
)

var ForeverTestTimeout = time.Second * 30

var NeverStop <-chan struct{} = make(chan struct{})

type Group struct {
	wg sync.WaitGroup
}

func (g *Group) Wait() {
	g.wg.Wait()
}

func (g *Group) StartWithChannel(stopCh <-chan struct{}, f func(stopCh <-chan struct{})) {
	g.Start(func() {
		f(stopCh)
	})
}

func (g *Group) StartWithContext(ctx context.Context, f func(context.Context)) {
	g.Start(func() {
		f(ctx)
	})
}

// Start starts f in a new goroutine in the group.
func (g *Group) Start(f func()) {
	g.wg.Add(1)
	go func() {
		defer g.wg.Done()
		f()
	}()
}

func Forever(f func(), period time.Duration) {
	Until(f, period, NeverStop)
}

func Until(f func(), period time.Duration, stopCh <-chan struct{}) {
	JitterUntil(f, period, 0.0, true, stopCh)
}

func UntilWithContext(ctx context.Context, f func(context.Context), period time.Duration) {
	JitterUntilWithContext(ctx, f, period, 0.0, true)
}

func NonSlidingUntil(f func(), period time.Duration, stopCh <-chan struct{}) {
	JitterUntil(f, period, 0.0, false, stopCh)
}

func NonSlidingUntilWithContext(ctx context.Context, f func(context.Context), period time.Duration) {
	JitterUntilWithContext(ctx, f, period, 0.0, false)
}

func JitterUntil(f func(), period time.Duration, jitterFactor float64, sliding bool, stopCh <-chan struct{}) {
	var t *time.Timer
	var sawTimeout bool

	for {
		select {
		case <-stopCh:
			return
		default:
		}

		jitteredPeriod := period
		if jitterFactor > 0.0 {
			jitteredPeriod = Jitter(period, jitterFactor)
		}

		if !sliding {
			t = resetOrReuseTimer(t, jitteredPeriod, sawTimeout)
		}

		func() {
			f()
		}()

		if sliding {
			t = resetOrReuseTimer(t, jitteredPeriod, sawTimeout)
		}

		select {
		case <-stopCh:
			return
		case <-t.C:
			sawTimeout = true
		}
	}
}

func JitterUntilWithContext(ctx context.Context, f func(context.Context), period time.Duration, jitterFactor float64, sliding bool) {
	JitterUntil(func() { f(ctx) }, period, jitterFactor, sliding, ctx.Done())
}

func Jitter(duration time.Duration, maxFactor float64) time.Duration {
	if maxFactor <= 0.0 {
		maxFactor = 1.0
	}
	wait := duration + time.Duration(rand.Float64()*maxFactor*float64(duration))
	return wait
}

var ErrWaitTimeout = errors.New("timed out waiting for the condition")

type ConditionFunc func() (done bool, err error)

type Backoff struct {
	Duration time.Duration

	Factor float64
	Jitter float64
	Steps  int
	Cap    time.Duration
}

func (b *Backoff) Step() time.Duration {
	if b.Steps < 1 {
		if b.Jitter > 0 {
			return Jitter(b.Duration, b.Jitter)
		}
		return b.Duration
	}
	b.Steps--

	duration := b.Duration

	// calculate the next step
	if b.Factor != 0 {
		b.Duration = time.Duration(float64(b.Duration) * b.Factor)
		if b.Cap > 0 && b.Duration > b.Cap {
			b.Duration = b.Cap
			b.Steps = 0
		}
	}

	if b.Jitter > 0 {
		duration = Jitter(duration, b.Jitter)
	}
	return duration
}

func contextForChannel(parentCh <-chan struct{}) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		select {
		case <-parentCh:
			cancel()
		case <-ctx.Done():
		}
	}()
	return ctx, cancel
}

func ExponentialBackoff(backoff Backoff, condition ConditionFunc) error {
	for backoff.Steps > 0 {
		if ok, err := condition(); err != nil || ok {
			return err
		}
		if backoff.Steps == 1 {
			break
		}
		time.Sleep(backoff.Step())
	}
	return ErrWaitTimeout
}

func Poll(interval, timeout time.Duration, condition ConditionFunc) error {
	return pollInternal(poller(interval, timeout), condition)
}

func pollInternal(wait WaitFunc, condition ConditionFunc) error {
	done := make(chan struct{})
	defer close(done)
	return WaitFor(wait, condition, done)
}

func PollImmediate(interval, timeout time.Duration, condition ConditionFunc) error {
	return pollImmediateInternal(poller(interval, timeout), condition)
}

func pollImmediateInternal(wait WaitFunc, condition ConditionFunc) error {
	done, err := condition()
	if err != nil {
		return err
	}
	if done {
		return nil
	}
	return pollInternal(wait, condition)
}

func PollInfinite(interval time.Duration, condition ConditionFunc) error {
	done := make(chan struct{})
	defer close(done)
	return PollUntil(interval, condition, done)
}

func PollImmediateInfinite(interval time.Duration, condition ConditionFunc) error {
	done, err := condition()
	if err != nil {
		return err
	}
	if done {
		return nil
	}
	return PollInfinite(interval, condition)
}

func PollUntil(interval time.Duration, condition ConditionFunc, stopCh <-chan struct{}) error {
	ctx, cancel := contextForChannel(stopCh)
	defer cancel()
	return WaitFor(poller(interval, 0), condition, ctx.Done())
}

func PollImmediateUntil(interval time.Duration, condition ConditionFunc, stopCh <-chan struct{}) error {
	done, err := condition()
	if err != nil {
		return err
	}
	if done {
		return nil
	}
	select {
	case <-stopCh:
		return ErrWaitTimeout
	default:
		return PollUntil(interval, condition, stopCh)
	}
}

type WaitFunc func(done <-chan struct{}) <-chan struct{}

func WaitFor(wait WaitFunc, fn ConditionFunc, done <-chan struct{}) error {
	stopCh := make(chan struct{})
	defer close(stopCh)
	c := wait(stopCh)
	for {
		select {
		case _, open := <-c:
			ok, err := fn()
			if err != nil {
				return err
			}
			if ok {
				return nil
			}
			if !open {
				return ErrWaitTimeout
			}
		case <-done:
			return ErrWaitTimeout
		}
	}
}

func poller(interval, timeout time.Duration) WaitFunc {
	return WaitFunc(func(done <-chan struct{}) <-chan struct{} {
		ch := make(chan struct{})

		go func() {
			defer close(ch)

			tick := time.NewTicker(interval)
			defer tick.Stop()

			var after <-chan time.Time
			if timeout != 0 {
				timer := time.NewTimer(timeout)
				after = timer.C
				defer timer.Stop()
			}

			for {
				select {
				case <-tick.C:
					select {
					case ch <- struct{}{}:
					default:
					}
				case <-after:
					return
				case <-done:
					return
				}
			}
		}()

		return ch
	})
}

func resetOrReuseTimer(t *time.Timer, d time.Duration, sawTimeout bool) *time.Timer {
	if t == nil {
		return time.NewTimer(d)
	}
	if !t.Stop() && !sawTimeout {
		<-t.C
	}
	t.Reset(d)
	return t
}
