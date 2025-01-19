package exception

type AuthenticationError struct {
}

func (e *AuthenticationError) Error() string {
	return "authentication error"
}

var ErrAuthentication = &AuthenticationError{}
