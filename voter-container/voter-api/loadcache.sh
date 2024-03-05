#!/bin/bash
curl -d '{ "voter_id": 0, "name": "Moo Moo" }' -H "Content-Type: application/json" -X POST http://localhost:1080/voters/0
curl -d '{ "voter_id": 1, "name": "Totoro" }' -H "Content-Type: application/json" -X POST http://localhost:1080/voters/1

curl -d '{ "voter_id": 0, "poll_id": 0 }' -H "Content-Type: application/json" -X POST http://localhost:1080/voters/0/polls/0
