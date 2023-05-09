package main

import (
	"encoding/base64"
	"errors"
	"log"
	"net/http"
)

//variables that contain the errors
var (
	ErrValueTooLong = errors.New("Cookie value too long")
	ErrInvalidValue = errors.New("Invalid cookie value")
)

func main() {
	//multiplexer act lik a router
	mux := http.NewServeMux() // This multiplexer acts like a router and is used to route incoming requests to the appropriate handlers.
	// mux.HandleFunc("/set", setCookieHandler) /Two handler functions, setCookieHandler and getCookieHandler, are registered with the multiplexer using the HandleFunc method.
	mux.HandleFunc("/get", getCookieHandler) // "/= endpoint(url)where you wanna go" key "home = value"
	log.Print("starting server on :4000")    //server port
	err := http.ListenAndServe(":4000", mux) //we create a server
	// the job of the web server listen for request does not process
	// when it gets the request on the right port it will send it to the multiplexer/router and then it will send it to the handler
	log.Fatal(err)
}

//sets a cookie using the http.Cookie struct.
func setCookieHandler(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:     "pepto",
		Value:    "howdy!",
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &cookie) //SetCookie function from the net/http package is used to set the cookie in the HTTP response.
	w.Write([]byte("the cookies have been set!"))

}

func getCookieHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := Read(r, "pepto") //reads a cookie using the Read function, passing the cookie name "pepto" and the request object.

	// If the cookie is not found, an HTTP error with status code 400 is returned with a message indicating that a cookie was not found.
	//Otherwise, the cookie value is returned in the response.
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			http.Error(w, "a cookie was not found", http.StatusBadRequest)
		default:
			log.Println(err)
			http.Error(w, "server error", http.StatusInternalServerError)
		}
		return
	}
	w.Write([]byte(cookie))
}

//a utilty function that validates the cookie parameters and sets the cookie in the HTTP response.
func Write(w http.ResponseWriter, cookie http.Cookie) error {
	cookie.Value = base64.URLEncoding.EncodeToString([]byte(cookie.Value))
	//ensure the cookie is not greater than 4096 bytes
	if len(cookie.String()) > 4096 {
		return ErrValueTooLong
	}
	http.SetCookie(w, &cookie)
	return nil
}

func Read(r *http.Request, name string) (string, error) {
	//read the cookie
	//reads a cookie with a given name from a given HTTP request object
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}
	//It decodes the cookie value using base64 decoding and returns the value as a string.
	//If the cookie is not found or if the value cannot be decoded, an error is returned.

	//let's decode the cookie
	value, err := base64.URLEncoding.DecodeString(cookie.Value)
	if err != nil {
		return "", ErrInvalidValue
	}
	return string(value), nil
}
