package handlers

import "errors"

var ErrNoQuery = errors.New("no query parameter")
var ErrInvalidKey = errors.New("invalid key")
