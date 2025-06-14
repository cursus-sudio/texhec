package scopecleanup

import (
	"errors"
	"sync"
)

var (
	ErrAlreadyCleanedUp error = errors.New("ScopeCleanUp. Already cleaned up")
)

type CleanUpArgs struct {
	Error error
}

func NewCleanUpArgs(
	error error,
) CleanUpArgs {
	return CleanUpArgs{
		Error: error,
	}
}

type CleanUpListener func(args CleanUpArgs)

type ScopeCleanUp interface {
	// returns errors:
	// - ErrAlreadyCleanedUp
	Clean(Args CleanUpArgs) error

	// returns errors:
	// - ErrAlreadyCleanedUp
	AddCleanListener(listener CleanUpListener) error
}

type scopeCleanUp struct {
	listeners     []CleanUpListener
	alreadyCalled bool
	mutex         sync.Mutex
}

func (scopeCleanUp *scopeCleanUp) Clean(args CleanUpArgs) error {
	scopeCleanUp.mutex.Lock()
	defer scopeCleanUp.mutex.Unlock()
	if scopeCleanUp.alreadyCalled {
		return ErrAlreadyCleanedUp
	}

	scopeCleanUp.alreadyCalled = true

	for _, listener := range scopeCleanUp.listeners {
		listener(args)
	}

	return nil
}

func (scopeCleanUp *scopeCleanUp) AddCleanListener(listener CleanUpListener) error {
	scopeCleanUp.mutex.Lock()
	defer scopeCleanUp.mutex.Unlock()
	if scopeCleanUp.alreadyCalled {
		return ErrAlreadyCleanedUp
	}

	scopeCleanUp.alreadyCalled = true
	scopeCleanUp.listeners = append(scopeCleanUp.listeners, listener)

	return nil
}

func newScopeCleanUp() ScopeCleanUp {
	return &scopeCleanUp{
		listeners:     make([]CleanUpListener, 0),
		alreadyCalled: false,
	}
}
