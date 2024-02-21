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
	r := gin.Default()
	r.Use(cors.Default())

	apiHandler, err := api.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// TODO: Implement PUT and DELETE routes for extra credit
	r.GET("/voters", apiHandler.GetVoterList)
	r.GET("/voters/:id", apiHandler.GetVoter)
	r.POST("/voters/:id", apiHandler.AddVoter)
	// r.PUT("/voters/:id", apiHandler.UpdateVoter)
	r.DELETE("/voters/:id", apiHandler.DeleteVoter)
	r.GET("/voters/:id/polls", apiHandler.ListVoterPolls)
	r.GET("/voters/:id/polls/:pollid", apiHandler.GetVoterPoll)
	r.POST("/voters/:id/polls/:pollid", apiHandler.AddVoterPoll)
	// r.PUT("/voters/:id/polls/:pollid", apiHandler.UpdateVoterPoll)
	// r.DELETE("/voters/:id/polls/:pollid", apiHandler.DeleteVoterPoll)
	r.GET("/voters/health", apiHandler.HealthCheck)

	serverPath := fmt.Sprintf("%s:%d", hostFlag, portFlag)
	r.Run(serverPath)
}
