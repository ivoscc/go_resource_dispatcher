package main

import (
	"errors"
	"fmt"
	"net/http"
)

type Resource interface {
	Get(w http.ResponseWriter, r *http.Request)
	Post(w http.ResponseWriter, r *http.Request)
	Put(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

// VerifyAPIKey checks if the request header contains "X-Api-Key" and returns
// and error with the appropiate code if not.
// FIXME: Right now only presence of the header is checked.
func VerifyAPIKey(r *http.Request) (error, int) {
	if len(r.Header["X-Api-Key"]) == 0 {
		return errors.New("Must validate with X-Api-Key"), http.StatusUnauthorized
	}
	return nil, 200
}

// Resourcedispatcher bind the resource's methods to the corresponding HTTP verb.
// TODO: should provide a way to disable APIKeyChecking
func ResourceDispatcher(resource Resource) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err, status_code := VerifyAPIKey(r); err != nil {
			http.Error(w, err.Error(), status_code)
			return
		}
		switch method := r.Method; method {
		case "GET":
			resource.Get(w, r)
		case "PUT":
			resource.Put(w, r)
		case "POST":
			resource.Post(w, r)
		case "DELETE":
			resource.Delete(w, r)
		default:
			http.Error(w, "Sorry, I don't know how to handle that :(", http.StatusNotImplemented)
			return
		}
	}
}

// Sample Resource implementation
type SomeResource struct{}

func (word SomeResource) Get(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "GET from SomeResource")
}

func (word SomeResource) Post(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "POST from SomeResource")
}

func (word SomeResource) Put(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "PUT from SomeResource")
}

func (word SomeResource) Delete(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "DELETE from SomeResource")
}

func main() {
	http.HandleFunc("/", ResourceDispatcher(SomeResource{}))
	http.ListenAndServe(":8080", nil)
}
