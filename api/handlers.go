package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gonotes/types"
	"github.com/gorilla/mux"
)

func (s *APIServer) handleNote(w http.ResponseWriter, r *http.Request) error {

	switch r.Method {
	case http.MethodGet:
		return s.handleGetNote(w, r)
	case http.MethodPost:
		return s.handleCreateNote(w, r)
	// case http.MethodDelete:
	// 	return s.handleDeleteNote(w, r)
	default:
		return fmt.Errorf("method not allowed: %s", r.Method)
	}
}

func (s *APIServer) handleNoteTitle(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodGet:
		return s.handleGetNoteByTitle(w, r)
	default:
		return fmt.Errorf("method not allowed: %s", r.Method)
	}
}

func (s *APIServer) handleNoteID(w http.ResponseWriter, r *http.Request) error {
	fmt.Println(r.Method)
	switch r.Method {
	case http.MethodDelete:
		return s.handleDeleteNoteByID(w, r)
	default:
		return fmt.Errorf("method not allowed: %s", r.Method)
	}
}

func (s *APIServer) handleGetNote(w http.ResponseWriter, r *http.Request) error {
	title := mux.Vars(r)["title"]

	if title == "" {
		return s.handleGetAllNotes(w, r)
	}

	fmt.Println("GET note with title:", title)

	return WriteJSON(w, http.StatusOK, &types.Note{})
}

func (s *APIServer) handleCreateNote(w http.ResponseWriter, r *http.Request) error {
	// In a real application, you would parse and save the note to storage
	fmt.Println("POST create note")

	createNoteReq := &types.CreateNewNoteRequest{}

	if err := json.NewDecoder(r.Body).Decode(createNoteReq); err != nil {
		return err
	}

	note := types.NewNote(createNoteReq.Title, createNoteReq.Body)
	fmt.Printf("Creating note1: %+v\n", note)
	if err := s.store.CreateNote(note); err != nil {
		return err
	}

	fmt.Printf("Created note2: %+v\n", note)

	return WriteJSON(w, http.StatusCreated, note)
}

func (s *APIServer) handleDeleteNoteByID(w http.ResponseWriter, r *http.Request) error {
	idStr := mux.Vars(r)["id"]

	fmt.Printf("IdStr: %+v\n", idStr)

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid note ID", http.StatusBadRequest)
		return err
	}

	// Delete the note using the store.
	err = s.store.DeleteNoteByID(id)
	if err != nil {
		http.Error(w, "Failed to delete note", http.StatusInternalServerError)
		return err
	}

	// Log the deletion.
	fmt.Println("DELETE note with ID:", id)

	// Respond with a success message.
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Note deleted successfully"))

	return nil
}

func (s *APIServer) handleGetNoteByTitle(w http.ResponseWriter, r *http.Request) error {
	title := mux.Vars(r)["title"]

	note, err := s.store.GetNoteByTitle(title)

	if err != nil {
		return err
	}

	if title == "" {
		return s.handleGetAllNotes(w, r)
	}

	fmt.Println("GET note with title:", title)
	return WriteJSON(w, http.StatusOK, note)
}

func (s *APIServer) handleGetAllNotes(w http.ResponseWriter, r *http.Request) error {
	// In a real application, you would load all notes from storage
	notes, err := s.store.GetAllNotes()

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, notes)
}
