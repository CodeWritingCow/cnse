package tests

import (
	"encoding/json"
	"fmt"
	"testing"
	"voter-api/db"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

var (
	BASE_API = "http://localhost:1080"

	client = resty.New()
)

// Reset VoterList before each test
func TestMain(m *testing.M) {
	DeleteAllVotersResponse, error := client.R().Delete(BASE_API + "/voters")
	if DeleteAllVotersResponse.StatusCode() != 200 {
		fmt.Printf("error clearing database, %v", error)
	}

	AddSampleVotersResponse, error := client.R().Get(BASE_API + "/voters/add-sample-voters")
	if AddSampleVotersResponse.StatusCode() != 200 {
		fmt.Printf("error adding sample voters, %v", error)
	}

	code := m.Run()

	fmt.Println(code)
}

func Test_GetAllVoters(t *testing.T) {
	response, _ := client.R().Get(BASE_API + "/voters")
	voters := []db.Voter{}

	err := json.Unmarshal(response.Body(), &voters)

	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode())
	assert.Equal(t, 2, len(voters))
}

func Test_GetVoter(t *testing.T) {
	response, _ := client.R().Get(BASE_API + "/voters/1")
	voter := db.Voter{}

	err := json.Unmarshal(response.Body(), &voter)

	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode())
	assert.Equal(t, uint(1), voter.VoterId)
	assert.Equal(t, "Totoro", voter.Name)
}

func Test_AddVoter(t *testing.T) {
	response, _ := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(`{
			"voter_id": 3,
			"name": "Pikachu"
		}`).
		Post(BASE_API + "/voters/3")

	assert.Equal(t, 200, response.StatusCode())

	getResponse, _ := client.R().Get(BASE_API + "/voters")
	voters := []db.Voter{}

	err := json.Unmarshal(getResponse.Body(), &voters)

	assert.Nil(t, err)
	assert.Equal(t, 200, getResponse.StatusCode())
	assert.Equal(t, 3, len(voters))
}

func Test_DeleteVoter(t *testing.T) {
	DeleteAllVotersResponse, error := client.R().Delete(BASE_API + "/voters")
	if DeleteAllVotersResponse.StatusCode() != 200 {
		fmt.Printf("error clearing database, %v", error)
	}

	AddSampleVotersResponse, error := client.R().Get(BASE_API + "/voters/add-sample-voters")
	if AddSampleVotersResponse.StatusCode() != 200 {
		fmt.Printf("error adding sample voters, %v", error)
	}

	deleteResponse, _ := client.R().Delete(BASE_API + "/voters/1")

	assert.Equal(t, 200, deleteResponse.StatusCode())

	getResponse, _ := client.R().Get(BASE_API + "/voters")
	voters := []db.Voter{}

	err := json.Unmarshal(getResponse.Body(), &voters)

	assert.Nil(t, err)
	assert.Equal(t, 200, getResponse.StatusCode())
	assert.Equal(t, 1, len(voters))
}
