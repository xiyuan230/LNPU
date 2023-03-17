package errs

import "errors"

var (
	ErrParamMiss     = errors.New("missing param")
	ErrPathIllegal   = errors.New("path is illegal")
	ErrWxLoginFailed = errors.New("wechat login failed")
)
