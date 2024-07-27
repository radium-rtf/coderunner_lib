package container

import "errors"

var (
	ErrContainerHasAlreadyStarted = errors.New("container has already started")
	ErrContainerHasNotBeenStarted = errors.New("container has not been started")
	ErrContainerTimeout           = errors.New("container timeout")
)
