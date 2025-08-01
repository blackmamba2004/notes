package repository

import (
	"database/sql"
	"fmt"
	"notes/internal/notes/dto"
	"notes/internal/storage/postgresql"
)

type NoteRepo struct {
	db *sql.DB
}

func NewNoteRepo(s *postgresql.Storage) *NoteRepo {
	return &NoteRepo{db: s.GetDB()}
}

func (r *NoteRepo) Create(n *dto.NoteDTO) (*dto.NoteFromDB, error) {
	const op = "storage.postgresql.CreateNote"
	var noteFromDB dto.NoteFromDB

	stmt, err := r.db.Prepare("INSERT INTO notes(title, description) VALUES($1, $2) RETURNING id, title, description")
	if err != nil {
		return &noteFromDB, fmt.Errorf("%s: prepare stmt %w", op, err)
	}

	// var id int64
	// var title string
	// var description string

	// n := NoteDTO{
	// 	Id:          10,
	// 	Title:       "Алгоритмы",
	// 	Description: "Эвклида...",
	// }

	err = stmt.QueryRow(n.Title, n.Description).Scan(
		&noteFromDB.Id, &noteFromDB.Title, &noteFromDB.Description,
	)
	
	if err != nil {
		return &noteFromDB, fmt.Errorf("%s: scan row %w", op, err)
	}

	return &noteFromDB, nil
}


// TODO: поправить
// func (r *NoteRepo) ReadNote(id int) (int64, error) {
// 	const op = "storage.postgresql.ReadNote"
// 	stmt, err := r.db.Prepare(
// 		"SELECT title, description FROM notes WHERE id=$1",
// 	)
// 	if err != nil {
// 		return 0, fmt.Errorf("%s: %w", op, err)
// 	}

// 	err = stmt.QueryRow(id).Scan()

// 	return 1, nil
// }

