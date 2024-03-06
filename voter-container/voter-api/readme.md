## How to run Voter API with Docker and Redis

1. Open a terminal and go to /voter-api
2. Run build-docker.sh to containerize Voter API
3. Run `docker compose up` to start API
4. Run loadcache.sh to load Redis cache with data
5. Run `docker compose down` to shut down API
6. (Optional) Use commands from makefile to interact with API endpoints
