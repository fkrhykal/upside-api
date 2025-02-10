package exception

type AuthenticationError struct{}

func (e *AuthenticationError) Error() string {
	return "authentication error"
}

type AuthorizationError struct{}

func (e *AuthorizationError) Error() string {
	return "authorization error"
}

var ErrAuthentication = &AuthenticationError{}
var ErrAuthorization = &AuthorizationError{}
