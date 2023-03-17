package errs

import "errors"

var (
	ErrCookieExpire    = errors.New("cookie expire")
	ErrUserIllegal     = errors.New("username or password is illegal")
	ErrPasswordWrong   = errors.New("password is wrong")
	ErrJwxtLoginFailed = errors.New("login jwxt failed")
)
