package api

import (
	"log"
	"net/http"
	"strconv"
	"time"
	"voter-api/db"
	"voter-api/voter"

	"github.com/gin-gonic/gin"
)

// TODO: Refactor TodoAPI to replace ToDo API code with Voter API code
type VoterAPI struct {
	db        *db.ToDo
	voterList voter.VoterList
}

func New() (*VoterAPI, error) {
	dbHandler, err := db.New()
	if err != nil {
		return nil, err
	}

	return &VoterAPI{db: dbHandler}, nil
}

func (td *VoterAPI) GetVoterList(c *gin.Context) {
	if td.voterList.Voters == nil {
		td.voterList.Voters = make(map[uint]voter.Voter)

		// TODO: Delete code for adding sample voter
		td.AddSampleVoters(c)
	}

	c.JSON(http.StatusOK, td.voterList)
}

func (td *VoterAPI) GetVoter(c *gin.Context) {
	idS := c.Param("id")
	id64, err := strconv.ParseInt(idS, 10, 32)

	if err != nil {
		log.Println("Error converting id to int64: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	voter, ok := td.voterList.Voters[uint(id64)]
	if !ok {
		log.Println("Item not found")
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, voter)
}

func (td *VoterAPI) AddVoter(c *gin.Context) {
	var newVoter voter.Voter
	newVoter.VoteHistory = []voter.VoterHistory{}

	if err := c.ShouldBindJSON(&newVoter); err != nil {
		log.Println("Error binding JSON: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	_, doesVoterExist := td.voterList.Voters[newVoter.VoterId]
	if doesVoterExist {
		log.Println("Voter already exists")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	td.voterList.Voters[newVoter.VoterId] = newVoter

	c.JSON(http.StatusOK, newVoter)
}

func (td *VoterAPI) DeleteVoter(c *gin.Context) {
	id := c.Param("id")
	id64, err := strconv.ParseInt(id, 10, 32)

	if err != nil {
		log.Println("Error converting id to int64: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	_, doesVoterExist := td.voterList.Voters[uint(id64)]
	if !doesVoterExist {
		log.Println("Voter not found")
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	delete(td.voterList.Voters, uint(id64))

	c.JSON(http.StatusOK, gin.H{"message": "Voter successfully deleted"})
}

func (td *VoterAPI) ListVoterPolls(c *gin.Context) {
	idS := c.Param("id")
	id64, err := strconv.ParseInt(idS, 10, 32)

	if err != nil {
		log.Println("Error converting id to int64: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	voter, ok := td.voterList.Voters[uint(id64)]
	if !ok {
		log.Println("Item not found")
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, voter.VoteHistory)
}

func (td *VoterAPI) GetVoterPoll(c *gin.Context) {
	idS := c.Param("id")
	id64, err := strconv.ParseInt(idS, 10, 32)

	if err != nil {
		log.Println("Error converting id to int64: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	voter, ok := td.voterList.Voters[uint(id64)]
	if !ok {
		log.Println("Item not found")
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	pollId := c.Param("pollid")
	pollId64, pollErr := strconv.ParseInt(pollId, 10, 32)

	if pollErr != nil {
		log.Println("Error converting pollId to int64: ", pollErr)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	for _, poll := range voter.VoteHistory {
		if int64(poll.PollId) == pollId64 {
			c.JSON(http.StatusOK, poll)
			return
		}
	}

	log.Println("Item not found")
	c.AbortWithStatus(http.StatusNotFound)
}

func (td *VoterAPI) AddVoterPoll(c *gin.Context) {
	voterIdS := c.Param("id")
	voterId64, err := strconv.ParseInt(voterIdS, 10, 32)

	if err != nil {
		log.Println("Error converting voterId to int64: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	pollId := c.Param("pollid")
	pollId64, pollErr := strconv.ParseInt(pollId, 10, 32)

	if pollErr != nil {
		log.Println("Error converting pollId to int64: ", pollErr)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user, ok := td.voterList.Voters[uint(voterId64)]
	if !ok {
		log.Println("Voter not found")
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	for _, poll := range user.VoteHistory {
		if int64(poll.PollId) == pollId64 {
			log.Println("Voter poll already exists")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}

	newVoterPoll := voter.NewVoterHistory(uint(pollId64), time.Now())

	if err := c.ShouldBindJSON(&newVoterPoll); err != nil {
		log.Println("Error binding JSON: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user.VoteHistory = append(user.VoteHistory, *newVoterPoll)
	td.voterList.Voters[uint(voterId64)] = user

	c.JSON(http.StatusOK, newVoterPoll)
}

func (td *VoterAPI) DeleteVoterPoll(c *gin.Context) {
	voterId := c.Param("id")
	voterId64, err := strconv.ParseInt(voterId, 10, 32)

	if err != nil {
		log.Println("Error converting voterId to int64: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	pollId := c.Param("pollid")
	pollId64, pollErr := strconv.ParseInt(pollId, 10, 32)

	if pollErr != nil {
		log.Println("Error converting pollId to int64: ", pollErr)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user, ok := td.voterList.Voters[uint(voterId64)]
	if !ok {
		log.Println("Voter not found")
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	for i, poll := range user.VoteHistory {
		if int64(poll.PollId) == pollId64 {
			user.VoteHistory = append(user.VoteHistory[:i], user.VoteHistory[i+1:]...)
			td.voterList.Voters[uint(voterId64)] = user
			c.JSON(http.StatusOK, gin.H{"message": "Voter poll successfully deleted"})
			return
		}
	}

	log.Println("Voter poll not found")
	c.AbortWithStatus(http.StatusNotFound)
}

// TODO: Delete AddSampleVoter
func (td *VoterAPI) AddSampleVoters(c *gin.Context) {
	td.voterList.Voters[0] = voter.Voter{
		VoterId: 0,
		Name:    "Moo Moo",
		VoteHistory: []voter.VoterHistory{
			{
				PollId:   0,
				VoteDate: time.Now(),
			},
		},
	}

	td.voterList.Voters[1] = voter.Voter{
		VoterId: 1,
		Name:    "Totoro",
		VoteHistory: []voter.VoterHistory{
			{
				PollId:   0,
				VoteDate: time.Now(),
			},
		},
	}
}

// TODO: Remove unused boilerplate code

// implementation for PUT /todo
// func (td *ToDoAPI) UpdateToDo(c *gin.Context) {
// 	var todoItem db.ToDoItem
// 	if err := c.ShouldBindJSON(&todoItem); err != nil {
// 		log.Println("Error binding JSON: ", err)
// 		c.AbortWithStatus(http.StatusBadRequest)
// 		return
// 	}

// 	if err := td.db.UpdateItem(todoItem); err != nil {
// 		log.Println("Error updating item: ", err)
// 		c.AbortWithStatus(http.StatusInternalServerError)
// 		return
// 	}

// 	c.JSON(http.StatusOK, todoItem)
// }

// Deletes all todos
// func (td *ToDoAPI) DeleteAllToDo(c *gin.Context) {

// 	if err := td.db.DeleteAll(); err != nil {
// 		log.Println("Error deleting all items: ", err)
// 		c.AbortWithStatus(http.StatusInternalServerError)
// 		return
// 	}

// 	c.Status(http.StatusOK)
// }

func (td *VoterAPI) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK,
		gin.H{
			"status":             "ok",
			"version":            "1.0.0",
			"uptime":             100,
			"users_processed":    1000,
			"errors_encountered": 10,
		})
}
