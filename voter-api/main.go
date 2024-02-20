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

	r.GET("/voters", apiHandler.GetVoterList)
	r.GET("/voters/:id", apiHandler.GetVoter)
	r.POST("/voters/:id", apiHandler.AddVoter)
	r.GET("/voters/:id/polls", apiHandler.ListVoterPolls)
	r.GET("/voters/:id/polls/:pollid", apiHandler.GetVoterPoll)
	r.POST("/voters/:id/polls/:pollid", apiHandler.AddVoterPoll)
	r.GET("/voters/health", apiHandler.HealthCheck)

	// TODO: Remove unused boilerplate code
	// r.GET("/todo", apiHandler.ListAllTodos)
	// r.POST("/todo", apiHandler.AddToDo)
	// r.PUT("/todo", apiHandler.UpdateToDo)
	// r.DELETE("/todo", apiHandler.DeleteAllToDo)
	// r.DELETE("/todo/:id", apiHandler.DeleteToDo)
	// r.GET("/todo/:id", apiHandler.GetToDo)

	// r.GET("/crash", apiHandler.CrashSim)
	// r.GET("/health", apiHandler.HealthCheck)

	// v2 := r.Group("/v2")
	// v2.GET("/todo", apiHandler.ListSelectTodos)

	serverPath := fmt.Sprintf("%s:%d", hostFlag, portFlag)
	r.Run(serverPath)
}
