package voter

import (
	// "encoding/json"
	// "errors"
	// "fmt"
	"time"
)

type VoterHistory struct {
	PollId   uint
	VoteDate time.Time
}

type Voter struct {
	VoterId     uint
	Name        string
	VoteHistory []VoterHistory
}
type VoterList struct {
	Voters map[uint]Voter //A map of VoterIDs as keys and Voter structs as values
}

// constructor for VoterList struct
func NewVoter(id uint, name string) *Voter {
	return &Voter{
		Name:        name,
		VoteHistory: []VoterHistory{},
	}
}
