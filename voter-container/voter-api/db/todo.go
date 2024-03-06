package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/nitishm/go-rejson/v4"
	"github.com/redis/go-redis/v9"
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

// Constructor for VoterList struct
func NewVoter(id uint, name string) *Voter {
	return &Voter{
		Name:        name,
		VoteHistory: []VoterHistory{},
	}
}

// Constructor for VoterHistory struct
func NewVoterHistory(pollId uint, voteDate time.Time) *VoterHistory {
	return &VoterHistory{
		PollId:   pollId,
		VoteDate: voteDate,
	}
}

const (
	RedisNilError        = "redis: nil"
	RedisDefaultLocation = "0.0.0.0:6379"
	RedisKeyPrefix       = "voter:"
)

type cache struct {
	cacheClient *redis.Client
	jsonHelper  *rejson.Handler
	context     context.Context
}

// DbMap is a type alias for a map of ToDoItems.  The key
// will be the ToDoItem.Id and the value will be the ToDoItem
type DbMap map[int]Voter

// ToDo is the struct that represents the main object of our
// todo app.  It contains a map of ToDoItems and the name of
// the file that is used to store the items.
//
// This is just a mock, so we will only be managing an in memory
// map
type ToDo struct {
	toDoMap DbMap
	//more things would be included in a real implementation

	cache
}

func New() (*ToDo, error) {
	redisUrl := os.Getenv("REDIS_URL")
	if redisUrl == "" {
		redisUrl = RedisDefaultLocation
	}
	return NewWithCacheInstance(redisUrl)
}

func NewWithCacheInstance(location string) (*ToDo, error) {
	//Connect to redis.  Other options can be provided, but the
	//defaults are OK
	client := redis.NewClient(&redis.Options{
		Addr: location,
	})

	//We use this context to coordinate betwen our go code and
	//the redis operaitons
	ctx := context.Background()

	//This is the reccomended way to ensure that our redis connection
	//is working
	err := client.Ping(ctx).Err()
	if err != nil {
		log.Println("Error connecting to redis" + err.Error() + "cache might not be available, continuing...")
	}

	//By default, redis manages keys and values, where the values
	//are either strings, sets, maps, etc.  Redis has an extension
	//module called ReJSON that allows us to store JSON objects
	//however, we need a companion library in order to work with it
	//Below we create an instance of the JSON helper and associate
	//it with our redis connnection
	jsonHelper := rejson.NewReJSONHandler()
	jsonHelper.SetGoRedisClientWithContext(ctx, client)

	//Return a pointer to a new ToDo struct
	return &ToDo{
		cache: cache{
			cacheClient: client,
			jsonHelper:  jsonHelper,
			context:     ctx,
		},
	}, nil
}

//------------------------------------------------------------
// REDIS HELPERS
//------------------------------------------------------------

func isRedisNilError(err error) bool {
	return errors.Is(err, redis.Nil) || err.Error() == RedisNilError
}

func redisKeyFromId(id int) string {
	return fmt.Sprintf("%s%d", RedisKeyPrefix, id)
}

func (t *ToDo) getItemFromRedis(key string, item *Voter) error {
	itemObject, err := t.jsonHelper.JSONGet(key, ".")
	if err != nil {
		return err
	}

	err = json.Unmarshal(itemObject.([]byte), item)
	if err != nil {
		return err
	}

	return nil
}

func (t *ToDo) AddVoter(voter Voter) error {
	redisKey := redisKeyFromId(int(voter.VoterId))
	var existingVoter Voter
	if err := t.getItemFromRedis(redisKey, &existingVoter); err == nil {
		return errors.New("Voter already exists")
	}

	if _, err := t.jsonHelper.JSONSet(redisKey, ".", voter); err != nil {
		return err
	}

	return nil
}

func (t *ToDo) GetAllVoters() ([]Voter, error) {
	var voters []Voter
	var voter Voter

	pattern := RedisKeyPrefix + "*"
	ks, _ := t.cacheClient.Keys(t.context, pattern).Result()
	for _, key := range ks {
		err := t.getItemFromRedis(key, &voter)
		if err != nil {
			return nil, err
		}
		voters = append(voters, voter)
	}

	return voters, nil
}

