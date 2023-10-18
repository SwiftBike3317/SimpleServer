package api

import (
	"encoding/json"
	"log"
	"net/http"
	users "simpleServer/data"
	"simpleServer/types"

	"github.com/gorilla/mux"
)

type APIServer struct {
	listenAddr string
}
type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string
}

func NewAPIServer(listenAddr string) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
	}
}

func WriteJSON(w http.ResponseWriter, status int, v interface{}) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/get-items", makeHTTPHandleFunc(s.getItems)).Methods("GET")

	log.Println("JSON API SERVER RUNNING ON PORT: ", s.listenAddr)
	srv := &http.Server{Addr: s.listenAddr, Handler: router}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("listenAndServe failed: %v", err)
	}
}

func (s *APIServer) getItems(w http.ResponseWriter, r *http.Request) error {

	req := new(types.GetItemsRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	users, isDone := users.GetUsers(req.IDs)
	if !isDone {
		return WriteJSON(w, http.StatusForbidden, ApiError{Error: "user not exist"})
	}

	WriteJSON(w, http.StatusOK, users)

	return nil

}
