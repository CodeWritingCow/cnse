SHELL := /bin/bash

.PHONY: help
help:
	@echo "Usage make <TARGET>"
	@echo ""
	@echo "  Targets:"
	@echo "	   build				Build the todo executable"
	@echo "	   run					Run the todo program from code"
	@echo "	   run-bin				Run the todo executable"
	@echo "	   load-db				Add sample data via curl"
	@echo "	   get-by-id			Get a todo by id pass id=<id> on command line"
	@echo "	   get-all				Get all todos"
	@echo "	   update-2				Update record 2, pass a new title in using title=<title> on command line"
	@echo "	   delete-all			Delete all todos"
	@echo "	   delete-by-id			Delete a todo by id pass id=<id> on command line"
	@echo "	   get-v2				Get all todos by done status pass done=<true|false> on command line"
	@echo "	   get-v2-all			Get all todos using version 2"
	@echo "	   build-amd64-linux	Build amd64/Linux executable"
	@echo "	   build-arm64-linux	Build arm64/Linux executable"
	@echo "	   get-voters			Get all voters"
	@echo "	   get-voter-by-id		Get a voter by id, pass id=<voter_id> on command line"
	@echo "	   add-voter			Add a voter by id and name, pass id=<voter_id> and name=\"<name>\" on command line"
	@echo "	   get-polls			Get all polls for a voter, pass id=<voter_id> on command line"
	@echo "	   get-poll-by-id		Get a poll for a voter, pass id=<voter_id> and pollid=<poll_id> on command line"
	@echo "	   add-poll				Add a poll to a voter, pass id=<voter_id> and pollid=<poll_id> on command line"



.PHONY: build
build:
	go build .

.PHONY: build-amd64-linux
build-amd64-linux:
	GOOS=linux GOARCH=amd64 go build -o ./todo-linux-amd64 .

.PHONY: build-arm64-linux
build-arm64-linux:
	GOOS=linux GOARCH=arm64 go build -o ./todo-linux-arm64 .

	
.PHONY: run
run:
	go run main.go

.PHONY: run-bin
run-bin:
	./todo

.PHONY: restore-db
restore-db:
	(cp ./data/todo.json.bak ./data/todo.json)

.PHONY: restore-db-windows
restore-db-windows:
	(copy.\data\todo.json.bak .\data\todo.json)

.PHONY: load-db
load-db:
	curl -d '{ "id": 1, "title": "Learn Go / GoLang", "done": false }' -H "Content-Type: application/json" -X POST http://localhost:1080/todo 
	curl -d '{ "id": 2, "title": "Learn Kubernetes", "done": true}' -H "Content-Type: application/json" -X POST http://localhost:1080/todo 
	curl -d '{ "id": 3, "title": "Learn Cloud Native Architecture","done": false}' -H "Content-Type: application/json" -X POST http://localhost:1080/todo 
	curl -d '{"id": 4,"title": "Learn Why Professor Mitchell is the BEST! :-)","done": true}' -H "Content-Type: application/json" -X POST http://localhost:1080/todo

.PHONY: update-2
update-2:
	curl -d '{ "id": 2, "title": "$(title)", "done": false }' -H "Content-Type: application/json" -X PUT http://localhost:1080/todo 

.PHONY: get-by-id
get-by-id:
	curl -w "HTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X GET http://localhost:1080/todo/$(id) 

.PHONY: get-all
get-all:
	curl -w "HTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X GET http://localhost:1080/todo 

.PHONY: delete-all
delete-all:
	curl -w "HTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X DELETE http://localhost:1080/todo 

.PHONY: delete-by-id
delete-by-id:
	curl -w "HTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X DELETE http://localhost:1080/todo/$(id) 

.PHONY: get-v2
get-v2:
	curl -w "HTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X GET http://localhost:1080/v2/todo?done=$(done) 

.PHONY: get-v2-all
get-v2-all:
	curl -w "HTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X GET http://localhost:1080/v2/todo

.PHONY: get-voters
get-voters:
	curl -w "HTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X GET http://localhost:1080/voters 

.PHONY: get-voter-by-id
get-voter-by-id:
	curl -w "HTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X GET http://localhost:1080/voters/$(id) 

.PHONY: add-voter
add-voter:
	curl -d '{ "voter_id": $(id), "name": "$(name)" }' -H "Content-Type: application/json" -X POST http://localhost:1080/voters/$(id)

.PHONY: get-polls
get-polls:
	curl -w "HTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X GET http://localhost:1080/voters/$(id)/polls

.PHONY: get-poll-by-id
get-poll-by-id:
	curl -w "HTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X GET http://localhost:1080/voters/$(id)/polls/$(pollid)

.PHONY: add-poll
add-poll:
	curl -d '{ "voter_id": $(id), "poll_id": $(pollid) }' -H "Content-Type: application/json" -X POST http://localhost:1080/voters/$(id)/polls/$(pollid)
