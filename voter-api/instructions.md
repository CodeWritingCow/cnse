# First API Assignment - the "Voter" API
For this assignment you will be implementing a Voter API.  This API will be part of the final project deliverable.  Note that this is the initial version, we will be extending it a few more times. 

## BASIC REQUIREMENTS
A voter is a person that has the ability to vote in one or more "polls".   For example, a voter can vote in a poll - "What is better, coke or pepsi?", or a poll "What do you like better, dogs or cats?.  The ultimate course project will have very simple polls that a voter can vote on.  For now, we are implementing the first version of the Voter API.  We will be tracking voter history via the voterPoll struct but not what the poll was about, nor what the "voter voted for - e.g., I like coke better" here - it will just manage the fact that a particular voter voted in a particular poll.  For now the data structure of this API should be something similiar to the below (note feel free to change anything you want):
```
type VoterHistory struct{
  PollId uint
  VoteId uint 
  VoteDate time.Time
}

type Voter struct{
  VoterId uint
  Name string
  Email string
  VoteHistory []VoterHistory
}
type VoterList struct {
  Voters [uint]Voter //A map of VoterIDs as keys and Voter structs as values
}

//constructor for VoterList struct
func NewVoterList(){
  ...
}
//Add receivers to any structs you want, but at the minimum you should add the API behavior to the
//VoterList struct as its managing the collection of voters.  Also dont forget in the constructor
//that you need to make the map before you can use it - make map[uint]Voter   
```
Note, we will get into RESTFUL API design best practices a little later, but for now, just think about you are managing "a collection" of resources.  Thus, best practices states to name your REST endpoints with meaningful plural names.  Thus your API interface should follow (you can change, this is just a strong suggestion):
GET /voters - Get all voter resources including all voter history for each voter (note we will discuss the concept of "paging" later, for now you can ignore)
GET&POST /voters/:id - Get a single voter resource with voterID=:id including their entire voting history.  POST version adds one to the "database" 
GET /voters/:id/polls - Gets the JUST the voter history for the voter with VoterID = :id 
GET&POST /voters/:id/polls/:pollid - Gets JUST the single voter poll data with PollID = :id and VoterID = :id.  POST version adds one to the "database" 
GET /voters/health - Returns a "health" record indicating that the voter API is functioning properly and some metadata about the API.  Note the payload can be hard coded, we are mainly looking for a HTTP status code of 200, which means the API is functioning properly. 

Note since REST based APIs use HTTP, we use GET for Queries and POST for Adding records.  The above shows the minimal requirements of the API endpoints that need to support just GET and the endpoints that must also support GET and POST.
You must also create a test suite to showcase the expected operation of your voter API.  There is a good starting point example for you to use to build upon here:  https://github.com/ArchitectingSoftware/CNSE-Class-Demo-Files/tree/main/todo-api/tests.  We will be running your test suite to evaluate the correct function of your API using the following command from the root of your project folder:  go test ./... -v.   This command searches all of your packages for test suites and executes them all.  If you follow the directory structure in the todo-api sample, you will notice the sample test files are in the "tests" package under the "/tests" directory.  You can also run your test suites if they are all in the /tests directory by running go test ./tests -v

## HINTS/REMINDERS
1. For now, you don't have to save the data in a database or file.  You can save the information in your API using an in-memory slice of "Voter".  You can follow the pattern of the sample code I provide (see below)
2. Don't forget to "make" your slices in constructors, and then use append() to add elements to it.  AKA you will need to do this for the VoteHistory
3. You may need to do a little investigation of your own (this is by design) on working with time based structures in Go.  e.g., time.Time.  There is a lot of information online, or you can aks me questions.

## EXTRA CREDIT
There are a number of ways to get extra credit for this assignment, you can pick either one or both for more extra credit:
1. Add some json tags to the key Voter data structures (see the example in the todo API) to give your API payloads more friendly name for the json structures.  For example, tag the VoterId field to allow the json to be represented as voter_id, which is more inline with JSON best practices.
2. Add PUT and DELETE endpoints for \voters\:id and \voters\:id\polls\:id where PUT is used to update a record and DELETE is used to remove a record. Note that PUT/DELETE for /voters/:id should only modify the voter resource data, and PUT/DELETE /voters/:id/polls/:id should only modify the poll data for the selected voter
3. Add a realistic body to GET /voters/health.  Thus, instead of returning hard coded values, return useful metadata including API uptime, total API calls, total API calls with errors, etc.  You can get creative here.  Note that you will need to add additional fields to the structures for this.  For example, when the API first starts you can add a time.Time field called bootTime and set it to time.Now() in the constructor (aka New() function).  Then you can calculate the uptime duration with something like time.Now().Sub(bootTime).  You may need to research the APIs for this, it should be fun :-)
## TIPS & GETTING STARTED
Note that I will not be providing any scaffold code, however, the implementation should be able to closely follow the sample todo-api that we will cover in class.  See: https://github.com/ArchitectingSoftware/CNSE-Class-Demo-Files/tree/main/todo-api.  I highly suggest that you copy this code and then start modifying it to implement this assignment.  It should have just about everything that you need. 

This might seem like a lot of work based on reading above, but dont panic, its really not much work at all if you follow my todo-api as a reference. 
