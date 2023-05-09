package main

import (
	"errors"
	"log"
	"net/http"

	cookies "github.com/Moises/FinalProject/internal"
)

func main() {
	// Start a web server with the two endpoints.
	mux := http.NewServeMux()                // Use the http.NewServeMux() function to create an empty servemux.
	mux.HandleFunc("/set", setCookieHandler) //function to register this with our new servemux, so it acts as the handler for all incoming requests
	//with the URL path stated
	mux.HandleFunc("/get", getCookieHandler)
	log.Print("Listening...")
	// Then we create a new server and start listening for incoming requests
	// with the http.ListenAndServe() function, passing in our servemux for it to
	// match requests against as the second parameter.
	err := http.ListenAndServe(":3000", mux)
	if err != nil {
		log.Fatal(err)
	}
}

func setCookieHandler(w http.ResponseWriter, r *http.Request) {
	// Initialize the cookie as normal.
	cookie := http.Cookie{
		Name:     "exampleCookie",
		Value:    "Hello Zoë!",
		//Hello Zo!
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	// Write the cookie. If there is an error (due to an encoding failure or it
	// being too long) then log the error and send a 500 Internal Server Error
	// response.
	err := cookies.Write(w, cookie)
	if err != nil {
		log.Println(err)
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("cookie set!"))
}

func getCookieHandler(w http.ResponseWriter, r *http.Request) {
	// Use the Read() function to retrieve the cookie value, additionally
	// checking for the ErrInvalidValue error and handling it as necessary.
	value, err := cookies.Read(r, "exampleCookie")
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			http.Error(w, "cookie not found", http.StatusBadRequest)
		case errors.Is(err, cookies.ErrInvalidValue):
			http.Error(w, "invalid cookie", http.StatusBadRequest)
		default:
			log.Println(err)
			http.Error(w, "server error", http.StatusInternalServerError)
		}
		return
	}

	w.Write([]byte(value))
}
