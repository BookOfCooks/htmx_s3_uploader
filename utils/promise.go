package utils

// Promise is meant to start immediate execution of a function, but not wait for it to finish
// Call Wait() to wait for the `compute` function to finish or get the cached result if it has already finished
// In addition, panics are swallowed by Promise, and will be rethrown after calling Wait()
type Promise[T any] struct {
	finished   chan struct{}
	compute    func() T
	cacheValue T
	panicked   bool
	cachePanic any
}

func NewPromise[T any](compute func() T) *Promise[T] {
	p := &Promise[T]{
		finished: make(chan struct{}),
		compute:  compute,
	}

	go func() {
		defer func() {
			if e := recover(); e != nil {
				p.cachePanic = e
				p.panicked = true
				close(p.finished)
			}
		}()

		t := compute()
		p.cacheValue = t
		close(p.finished)
	}()

	return p
}

// Returns the result of the `compute` function, or the cached result if it has already finished
// Will rethrow any panics that comes from `compute`
func (p *Promise[T]) Wait() T {
	<-p.finished
	if p.panicked {
		panic(p.cachePanic)
	}
	return p.cacheValue
}
