package api

import (
	"log"
	"net/http"
	"strconv"
	"time"
	"voter-api/db"

	"github.com/gin-gonic/gin"
)

type VoterAPI struct {
	db        *db.ToDo
	voterList db.VoterList
}

func New() (*VoterAPI, error) {
	dbHandler, err := db.New()
	if err != nil {
		return nil, err
	}

	return &VoterAPI{db: dbHandler}, nil
}

func (td *VoterAPI) GetAllVoters(c *gin.Context) {
	voters, err := td.db.GetAllVoters()
	if err != nil {
		log.Println("Error Getting All Items: ", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if voters == nil {
		voters = make([]db.Voter, 0)
	}

	c.JSON(http.StatusOK, voters)
}

func (td *VoterAPI) GetVoter(c *gin.Context) {
	idS := c.Param("id")
	id64, err := strconv.ParseInt(idS, 10, 32)

	if err != nil {
		log.Println("Error converting id to int64: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	voter, err := td.db.GetVoter(int(id64))

	if err != nil {
		log.Println("Item not found: ", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, voter)
}

func (td *VoterAPI) AddVoter(c *gin.Context) {
	var newVoter db.Voter
	newVoter.VoteHistory = []db.VoterHistory{}

	if err := c.ShouldBindJSON(&newVoter); err != nil {
		log.Println("Error binding JSON: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := td.db.AddVoter(newVoter); err != nil {
		log.Println("Error adding item: ", err)
		c.AbortWithStatus(http.StatusConflict)
		return
	}

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

	if err := td.db.DeleteVoter(int(id64)); err != nil {
		log.Println("Error deleting voter: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Voter successfully deleted"})
}

func (td *VoterAPI) DeleteAllVoters(c *gin.Context) {
	if err := td.db.DeleteAll(); err != nil {
		log.Println("Error deleting all voters: ", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "All voters successfully deleted"})
}

func (td *VoterAPI) GetVoterPolls(c *gin.Context) {
	idS := c.Param("id")
	id64, err := strconv.ParseInt(idS, 10, 32)

	if err != nil {
		log.Println("Error converting id to int64: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	voterHistory, err := td.db.GetVoterPolls(int(id64))
	if err != nil {
		log.Println("Item not found: ", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, voterHistory)
}

func (td *VoterAPI) GetVoterPoll(c *gin.Context) {
	idS := c.Param("id")
	id64, err := strconv.ParseInt(idS, 10, 32)

	if err != nil {
		log.Println("Error converting id to int64: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	voter, err := td.db.GetVoter(int(id64))
	if err != nil {
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

	voter, err := td.db.GetVoter(int(voterId64))
	if err != nil {
		log.Println("Voter not found: ", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	for _, poll := range voter.VoteHistory {
		if int64(poll.PollId) == pollId64 {
			log.Println("Voter poll already exists")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}

	var currentTime = time.Now()

	newVoterPoll := db.NewVoterHistory(uint(pollId64), currentTime)

	if err := c.ShouldBindJSON(&newVoterPoll); err != nil {
		log.Println("Error binding JSON: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	td.db.AddVoterPollHistory(int(voterId64), int(pollId64), currentTime)

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

	voter, err := td.db.GetVoter(int(voterId64))
	if err != nil {
		log.Println("Voter not found")
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	for i, poll := range voter.VoteHistory {
		if int64(poll.PollId) == pollId64 {
			voter.VoteHistory = append(voter.VoteHistory[:i], voter.VoteHistory[i+1:]...)
			td.db.DeleteVoterPoll(int(voterId64), int(pollId64))
			c.JSON(http.StatusOK, gin.H{"message": "Voter poll successfully deleted"})
			return
		}
	}

	log.Println("Voter poll not found")
	c.AbortWithStatus(http.StatusNotFound)
}

func (td *VoterAPI) AddSampleVoters(c *gin.Context) {
	td.voterList.Voters[0] = db.Voter{
		VoterId: 0,
		Name:    "Moo Moo",
		VoteHistory: []db.VoterHistory{
			{
				PollId:   0,
				VoteDate: time.Now(),
			},
		},
	}

	td.voterList.Voters[1] = db.Voter{
		VoterId: 1,
		Name:    "Totoro",
		VoteHistory: []db.VoterHistory{
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
