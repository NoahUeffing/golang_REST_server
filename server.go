package main

import (
	"log"
	"os"
    "fmt"
	"net/http"
    "database/sql"
	"encoding/json"
    _ "github.com/mattn/go-sqlite3"
)

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

// MarshalJSON for NullString
func (ni *NullInt64) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ni.Int64)
}

// MarshalJSON for NullString
func (nf *NullFloat64) MarshalJSON() ([]byte, error) {
	if !nf.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nf.Float64)
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.SetOutput(os.Stdout)
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed) // Return 405 Method Not Allowed.
		log.Println("405 Error: Method not allowed")
		return
	}
	// URL search term keys
	keys, ok := r.URL.Query()["search"]

	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("400 Error: Bad request")
		return
	}

	if len(keys[0]) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("400 Error: No valid search criteria")
		return
	}

	// Query()["key"] will return an array of items, 
	// we only want the single item.
	key := keys[0]

	log.Println("Received search query for: " + string(key))

    // Open up our database connection.
    db, err := sql.Open("sqlite3", "./Chinook_Sqlite.sqlite")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("500 Error: Database connection error")
		return
	}
	// query)
	results, err := db.Query("SELECT track.TrackId, track.Name, artist.Name, album.Title, track.AlbumId, track.MediaTypeId, track.GenreId, track.Composer, track.Milliseconds, track.Bytes, track.UnitPrice FROM track INNER JOIN album ON track.AlbumId = album.AlbumId INNER JOIN artist ON album.ArtistId = artist.ArtistId WHERE track.Name LIKE '%" + string(key) + "%'")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("500 Error: Database error")
		return
	}
	var track Track
	fmt.Fprintf(w, "[")
    for results.Next() {
		if err = results.Scan(&track.TrackId, &track.Name, &track.Artist, &track.Album, &track.AlbumId, &track.MediaTypeId, &track.GenreId,
			&track.Composer, &track.Milliseconds, &track.Bytes, &track.UnitPrice); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println("500 Error: Server error")
				return
			}	

		trackJSON, err := json.Marshal(&track)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("500 Error: Encoding error")
			return
		} else {
			fmt.Fprintf(w, "%s,", trackJSON)
		}
	}

	fmt.Fprintf(w, "]")
    results.Close() 
	db.Close()
	log.Println("Search query completed for: " + string(key))
	return
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./note.ico")
}

func main() {
	http.HandleFunc("/favicon.ico", faviconHandler)
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":4041", nil))
}