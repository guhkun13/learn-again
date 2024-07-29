package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// User represents a user entity
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UserHandler implements the Handler interface
type UserHandler struct {
	users map[int]User
}

// NewUserHandler creates a new UserHandler
func NewUserHandler() *UserHandler {
	return &UserHandler{
		users: make(map[int]User),
	}
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user.ID = len(h.users) + 1
	h.users[user.ID] = user

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) Read(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, ok := h.users[id]
	if !ok {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user.ID = id
	h.users[id] = user

	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if _, ok := h.users[id]; !ok {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	delete(h.users, id)
	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	users := make([]User, 0, len(h.users))
	for _, user := range h.users {
		users = append(users, user)
	}

	json.NewEncoder(w).Encode(users)
}
