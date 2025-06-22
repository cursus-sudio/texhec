package scopes

import "sync"

type UserSessionEndArgs struct {
}

func NewUserSessionEndArgs() UserSessionEndArgs {
	return UserSessionEndArgs{}
}

type UserSessionEndListener func(args UserSessionEndArgs)

type UserSessionEnd interface {
	// returns errors:
	// - ErrAlreadyCleanedUp
	Clean(Args UserSessionEndArgs) error

	// returns errors:
	// - ErrAlreadyCleanedUp
	AddCleanListener(listener UserSessionEndListener) error
}

type userSessionEnd struct {
	listeners     []UserSessionEndListener
	alreadyCalled bool
	mutex         sync.Mutex
}

func (sessionEnd *userSessionEnd) Clean(args UserSessionEndArgs) error {
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

func (sessionEnd *userSessionEnd) AddCleanListener(listener UserSessionEndListener) error {
	sessionEnd.mutex.Lock()
	defer sessionEnd.mutex.Unlock()
	if sessionEnd.alreadyCalled {
		return ErrAlreadyCleanedUp
	}

	sessionEnd.listeners = append(sessionEnd.listeners, listener)

	return nil
}

func newSessionEnd() UserSessionEnd {
	return &userSessionEnd{
		listeners:     []UserSessionEndListener{},
		alreadyCalled: false,
		mutex:         sync.Mutex{},
	}
}
