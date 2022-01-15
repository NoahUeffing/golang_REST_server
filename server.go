/* 
A basic Go server used to query tracks from the Chinook_Sqlite database. 
Written by Noah Ueffing 
*/

package main

import (
	"log"
	"os"
	"fmt"
	"net/http"
	"database/sql"
	"encoding/json"
	"strings"
	"strconv"
	_ "github.com/mattn/go-sqlite3"
)

// struct used for converting track data to json form
// Null datatypes are used in case attributes are missing for track instances
type Track struct {
	TrackId NullInt64 `json:"TrackId"`
	Name NullString `json:"Name"`
	Artist NullString `json:"Artist"`
	Album NullString `json:"Album"`
	AlbumId NullInt64 `json:"AlbumId"`
	MediaTypeId NullInt64 `json:"MediaTypeId"`
	GenreId NullInt64 `json:"GenreId"`
	Composer NullString `json:"Composer"`
	Milliseconds NullInt64 `json:"Milliseconds"`
	Bytes NullInt64 `json:"Bytes"`
	UnitPrice NullFloat64 `json:"UnitPrice"`
}

// NullString is an alias for sql.NullString data type
type NullString struct {
	sql.NullString
}

// NullInt64 is an alias for sql.NullString data type
type NullInt64 struct {
	sql.NullInt64
}

// NullFloat64 is an alias for sql.NullString data type
type NullFloat64 struct {
	sql.NullFloat64
}

// MarshalJSON for NullString
func (ns *NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.String)
}

// MarshalJSON for NullInt64
func (ni *NullInt64) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ni.Int64)
}

// MarshalJSON for NullFloat64
func (nf *NullFloat64) MarshalJSON() ([]byte, error) {
	if !nf.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nf.Float64)
}

// Function to send http error response and print error message to log
func errorHandler(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	log.Println(message)
}

