package api

import (
	"log"
	"net/http"

	"github.com/gonotes/storage"

	"github.com/gorilla/mux"
)

type APIServer struct {
	listenAddress string
	store         storage.Storage
}

func NewAPIServer(listenAddress string, store storage.Storage) *APIServer {
	return &APIServer{
		listenAddress: listenAddress,
		store:         store,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/notes", makeHTTPHandleFunc(s.handleNote))
	router.HandleFunc("/notes/{title}", makeHTTPHandleFunc(s.handleNote))
	router.HandleFunc("/notes/delete/{id}", makeHTTPHandleFunc(s.handleNoteID))
	router.HandleFunc("/notes/title/{title}", makeHTTPHandleFunc(s.handleNoteTitle))

	log.Printf("Server is running on %s\n", s.listenAddress)
	err := http.ListenAndServe(s.listenAddress, router)
	if err != nil {
		log.Fatalf("Could not start server: %v\n", err)
	}
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string `json:"error"`
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}
