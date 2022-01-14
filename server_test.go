package main

import (
    "net/http"
    "net/http/httptest"
    "testing"
)

// Function to check the http response code for errors 
func ResponseCodeTest(rec *httptest.ResponseRecorder, req *http.Request, err error, status int, t *testing.T) {
	if err != nil {
        t.Fatal(err)
    }

	handler(rec, req)

	// Check the status code is expected code
	if rec.Code != status {
		t.Errorf("handler returned wrong status code: \n\ngot\n\n%v\n\nwant\n\n%v",
			rec.Code, status)
	}
}

// Function to check the if content from JSON response matches expected
func ResponseJSONTest(rec *httptest.ResponseRecorder, ctype string, expected string, t *testing.T) {
	// Check content type
	if ctype != "application/json" {
		t.Errorf("content type header does not match: \n\ngot\n\n%v\n\nwant\n\n%v",
			ctype, "application/json")
	}

	// Check body of response
	if rec.Body.String() != expected {
		t.Errorf("handler returned unexpected body. \n\ngot:\n\n%v\n\nwant\n\n%v",
			rec.Body.String(), expected)
	}
}

func TestHandler(t *testing.T) {
	// Request 1: Ensure correct behaviour when searching for "Jesus of Suburbia" 
	// Create a ResponseRecorder 
	rec1 := httptest.NewRecorder()

	// Create a request to pass to our handler. 
	req1, err1 := http.NewRequest(http.MethodGet, "http://localhost4041/?search=jesus%20of%20suburbia", nil)

	ResponseCodeTest(rec1, req1, err1, http.StatusOK, t)

	// Get content type
	ctype1 := rec1.Header().Get("Content-Type");

	// expected body to check against response body
	// Must be four spaces from margin (not tab) to correctly match output from server.go
    expected1 := `[{
    "TrackId": 1134,
    "Name": "Jesus Of Suburbia / City Of The Damned / I Don't Care / Dearly Beloved / Tales Of Another Broken Home",
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
	req2, err2 := http.NewRequest(http.MethodGet, "http://localhost4041/?search=london", nil)
	
	ResponseCodeTest(rec2, req2, err2, http.StatusOK, t)

	// Get content type
	ctype2 := rec2.Header().Get("Content-Type")

	// expected body to check against response body
	// Must be four spaces from margin (not tab) to correctly match output from server.go
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
	req3, err3 := http.NewRequest(http.MethodGet, "http://localhost4041/?search=", nil)
	
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

	// Request 6: Ensure correct behaviour when sending apostrophe for database query
	// Create a ResponseRecorder to record the response.
	rec6 := httptest.NewRecorder()

	// Create a request to pass to our handler. 
	req6, err6 := http.NewRequest(http.MethodGet, "http://localhost4041/?search=please%20don%27t%20touch", nil)
	
	ResponseCodeTest(rec6, req6, err6, http.StatusOK, t)
	
	// Get content type
	ctype3 := rec6.Header().Get("Content-Type")

	// expected body to check against response body
	// Must be four spaces from margin (not tab) to correctly match output from server.go
	expected3 := `[{
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
	
	ResponseJSONTest(rec6, ctype3, expected3, t)
}