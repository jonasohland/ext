package errors

import "sync/atomic"

// Simple structure that contains an error that can be atomically set and cleared
type LastError struct {
	v atomic.Value
}

type lastErrorWrapper struct {
	v error
}

// Set the error atomically, future calls to Get() will return this error
func (l *LastError) Set(err error) error {
	v := l.v.Swap(lastErrorWrapper{err})
	wrapper, ok := v.(lastErrorWrapper)
	if !ok {
		return nil
	}

	return wrapper.v
}

// Clear the error, future calls to Get() will return nil
func (l *LastError) Clear() error {
	v := l.v.Swap(lastErrorWrapper{})
	wrapper, ok := v.(lastErrorWrapper)
	if !ok {
		return nil
	}

	return wrapper.v
}

// Get the error that was last stored with Store(), or nil if Clear() was called
func (l *LastError) Get() error {
	v := l.v.Load()
	wrapper, ok := v.(lastErrorWrapper)
	if !ok {
		return nil
	}

	return wrapper.v
}
