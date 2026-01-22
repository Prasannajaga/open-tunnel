package constants

const (
	ErrServerUnreachable  = "tunnel server at %s is not responding or unreachable"
	ErrInvalidPort        = "invalid port argument: %s"
	ErrLocalServiceDown   = "could not reach local service at %s (port %d)"
	ErrDataChannelFailed  = "failed to open data channel to server"
	ErrControlWriteFailed = "control write failed"
	ErrUnknownCommand     = "unknown command: %s"
	ErrCredentialsExpired = "credentials expired. Please login again"
)