func (t *ToDo) DeleteVoter(id int) error {
	pattern := redisKeyFromId(id)
	numDeleted, err := t.cacheClient.Del(t.context, pattern).Result()
	if err != nil {
		return err
	}
	if numDeleted == 0 {
		return errors.New("attempted to delete non-existent item")
	}

	return nil
}

func (t *ToDo) DeleteAll() error {
	pattern := RedisKeyPrefix + "*"
	keyStrings, _ := t.cacheClient.Keys(t.context, pattern).Result()
	numDeleted, err := t.cacheClient.Del(t.context, keyStrings...).Result()
	if err != nil {
		return err
	}

	if numDeleted != int64(len(keyStrings)) {
		return errors.New("one or more items could not be deleted")
	}

	return nil
}

func (t *ToDo) GetVoterPolls(voterId int) ([]VoterHistory, error) {
	redisKey := redisKeyFromId(voterId)
	var voter Voter
	if err := t.getItemFromRedis(redisKey, &voter); err != nil {
		return nil, err
	}

	return voter.VoteHistory, nil
}

func (t *ToDo) GetVoterPoll(voterId int, pollId int) (VoterHistory, error) {
	redisKey := redisKeyFromId(voterId)
	var voter Voter
	if err := t.getItemFromRedis(redisKey, &voter); err != nil {
		return VoterHistory{}, err
	}

	for _, poll := range voter.VoteHistory {
		if int(poll.PollId) == pollId {
			return poll, nil
		}
	}

	return VoterHistory{}, errors.New("poll not found")
}

func (t *ToDo) AddVoterPollHistory(voterId int, pollId int, voteDate time.Time) error {
	redisKey := redisKeyFromId(voterId)
	var voter Voter
	if err := t.getItemFromRedis(redisKey, &voter); err != nil {
		return err
	}

	voter.VoteHistory = append(voter.VoteHistory, VoterHistory{PollId: uint(pollId), VoteDate: voteDate})
	if _, err := t.jsonHelper.JSONSet(redisKey, ".", voter); err != nil {
		return err
	}

	return nil
}

// UpdateItem accepts a ToDoItem and updates it in the DB.
// Preconditions:   (1) The database file must exist and be a valid
//
//					(2) The item must exist in the DB
//	    				because we use the item.Id as the key, this
//						function must check if the item already
//	    				exists in the DB, if not, return an error
//
// Postconditions:
//
//	 (1) The item will be updated in the DB
//		(2) The DB file will be saved with the item updated
//		(3) If there is an error, it will be returned
func (t *ToDo) UpdateItem(item Voter) error {

	// Check if item exists before trying to update it
	// this is a good practice, return an error if the
	// item does not exist
	_, ok := t.toDoMap[int(item.VoterId)]
	if !ok {
		return errors.New("item does not exist")
	}

	//Now that we know the item exists, lets update it
	t.toDoMap[int(item.VoterId)] = item

	return nil
}

func (t *ToDo) GetVoter(id int) (Voter, error) {
	var voter Voter
	pattern := redisKeyFromId(id)
	err := t.getItemFromRedis(pattern, &voter)
	if err != nil {
		return Voter{}, err
	}

	return voter, nil
}

// PrintItem accepts a ToDoItem and prints it to the console
// in a JSON pretty format. As some help, look at the
// json.MarshalIndent() function from our in class go tutorial.
func (t *ToDo) PrintItem(item Voter) {
	jsonBytes, _ := json.MarshalIndent(item, "", "  ")
	fmt.Println(string(jsonBytes))
}

// PrintAllItems accepts a slice of ToDoItems and prints them to the console
// in a JSON pretty format.  It should call PrintItem() to print each item
// versus repeating the code.
func (t *ToDo) PrintAllItems(itemList []Voter) {
	for _, item := range itemList {
		t.PrintItem(item)
	}
}

// JsonToItem accepts a json string and returns a ToDoItem
// This is helpful because the CLI accepts todo items for insertion
// and updates in JSON format.  We need to convert it to a ToDoItem
// struct to perform any operations on it.
func (t *ToDo) JsonToItem(jsonString string) (Voter, error) {
	var item Voter
	err := json.Unmarshal([]byte(jsonString), &item)
	if err != nil {
		return Voter{}, err
	}

	return item, nil
}
