package main

import (
	"flag"
	"fmt"
	"os"

	"voter-api/api"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	hostFlag string
	portFlag uint
)

func processCmdLineFlags() {
	flag.StringVar(&hostFlag, "h", "0.0.0.0", "Listen on all interfaces")
	flag.UintVar(&portFlag, "p", 1080, "Default Port")

	flag.Parse()
}

func main() {
	processCmdLineFlags()
	router := gin.Default()
	router.Use(cors.Default())

	apiHandler, err := api.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// TODO: Implement PUT routes for extra credit
	router.GET("/voters", apiHandler.GetAllVoters)
	router.GET("/voters/:id", apiHandler.GetVoter)
	router.POST("/voters/:id", apiHandler.AddVoter)
	// router.PUT("/voters/:id", apiHandler.UpdateVoter)
	router.DELETE("/voters/:id", apiHandler.DeleteVoter)
	router.DELETE("/voters", apiHandler.DeleteAllVoters)
	router.GET("/voters/:id/polls", apiHandler.GetVoterPolls)
	router.GET("/voters/:id/polls/:pollid", apiHandler.GetVoterPoll)
	router.POST("/voters/:id/polls/:pollid", apiHandler.AddVoterPoll)
	// router.PUT("/voters/:id/polls/:pollid", apiHandler.UpdateVoterPoll)
	router.DELETE("/voters/:id/polls/:pollid", apiHandler.DeleteVoterPoll)
	router.GET("/voters/health", apiHandler.HealthCheck)

	serverPath := fmt.Sprintf("%s:%d", hostFlag, portFlag)
	router.Run(serverPath)
}
