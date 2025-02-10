package exception

type AlreadyJoinSideError struct{}

func (al *AlreadyJoinSideError) Error() string {
	return "already joined to side"
}

var (
	ErrAlreadyJoinSide = &AlreadyJoinSideError{}
)
