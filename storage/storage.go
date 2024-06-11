package storage

import (
	"database/sql"
	"fmt"

	"github.com/gonotes/types"
	_ "github.com/lib/pq"
)

type Storage interface {
	CreateNote(*types.Note) error
	DeleteNoteByID(int) error
	UpdateNote(*types.Note) error
	GetNoteByID(int) (*types.Note, error)
	GetAllNotes() ([]*types.Note, error)
	GetNoteByTitle(string) (*types.Note, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "user= postgres dbname=postgres password=gonotes sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) Init() error {
	return s.CreateNoteTable()
}

func (s *PostgresStore) CreateNoteTable() error {
	query := `CREATE TABLE IF NOT EXISTS note (
		id SERIAL PRIMARY KEY,
		title TEXT,
		body TEXT,
		created_at TIMESTAMP,
		updated_at TIMESTAMP
	)`

	_, err := s.db.Exec(query)

	return err
}

func (s *PostgresStore) CreateNote(note *types.Note) error {
	query := `INSERT INTO note (title, body, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id`

	err := s.db.QueryRow(query, note.Title, note.Body, note.CreatedAt, note.UpdatedAt).Scan(&note.ID)
	if err != nil {
		return err
	}

	fmt.Printf("New note ID: %d\n", note.ID)
	fmt.Printf("%+v\n", note)

	return nil
}

func (s *PostgresStore) DeleteNoteByID(id int) error {
	query := `DELETE FROM note WHERE id = $1`

	result, err := s.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	// If no rows were affected, it means the note was not found
	if rowsAffected == 0 {
		return fmt.Errorf("no note found with id %d", id)
	}

	return nil
}

func (s *PostgresStore) UpdateNote(note *types.Note) error {
	return nil
}

func (s *PostgresStore) GetNoteByID(id int) (*types.Note, error) {
	return nil, nil
}

func (s *PostgresStore) GetNoteByTitle(title string) (*types.Note, error) {
	query := `SELECT id, title, body, created_at, updated_at FROM note WHERE title = $1`

	note := &types.Note{}
	err := s.db.QueryRow(query, title).Scan(&note.ID, &note.Title, &note.Body, &note.CreatedAt, &note.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return note, nil
}

// getAllNotes returns all notes from the database
func (s *PostgresStore) GetAllNotes() ([]*types.Note, error) {
	// define the query
	query := `select id, title, body, created_at, updated_at from note`

	// execute the query
	rows, err := s.db.Query(query)

	if err != nil {
		return nil, err
	}

	// close the rows when the function returns
	defer rows.Close()

	// return the list of notes::
	notes := []*types.Note{}

	for rows.Next() {
		note := &types.Note{}
		err := rows.Scan(&note.ID, &note.Title, &note.Body, &note.CreatedAt, &note.UpdatedAt)
		if err != nil {
			return nil, err
		}

		notes = append(notes, note)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return notes, nil
}
