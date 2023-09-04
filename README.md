# GoQuests

This is a simple REST API in Golang that manages quests in a database.
The purpose of this API is learning the GO language.

## Setup

- make sure you have Docker installed and Postgres installed and running
- run `docker run --name postgres_db  -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=quests -d postgres`
- copy example.env into .env
- run `go run main.go`
- visit `http://localhost:8008/quests`

## Endpoints

/quests - GET
/quest/{id} - GET
/quest - POST
/quest/{id} - PUT
/quest/{id} - DELETE
