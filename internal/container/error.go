package container

import "errors"

var (
	ErrContainerHasAlreadyStarted = errors.New("container has already started")
)
