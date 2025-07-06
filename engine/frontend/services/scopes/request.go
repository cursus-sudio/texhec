package scopes

import "sync"

type RequestEndArgs struct {
	Error error
}

func NewRequestEndArgs(err error) RequestEndArgs {
	return RequestEndArgs{
		Error: err,
	}
}

type RequestEndListener func(args RequestEndArgs)

type RequestService interface {
	// returns errors:
	// - ErrAlreadyCleanedUp
	Clean(Args RequestEndArgs) error

	// returns errors:
	// - ErrAlreadyCleanedUp
	AddCleanListener(listener RequestEndListener) error
}

type requestEnd struct {
	listeners     []RequestEndListener
	alreadyCalled bool
	mutex         sync.Mutex
}

func (sessionEnd *requestEnd) Clean(args RequestEndArgs) error {
	sessionEnd.mutex.Lock()
	defer sessionEnd.mutex.Unlock()
	if sessionEnd.alreadyCalled {
		return ErrAlreadyCleanedUp
	}

	sessionEnd.alreadyCalled = true

	for _, listener := range sessionEnd.listeners {
		listener(args)
	}

	return nil
}

func (sessionEnd *requestEnd) AddCleanListener(listener RequestEndListener) error {
	sessionEnd.mutex.Lock()
	defer sessionEnd.mutex.Unlock()
	if sessionEnd.alreadyCalled {
		return ErrAlreadyCleanedUp
	}

	sessionEnd.listeners = append(sessionEnd.listeners, listener)

	return nil
}

func newRequestService() RequestService {
	return &requestEnd{
		listeners:     []RequestEndListener{},
		alreadyCalled: false,
		mutex:         sync.Mutex{},
	}
}
