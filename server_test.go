package main

import (
    "net/http"
    "net/http/httptest"
    "testing"
)

// Function to check the http response code for errors 
func ResponseCodeTest(rec *httptest.ResponseRecorder, req *http.Request,
	 err error, status int, t *testing.T) {
	if err != nil {
        t.Fatal(err)
    }

	handler(rec, req)

	// Check the status code is expected code
	if rec.Code != status {
		t.Errorf(
			"handler returned wrong status code: \n\ngot\n\n%v\n\nwant\n\n%v",
			rec.Code, status)
	}
}

// Function to check the if content from JSON response matches expected
func ResponseJSONTest(rec *httptest.ResponseRecorder, ctype string,
	 expected string, t *testing.T) {
	// Check content type
	if ctype != "application/json" {
		t.Errorf(
			"content type header does not match: \n\ngot\n\n%v\n\nwant\n\n%v",
			ctype, "application/json")
	}

	// Check body of response
	if rec.Body.String() != expected {
		t.Errorf(
			"handler returned unexpected body. \n\ngot:\n\n%v\n\nwant\n\n%v",
			rec.Body.String(), expected)
	}
}

func TestHandler(t *testing.T) {
	// Request 1: Ensure correct behaviour when searching 
	// for "Jesus of Suburbia" 
	// Create a ResponseRecorder 
	rec1 := httptest.NewRecorder()

	// Create a request to pass to our handler. 
	req1, err1 := http.NewRequest(http.MethodGet, 
		"http://localhost4041/?search=jesus%20of%20suburbia", nil)

	ResponseCodeTest(rec1, req1, err1, http.StatusOK, t)

	// Get content type
	ctype1 := rec1.Header().Get("Content-Type");

	// expected body to check against response body
	// Must be four spaces from margin (not tab) to correctly match output
	// from server.go
    expected1 := `[{
    "TrackId": 1134,
    "Name": "Jesus Of Suburbia / City Of The Damned / I Don't Care /` + 
	` Dearly Beloved / Tales Of Another Broken Home",
    "Artist": "Green Day",
    "Album": "American Idiot",
    "AlbumId": 89,
    "MediaTypeId": 1,
    "GenreId": 4,
    "Composer": "Billie Joe Armstrong/Green Day",
    "Milliseconds": 548336,
    "Bytes": 17875209,
    "UnitPrice": 0.99
}]`

	ResponseJSONTest(rec1, ctype1, expected1, t)

	// Request 2: Ensure correct behaviour when searching for "London"
	// Create a ResponseRecorder to record the response.
	rec2 := httptest.NewRecorder()

	// Create a request to pass to our handler. 
	req2, err2 := http.NewRequest(http.MethodGet, 
		"http://localhost4041/?search=london", nil)
	
	ResponseCodeTest(rec2, req2, err2, http.StatusOK, t)

	// Get content type
	ctype2 := rec2.Header().Get("Content-Type")

	// expected body to check against response body
	// Must be four spaces from margin (not tab) to correctly match output 
	// from server.go
	expected2 := `[{
    "TrackId": 2599,
    "Name": "London Calling",
    "Artist": "The Clash",
    "Album": "The Singles",
    "AlbumId": 211,
    "MediaTypeId": 1,
    "GenreId": 4,
    "Composer": "Joe Strummer/Mick Jones",
    "Milliseconds": 199706,
    "Bytes": 6569007,
    "UnitPrice": 0.99
},
{
    "TrackId": 3414,
    "Name": "Symphony No. 104 in D Major \"London\": IV. Finale: Spiritoso",
    "Artist": "Royal Philharmonic Orchestra \u0026 Sir Thomas Beecham",
    "Album": "Haydn: Symphonies 99 - 104",
    "AlbumId": 283,
    "MediaTypeId": 4,
    "GenreId": 24,
    "Composer": "Franz Joseph Haydn",
    "Milliseconds": 306687,
    "Bytes": 10085867,
    "UnitPrice": 0.99
}]`
	
	ResponseJSONTest(rec2, ctype2, expected2, t)

	// Request 3: Ensure correct behaviour when sending an empty search param
	// Create a ResponseRecorder to record the response.
	rec3 := httptest.NewRecorder()

	// Create a request to pass to our handler. 
	req3, err3 := http.NewRequest(http.MethodGet, 
		"http://localhost4041/?search=", nil)
	
	ResponseCodeTest(rec3, req3, err3, http.StatusBadRequest, t)

	// Request 4: Ensure correct behaviour when sending an unsupported method
	// Create a ResponseRecorder to record the response.
	rec4 := httptest.NewRecorder()

	// Create a request to pass to our handler. 
	req4, err4 := http.NewRequest("POST", "http://localhost4041/", nil)

	ResponseCodeTest(rec4, req4, err4, http.StatusMethodNotAllowed, t)

	// Request 5: Ensure correct behaviour when sending no search param
	// Create a ResponseRecorder to record the response.
	rec5 := httptest.NewRecorder()

	// Create a request to pass to our handler. 
	req5, err5 := http.NewRequest(http.MethodGet, "http://localhost4041/", nil)

	ResponseCodeTest(rec5, req5, err5, http.StatusBadRequest, t)	

	// Request 6: Ensure correct behaviour when sending apostrophe 
	// for database query
	// Create a ResponseRecorder to record the response.
	rec6 := httptest.NewRecorder()

	// Create a request to pass to our handler. 
	req6, err6 := http.NewRequest(http.MethodGet, 
		"http://localhost4041/?search=please%20don%27t%20touch", nil)
	
	ResponseCodeTest(rec6, req6, err6, http.StatusOK, t)
	
	// Get content type
	ctype6 := rec6.Header().Get("Content-Type")

	// expected body to check against response body
	// Must be four spaces from margin (not tab) to correctly match output 
	// from server.go
	expected6 := `[{
    "TrackId": 1955,
    "Name": "Please Don't Touch",
    "Artist": "Motörhead",
    "Album": "Ace Of Spades",
    "AlbumId": 160,
    "MediaTypeId": 1,
    "GenreId": 3,
    "Composer": "Heath/Robinson",
    "Milliseconds": 169926,
    "Bytes": 5557002,
    "UnitPrice": 0.99
}]`
	
	ResponseJSONTest(rec6, ctype6, expected6, t)

	// Request 7: Ensure correct behaviour when query for "Jump" 
	// Create a ResponseRecorder to record the response.
	rec7 := httptest.NewRecorder()

	// Create a request to pass to our handler. 
	req7, err7 := http.NewRequest(http.MethodGet, 
		"http://localhost4041/?search=jump", nil)
	
	ResponseCodeTest(rec7, req7, err7, http.StatusOK, t)
	
	// Get content type
	ctype7 := rec7.Header().Get("Content-Type")

	// expected body to check against response body
	// Must be four spaces from margin (not tab) to correctly match output 
	// from server.go
	expected7 := `[{
    "TrackId": 3070,
    "Name": "Jump",
    "Artist": "Van Halen",
    "Album": "The Best Of Van Halen, Vol. I",
    "AlbumId": 243,
    "MediaTypeId": 1,
    "GenreId": 1,
    "Composer": "Edward Van Halen, Alex Van Halen, David Lee Roth",
    "Milliseconds": 241711,
    "Bytes": 7911090,
    "UnitPrice": 0.99
},
{
    "TrackId": 3300,
    "Name": "Jump Around",
    "Artist": "House Of Pain",
    "Album": "House of Pain",
    "AlbumId": 258,
    "MediaTypeId": 1,
    "GenreId": 17,
    "Composer": "E. Schrody/L. Muggerud",
    "Milliseconds": 217835,
    "Bytes": 8715653,
    "UnitPrice": 0.99
},
{
    "TrackId": 3317,
    "Name": "Jump Around (Pete Rock Remix)",
    "Artist": "House Of Pain",
    "Album": "House of Pain",
    "AlbumId": 258,
    "MediaTypeId": 1,
    "GenreId": 17,
    "Composer": "E. Schrody/L. Muggerud",
    "Milliseconds": 236120,
    "Bytes": 9447101,
    "UnitPrice": 0.99
},
{
    "TrackId": 1832,
    "Name": "Jump In The Fire",
    "Artist": "Metallica",
    "Album": "Kill 'Em All",
    "AlbumId": 150,
    "MediaTypeId": 1,
    "GenreId": 3,
    "Composer": "James Hetfield, Lars Ulrich, Dave Mustaine",
    "Milliseconds": 281573,
    "Bytes": 9135755,
    "UnitPrice": 0.99
},
{
    "TrackId": 632,
    "Name": "Disc Jockey Jump",
    "Artist": "Gene Krupa",
    "Album": "Up An' Atom",
    "AlbumId": 51,
    "MediaTypeId": 1,
    "GenreId": 2,
    "Composer": null,
    "Milliseconds": 193149,
    "Bytes": 6260820,
    "UnitPrice": 0.99
},
{
    "TrackId": 198,
    "Name": "When My Left Eye Jumps",
    "Artist": "Buddy Guy",
    "Album": "The Best Of Buddy Guy - The Millenium Collection",
    "AlbumId": 20,
    "MediaTypeId": 1,
    "GenreId": 6,
    "Composer": "Al Perkins/Willie Dixon",
    "Milliseconds": 235311,
    "Bytes": 7685363,
    "UnitPrice": 0.99
}]`
	
	ResponseJSONTest(rec7, ctype7, expected7, t)

	// Request 8: Ensure correct behaviour with empty limit and offset 
	// Create a ResponseRecorder to record the response.
	rec8 := httptest.NewRecorder()

	// Create a request to pass to our handler. 
	req8, err8 := http.NewRequest(http.MethodGet, 
		"http://localhost:4041/?search=jump&limit=&offset=", nil)
	
	ResponseCodeTest(rec8, req8, err8, http.StatusOK, t)
	
	// Get content type
	ctype8 := rec8.Header().Get("Content-Type")

	// expected body to check against response body
	// Must be four spaces from margin (not tab) to correctly match output 
	// from server.go
	expected8 := `[{
    "TrackId": 3070,
    "Name": "Jump",
    "Artist": "Van Halen",
    "Album": "The Best Of Van Halen, Vol. I",
    "AlbumId": 243,
    "MediaTypeId": 1,
    "GenreId": 1,
    "Composer": "Edward Van Halen, Alex Van Halen, David Lee Roth",
    "Milliseconds": 241711,
    "Bytes": 7911090,
    "UnitPrice": 0.99
},
{
    "TrackId": 3300,
    "Name": "Jump Around",
    "Artist": "House Of Pain",
    "Album": "House of Pain",
    "AlbumId": 258,
    "MediaTypeId": 1,
    "GenreId": 17,
    "Composer": "E. Schrody/L. Muggerud",
    "Milliseconds": 217835,
    "Bytes": 8715653,
    "UnitPrice": 0.99
},
{
    "TrackId": 3317,
    "Name": "Jump Around (Pete Rock Remix)",
    "Artist": "House Of Pain",
    "Album": "House of Pain",
    "AlbumId": 258,
    "MediaTypeId": 1,
    "GenreId": 17,
    "Composer": "E. Schrody/L. Muggerud",
    "Milliseconds": 236120,
    "Bytes": 9447101,
    "UnitPrice": 0.99
},
{
    "TrackId": 1832,
    "Name": "Jump In The Fire",
    "Artist": "Metallica",
    "Album": "Kill 'Em All",
    "AlbumId": 150,
    "MediaTypeId": 1,
    "GenreId": 3,
    "Composer": "James Hetfield, Lars Ulrich, Dave Mustaine",
    "Milliseconds": 281573,
    "Bytes": 9135755,
    "UnitPrice": 0.99
},
{
    "TrackId": 632,
    "Name": "Disc Jockey Jump",
    "Artist": "Gene Krupa",
    "Album": "Up An' Atom",
    "AlbumId": 51,
    "MediaTypeId": 1,
    "GenreId": 2,
    "Composer": null,
    "Milliseconds": 193149,
    "Bytes": 6260820,
    "UnitPrice": 0.99
},
{
    "TrackId": 198,
    "Name": "When My Left Eye Jumps",
    "Artist": "Buddy Guy",
    "Album": "The Best Of Buddy Guy - The Millenium Collection",
    "AlbumId": 20,
    "MediaTypeId": 1,
    "GenreId": 6,
    "Composer": "Al Perkins/Willie Dixon",
    "Milliseconds": 235311,
    "Bytes": 7685363,
    "UnitPrice": 0.99
}]`
	
	ResponseJSONTest(rec8, ctype8, expected8, t)

	// Request 9: Ensure correct behaviour with non-digits as offset
	// Create a ResponseRecorder to record the response.
	rec9 := httptest.NewRecorder()

	// Create a request to pass to our handler. 
	req9, err9 := http.NewRequest(http.MethodGet, 
		"http://localhost:4041/?search=jump&limit=5&offset=a", nil)
	
	ResponseCodeTest(rec9, req9, err9, http.StatusBadRequest, t)
	
	// Request 10: Ensure correct behaviour with non-digits as limit
	// Create a ResponseRecorder to record the response.
	rec10 := httptest.NewRecorder()

	// Create a request to pass to our handler. 
	req10, err10 := http.NewRequest(http.MethodGet, 
		"http://localhost:4041/?search=jump&limit=a&offset=1", nil)
	
	ResponseCodeTest(rec10, req10, err10, http.StatusBadRequest, t)

	// Request 11: Ensure correct behaviour when using valid limit
	// Create a ResponseRecorder to record the response.
	rec11 := httptest.NewRecorder()

	// Create a request to pass to our handler. 
	req11, err11 := http.NewRequest(http.MethodGet, 
		"http://localhost:4041/?search=jump&limit=1", nil)
	
	ResponseCodeTest(rec11, req11, err11, http.StatusOK, t)
	
	// Get content type
	ctype11 := rec11.Header().Get("Content-Type")

	// expected body to check against response body
	// Must be four spaces from margin (not tab) to correctly match output 
	// from server.go
	expected11 := `[{
    "TrackId": 3070,
    "Name": "Jump",
    "Artist": "Van Halen",
    "Album": "The Best Of Van Halen, Vol. I",
    "AlbumId": 243,
    "MediaTypeId": 1,
    "GenreId": 1,
    "Composer": "Edward Van Halen, Alex Van Halen, David Lee Roth",
    "Milliseconds": 241711,
    "Bytes": 7911090,
    "UnitPrice": 0.99
}]`
	
	ResponseJSONTest(rec11, ctype11, expected11, t)

	// Request 12: Ensure correct behaviour when using valid limit and offset
	// Create a ResponseRecorder to record the response.
	rec12 := httptest.NewRecorder()

	// Create a request to pass to our handler. 
	req12, err12 := http.NewRequest(http.MethodGet, 
		"http://localhost:4041/?search=jump&limit=2&offset=3", nil)
	
	ResponseCodeTest(rec12, req12, err12, http.StatusOK, t)
	
	// Get content type
	ctype12 := rec12.Header().Get("Content-Type")

	// expected body to check against response body
	// Must be four spaces from margin (not tab) to correctly match output 
	// from server.go
	expected12 := `[{
    "TrackId": 1832,
    "Name": "Jump In The Fire",
    "Artist": "Metallica",
    "Album": "Kill 'Em All",
    "AlbumId": 150,
    "MediaTypeId": 1,
    "GenreId": 3,
    "Composer": "James Hetfield, Lars Ulrich, Dave Mustaine",
    "Milliseconds": 281573,
    "Bytes": 9135755,
    "UnitPrice": 0.99
},
{
    "TrackId": 632,
    "Name": "Disc Jockey Jump",
    "Artist": "Gene Krupa",
    "Album": "Up An' Atom",
    "AlbumId": 51,
    "MediaTypeId": 1,
    "GenreId": 2,
    "Composer": null,
    "Milliseconds": 193149,
    "Bytes": 6260820,
    "UnitPrice": 0.99
}]`
	
	ResponseJSONTest(rec12, ctype12, expected12, t)

	// Request 13: Ensure correct behaviour when offset without limit
	// Create a ResponseRecorder to record the response.
	rec13 := httptest.NewRecorder()

	// Create a request to pass to our handler. 
	req13, err13 := http.NewRequest(http.MethodGet, 
		"http://localhost:4041/?search=jump&offset=3", nil)
	
	ResponseCodeTest(rec13, req13, err13, http.StatusOK, t)
	
	// Get content type
	ctype13 := rec13.Header().Get("Content-Type")

	// expected body to check against response body
	// Must be four spaces from margin (not tab) to correctly match output 
	// from server.go
	expected13 := `[{
    "TrackId": 3070,
    "Name": "Jump",
    "Artist": "Van Halen",
    "Album": "The Best Of Van Halen, Vol. I",
    "AlbumId": 243,
    "MediaTypeId": 1,
    "GenreId": 1,
    "Composer": "Edward Van Halen, Alex Van Halen, David Lee Roth",
    "Milliseconds": 241711,
    "Bytes": 7911090,
    "UnitPrice": 0.99
},
{
    "TrackId": 3300,
    "Name": "Jump Around",
    "Artist": "House Of Pain",
    "Album": "House of Pain",
    "AlbumId": 258,
    "MediaTypeId": 1,
    "GenreId": 17,
    "Composer": "E. Schrody/L. Muggerud",
    "Milliseconds": 217835,
    "Bytes": 8715653,
    "UnitPrice": 0.99
},
{
    "TrackId": 3317,
    "Name": "Jump Around (Pete Rock Remix)",
    "Artist": "House Of Pain",
    "Album": "House of Pain",
    "AlbumId": 258,
    "MediaTypeId": 1,
    "GenreId": 17,
    "Composer": "E. Schrody/L. Muggerud",
    "Milliseconds": 236120,
    "Bytes": 9447101,
    "UnitPrice": 0.99
},
{
    "TrackId": 1832,
    "Name": "Jump In The Fire",
    "Artist": "Metallica",
    "Album": "Kill 'Em All",
    "AlbumId": 150,
    "MediaTypeId": 1,
    "GenreId": 3,
    "Composer": "James Hetfield, Lars Ulrich, Dave Mustaine",
    "Milliseconds": 281573,
    "Bytes": 9135755,
    "UnitPrice": 0.99
},
{
    "TrackId": 632,
    "Name": "Disc Jockey Jump",
    "Artist": "Gene Krupa",
    "Album": "Up An' Atom",
    "AlbumId": 51,
    "MediaTypeId": 1,
    "GenreId": 2,
    "Composer": null,
    "Milliseconds": 193149,
    "Bytes": 6260820,
    "UnitPrice": 0.99
},
{
    "TrackId": 198,
    "Name": "When My Left Eye Jumps",
    "Artist": "Buddy Guy",
    "Album": "The Best Of Buddy Guy - The Millenium Collection",
    "AlbumId": 20,
    "MediaTypeId": 1,
    "GenreId": 6,
    "Composer": "Al Perkins/Willie Dixon",
    "Milliseconds": 235311,
    "Bytes": 7685363,
    "UnitPrice": 0.99
}]`
	
	ResponseJSONTest(rec13, ctype13, expected13, t)

}