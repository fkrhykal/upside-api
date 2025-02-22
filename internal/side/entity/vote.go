package entity

import "github.com/google/uuid"

type Voter struct {
	ID uuid.UUID
}

type VoteKind int8

const (
	UpVote   VoteKind = 1
	DownVote VoteKind = -1
)

type Vote struct {
	ID    uuid.UUID
	Voter *Voter
	Post  *Post
	Kind  VoteKind
}
