package exception

type AlreadyVoteError struct {
}

func (ar *AlreadyVoteError) Error() string {
	return "already vote error"
}

var (
	ErrAlreadyVote = &AlreadyVoteError{}
)
