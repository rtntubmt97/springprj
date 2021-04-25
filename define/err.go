package define

import "errors"

// Some common errors used in this project

var ErrWrongInitBytes = errors.New("WrongInitBytes")
var ErrWrongCmd = errors.New("WrongCmd")
var ErrFailGreeting = errors.New("FailGreeting")
