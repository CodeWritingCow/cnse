package tests

import (
	"encoding/json"
	"fmt"
	"testing"
	"voter-api/voter"

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
	myResponse := voter.VoterList{}

	err := json.Unmarshal(response.Body(), &myResponse)

	// fmt.Println(myResponse)

	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode())
	assert.Equal(t, 2, len(myResponse.Voters))
}
