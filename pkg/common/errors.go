package common

import "errors"

// ErrRouteHandleWiring represent an error which suggests that a handler has been attached to an incompatible route.
var ErrRouteHandleWiring = errors.New("method wiring error")