// Request handler function for search queries
func handler(w http.ResponseWriter, r *http.Request) {
	// Set log ouput to Stdout
	log.SetOutput(os.Stdout)

	// Make sure the request is a GET request, otherwise give error
	if r.Method != http.MethodGet {
		errorHandler(w, http.StatusMethodNotAllowed, 
			"405 Error: Method not allowed")
		return
	}
	// Read the URL for parameters
	searchTerms, ok := r.URL.Query()["search"]
	inputLimit := r.URL.Query()["limit"]
	inputOffset := r.URL.Query()["offset"]

	// Error checking for search parameter
	if !ok {
		errorHandler(w, http.StatusBadRequest, "400 Error: Bad request")
		return
	}

	// If search parameter is empty, give error
	if len(searchTerms[0]) < 1 {
		errorHandler(w, http.StatusBadRequest, 
			"400 Error: No valid search criteria")
		return
	}

	// Query()["search"] will return an array of parameters, 
	// we only want a single parameter
	search := searchTerms[0]
	var limit string
	var offset string
	if len(inputLimit) > 0 {
		limit = inputLimit[0]
	}
	if len(inputOffset) > 0 {
		offset = inputOffset[0]
	}
	
	// Replace ' with '' for SQL query functionality
	search = strings.Replace(search, "'", "''", -1) 

	// log the recieved search query
	log.Println("Received search query for: " + string(search))

	// Open the database connection
	db, err := sql.Open("sqlite3", "./Chinook_Sqlite.sqlite")
	if err != nil {
		errorHandler(w, http.StatusInternalServerError, 
			"500 Error: Database connection error")
		return
	}

	// Create the query to search for tracks and order by best match
	// query for non limit or offset
	query := "SELECT track.TrackId, track.Name, artist.Name, album.Title, " + 
		"track.AlbumId, track.MediaTypeId, track.GenreId, track.Composer, " + 
		"track.Milliseconds, track.Bytes, track.UnitPrice " + 
		"FROM track " + 
		"INNER JOIN album ON track.AlbumId = album.AlbumId " +
		"INNER JOIN artist ON album.ArtistId = artist.ArtistId " + 
		"WHERE track.Name LIKE '%" + string(search) + "%' " +
		"ORDER BY (CASE WHEN track.Name = '" + string(search) + "' THEN 1 " +
		"WHEN track.name LIKE '" + string(search) + "%' THEN 2 ELSE 3 END), " + 
		"track.Name"
	// query with limit
	query1 := "SELECT track.TrackId, track.Name, artist.Name, album.Title, " + 
		"track.AlbumId, track.MediaTypeId, track.GenreId, track.Composer, " + 
		"track.Milliseconds, track.Bytes, track.UnitPrice " + 
		"FROM track " + 
		"INNER JOIN album ON track.AlbumId = album.AlbumId " +
		"INNER JOIN artist ON album.ArtistId = artist.ArtistId " + 
		"WHERE track.Name LIKE '%" + string(search) + "%' " +
		"ORDER BY (CASE WHEN track.Name = '" + string(search) + "' THEN 1 " +
		"WHEN track.name LIKE '" + string(search) + "%' THEN 2 ELSE 3 END), " + 
		"track.Name LIMIT " + string(limit)
	// query with limit and offset
	query2 := "SELECT track.TrackId, track.Name, artist.Name, album.Title, " + 
		"track.AlbumId, track.MediaTypeId, track.GenreId, track.Composer, " + 
		"track.Milliseconds, track.Bytes, track.UnitPrice " + 
		"FROM track " + 
		"INNER JOIN album ON track.AlbumId = album.AlbumId " +
		"INNER JOIN artist ON album.ArtistId = artist.ArtistId " + 
		"WHERE track.Name LIKE '%" + string(search) + "%' " +
		"ORDER BY (CASE WHEN track.Name = '" + string(search) + "' THEN 1 " +
		"WHEN track.name LIKE '" + string(search) + "%' THEN 2 ELSE 3 END), " + 
		"track.Name LIMIT " + string(limit) + " OFFSET " + string(offset)
	
	// Search the database
	var results *sql.Rows
	if len(limit) < 1 {
		results, err = db.Query(query) 
	} else if len(offset) < 1 {
		if _, err := strconv.Atoi(limit); err != nil {
			errorHandler(w, http.StatusBadRequest, "400 Error: Bad request")
			return
		} 
		results, err = db.Query(query1)
	} else {
		if _, err := strconv.Atoi(limit); err != nil {
			errorHandler(w, http.StatusBadRequest, "400 Error: Bad request")
			return
		} 
		if _, err := strconv.Atoi(offset); err != nil {
			errorHandler(w, http.StatusBadRequest, "400 Error: Bad request")
			return
		}
		results, err = db.Query(query2)
	}
	if err != nil {
		errorHandler(w, http.StatusInternalServerError, 
			"500 Error: Database error")
		return
	}

	// Create a Track item for each returned track from the search 
	// and write as an array of JSON objects
	var track Track
	count := 0
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "[")
	for results.Next() {
		if err = results.Scan(&track.TrackId, &track.Name, &track.Artist, 
			&track.Album, &track.AlbumId, &track.MediaTypeId, &track.GenreId,
			&track.Composer, &track.Milliseconds, &track.Bytes, 
			&track.UnitPrice); err != nil {
				errorHandler(w, http.StatusInternalServerError, 
					"500 Error: Server error")
				return
			}	

		trackJSON, err := json.MarshalIndent(&track, "", "    ")
		if err != nil {
			errorHandler(w, http.StatusInternalServerError, 
				"500 Error: Encoding error")
			return
		}
		// On first iteration, omit comma for array
		if count == 0 {
			fmt.Fprintf(w, "%s", trackJSON)
		} else {
			fmt.Fprintf(w, ",\n%s", trackJSON)
		}
		count = count + 1
	}

	fmt.Fprintf(w, "]")
	results.Close() 
	db.Close()
	log.Println("Search query completed for: " + string(search))
	return
}

// Function used to pass favicon
func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./note.ico")
}

// Driver function
func main() {
	// Pass favicon
	http.HandleFunc("/favicon.ico", faviconHandler)

	// Function to handle incoming requests
	http.HandleFunc("/", handler)

	// Listen for request on port 4041
	log.Fatal(http.ListenAndServe(":4041", nil))
}