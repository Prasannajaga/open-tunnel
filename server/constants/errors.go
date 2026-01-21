package constants

const (
	ErrControlAccept  = "control listener accept error"
	ErrExternalAccept = "error accepting external connection"
	ErrControlDown    = "external request received but control channel is down"
	ErrControlWrite   = "control write failed"
	ErrDataConnection = "data connection failed"
	ErrInvalidCreds   = "invalid username or password"
	ErrTokenExpired   = "token has expired"
	ErrTokenInvalid   = "invalid token"
	ErrMissingToken   = "missing authorization token"
	ErrUserNotFound   = "user not found"
)
