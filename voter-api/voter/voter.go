package voter

import (
	// "encoding/json"
	// "errors"
	// "fmt"
	"time"
)

type VoterHistory struct {
	PollId   uint      `json:"poll_id"`
	VoteDate time.Time `json:"vote_date"`
}

type Voter struct {
	VoterId     uint           `json:"voter_id"`
	Name        string         `json:"name"`
	VoteHistory []VoterHistory `json:"voter_history"`
}
type VoterList struct {
	Voters map[uint]Voter `json:"voters"` //A map of VoterIDs as keys and Voter structs as values
}

// constructor for VoterList struct
func NewVoter(id uint, name string) *Voter {
	return &Voter{
		Name:        name,
		VoteHistory: []VoterHistory{},
	}
}
