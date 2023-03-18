package errs

import "errors"

var (
	ErrParamMiss    = errors.New("missing param")
	ErrPathIllegal  = errors.New("path is illegal")
	ErrTokenIllegal = errors.New("token is illegal")
)
