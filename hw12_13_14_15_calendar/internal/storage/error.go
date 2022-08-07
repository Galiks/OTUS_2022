package storage

import "errors"

var (
	ErrDateBusy     = errors.New("this time is busy by another event")
	ErrUnknownEvent = errors.New("unknown event")
	ErrUnknownUser  = errors.New("unknown user")
)
