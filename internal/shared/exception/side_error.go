package exception

type AlreadyMemberError struct{}

func (am *AlreadyMemberError) Error() string {
	return "already joined to side"
}

type RequiredMembershipError struct{}

func (rm *RequiredMembershipError) Error() string {
	return "required membership"
}

var (
	ErrAlreadyMember      = &AlreadyMemberError{}
	ErrRequiredMembership = &RequiredMembershipError{}
)
