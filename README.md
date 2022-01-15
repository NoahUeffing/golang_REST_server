# Prerequisites/Dependencies:

- Go (I used go1.17.6)

- github.com/mattn/go-sqlite3

  Use commands "go get github.com/mattn/go-sqlite3" and then "go install github.com/mattn/go-sqlite3"

- SQLite (I used version 3.36.0)

- Docker (if building and running a Docker image as a container)

# Local Useage:

Run the server by using the command "go run server.go" in the project directory.

This will run the server on localhost:4041.

Append a search parameter to the URL in the following form: http://localhost:4041/?search=jesus%20of%20suburbia

All track names that contain the search parameter will be given in JSON array.

Log will display recieved and completed search queries as well as error codes for failed requests.

Note that "%20" is used to denote spaces in the URL search parameter, %27 for apostrophe, %3B for semicolon etc. See all character encodings [here.](https://www.w3schools.com/tags/ref_urlencode.ASP)

Run testing script by using the command "go test" in the project directory.

# Docker Useage:

To build the Docker image, with Docker running, use the command "docker build --tag golang-rest-server ./" in the project directory.

Then to run the Docker image as a container, use the command "docker run --publish 4041:4041 golang-rest-server" in the project directory.

# Details:

Manually tested using Google Chrome, Mozilla Firefox and Postman.

Implmentation includes adding artist and album names from their respective tables. [Database info](https://data-xtractor.com/knowledgebase/chinook-database-sample/)
