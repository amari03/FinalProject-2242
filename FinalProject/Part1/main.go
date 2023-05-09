package main

import (
	"errors"
	"log"
	"net/http"
)

//An HTTP cookie is a small piece of data that a server sends to the user’s web browser.
//The browser can store this data and send it back to the same server, even after the browser restart
//The first thing to keep in mind is that cookies in Go are represented by the http.Cookie type which we created below

func setCookieHandler(w http.ResponseWriter, r *http.Request) { // takes an http.ResponseWriter and an http.Request as parameters.
	// Initialize a new cookie containing the string "Hello world!" and some
	// non-default attributes.
	cookie := http.Cookie{
		Name: "exampleCookie", //cookie name. contains any US-ASCII characters except for charachters below -It is a mandatory field.
		// Characters not allowed: ( ) < > @ , ; : \ " / [ ? ] = { } and space, tab and control characters.

		Value: "Hello world!", //contains the data that you want to used. It is a mandatory field.

		//It can contain any US-ASCII characters except , ; \ " and space, tab and control characters.
		//The rest of the fields are optional and they just map directly to the respective cookie attributes for

		Path: "/", //The cookie is available to all paths on the website.

		MaxAge: 3600, //The cookie will expire after 3600 seconds (1 hour).

		HttpOnly: true, //  The cookie can only be accessed via HTTP(S) requests and cannot be accessed through client-side scripts such as JavaScript.
		//This helps to prevent cross-site scripting (XSS) attacks.

		Secure:   true, //The cookie can only be sent over a secure HTTPS connection.
		SameSite: http.SameSiteLaxMode,

		/*  The cookie will only be sent to the website that set the cookie and any third-party domains that the website
		trusts (as defined by the http.SameSiteLaxMode constant). This helps to prevent cross-site request forgery (CSRF) attacks.*/}

	http.SetCookie(w, &cookie) //sets the cookie in the response headers and sends it to the client's browser.

	/*
	 is a function provided by the net/http package that adds a Set-Cookie header to the response,
	 which tells the client's browser to store the cookie. The first argument to http.SetCookie() is the http.ResponseWriter,
	  which is used to send the response headers back to the client. The second argument is a pointer to the http.Cookie struct
	   that represents the cookie that should be set.
	*/

	w.Write([]byte("cookie set!")) //just a simple way to let the user know that the cookie has been set.
	//The w.Write() function writes the provided byte slice to the response body, which is sent back to the client's browser.
}

func getCookieHandler(w http.ResponseWriter, r *http.Request) { // takes an http.ResponseWriter and an http.Request as parameters.
	// Retrieve the cookie from the request using its name (which in our case is
	// "exampleCookie"). If no matching cookie is found, this will return a
	// http.ErrNoCookie error. We check for this, and return a 400 Bad Request
	// response to the client.
	cookie, err := r.Cookie("exampleCookie")
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			http.Error(w, "cookie not found", http.StatusBadRequest)
		default:
			log.Println(err)
			http.Error(w, "server error", http.StatusInternalServerError)
		}
		return
	}

	// Echo out the cookie value in the response body.
	w.Write([]byte(cookie.Value))
}

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
