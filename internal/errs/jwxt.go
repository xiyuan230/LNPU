package errs

import "errors"

var (
	ErrCookieExpire    = errors.New("cookie expire")
	ErrUserEmpty       = errors.New("username or password is empty")
	ErrPasswordWrong   = errors.New("password is wrong")
	ErrJwxtLoginFailed = errors.New("login jwxt failed")
)
