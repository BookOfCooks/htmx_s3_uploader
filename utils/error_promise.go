package utils

// ErrorPromise is meant to start immediate execution of a function, but not wait for it to finish
// Call Wait() to wait for the `compute` function to finish or get the cached result if it has already finished
// In addition, panics are swallowed by ErrorPromise, and will be rethrown after calling Wait()
type ErrorPromise[T any] struct {
	finished   chan struct{}
	compute    func() (T, error)
	cacheValue T
	cacheErr   error
	panicked   bool
	cachePanic any
}

func NewErrorPromise[T any](compute func() (T, error)) *ErrorPromise[T] {
	p := &ErrorPromise[T]{
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

		t, err := compute()
		p.cacheValue = t
		p.cacheErr = err
		close(p.finished)
	}()

	return p
}

// Returns the result of the `compute` function, or the cached result if it has already finished
// Will rethrow any panics that comes from `compute`
func (p *ErrorPromise[T]) Wait() (T, error) {
	<-p.finished
	if p.panicked {
		panic(p.cachePanic)
	}
	return p.cacheValue, p.cacheErr
}
