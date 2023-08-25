package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/AtilioBoher/simple-crud-app/pkg/server/database"
	"github.com/gorilla/mux"
)

// // userinfo is the information that is stored per user.
// type userinfo struct {
// 	email string
// 	age   int
// }

// Server is an HTTP server.
type Server struct {
	ctx context.Context
	db  database.Database
}

// New returns a new server
func New(ctx context.Context, db database.Database) *Server {
	return &Server{
		ctx: ctx,
		db:  db,
	}
}

var indexPage = `
<!DOCTYPE html>
<html>
	<body>
		<h1 style="text-align:center;" > User Database </h1>
		<p style="text-align:center;" > Welcome to the user database. </p>
	</body>
</html>
`

// HandleIndex handles the index path "/".
func (s *Server) HandleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(indexPage))
}

// HandleCreateUsers handles the path "/users/create".
// Create -> Post/Put.
func (s *Server) HandleCreateUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost, http.MethodGet:
		//Check that the input is JSON.
		if contenType := r.Header.Get("Content-Type"); contenType != "application/json" {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Could not read request body: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return

		}
		defer r.Body.Close()

		// Unmarshal the body.
		var u database.User
		err = json.Unmarshal(body, &u)
		if err != nil {
			log.Printf("Could not unmarshal request body: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Validation:
		// 1. User Name should not be empty.
		// 2. User must not exist in order to be created.
		if u.Name == "" {
			log.Printf("Empty username")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		got := s.db.Get(s.ctx, u.Name)
		if got != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("User already exists: %v", u.Name)))
			return
		}

		log.Printf("Create User: %v", u.Name)
		// Write to database.
		err = s.db.Create(s.ctx, u)
		if err != nil {
			log.Printf("Could not create user: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// HandleUsers handles the 'users/{name}' request.
func (s *Server) HandleUsers(w http.ResponseWriter, r *http.Request) {
	//Fetch the name from the query string. Common for all methods of this route.
	params := mux.Vars(r)
	name := params["name"]
	user := s.db.Get(s.ctx, name)
	if user == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		log.Printf("Get user: %v", name)
		
		msg, err := json.Marshal(user)
		if err != nil {
			log.Printf("Could not marshal request: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.Write(msg)
	case http.MethodPatch:
		// Partial updates.
		//Check that the input is JSON.
		if contenType := r.Header.Get("Content-Type"); contenType != "application/json" {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			return
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Could not read request body: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return

		}
		defer r.Body.Close()

		// Unmarshal the body.
		var u database.User
		err = json.Unmarshal(body, &u)
		if err != nil {
			log.Printf("Could not unmarshal request body: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		u.Name = name
		log.Printf("Update user: %v", name)

		user, err := s.db.Update(s.ctx, u)
		if err != nil {
			log.Printf("Could not update database: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// Return updatad value. 
		msg, err := json.Marshal(user)
		if err != nil {
			log.Printf("Could not marshal request: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.Write(msg)

	case http.MethodDelete:
		log.Printf("Delete user: %s", name)
		err := s.db.Delete(s.ctx, name)
		if err != nil {
			log.Printf("Could not delete user: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}
