# Prerequisites/Dependencies:

- Go (I used go1.17.6)

- github.com/mattn/go-sqlite3

  Use commands "go get github.com/mattn/go-sqlite3" and then "go install github.com/mattn/go-sqlite3"

- SQLite (I used version 3.36.0)

# Useage:

Run the server by using the command "go run server.go" in the project directory.

This will run the server on localhost:4041.

Append a search parameter to the URL in the following form: http://localhost:4041/?search=jesus%20of%20suburbia

All track names that contain the search parameter will be given in JSON array.

Log will display recieved and completed search queries as well as error codes for failed requests.

Note that "%20" is used to denote spaces in the URL search parameter.
Currently other URL character encoding is not supported (ex. %22 for double quotes, %27 for apostrophe, etc.).

# Details:

Manually tested using Google Chrome and Mozilla Firefox.

Run testing script by using the command "go test" in the project directory.

Implmentation includes adding artist and album names from their respective tables. [Database info](https://data-xtractor.com/knowledgebase/chinook-database-sample/)
