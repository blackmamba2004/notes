package repository

import (
	"database/sql"
	"fmt"
	"notes/internal/notes/dto"
	"notes/internal/storage/postgresql"

	sq "github.com/Masterminds/squirrel"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

type NoteRepo struct {
	db *sql.DB
	tx *sql.Tx
}

func New(s *postgresql.Storage, tx *sql.Tx) *NoteRepo {
	return &NoteRepo{db: s.GetDB(), tx: tx}
}

func (r *NoteRepo) Create(n *dto.CreateNoteDTO) (*dto.NoteFromDB, error) {
	const op = "notes.repository.Create"
	var note dto.NoteFromDB

	insertFields := make(map[string]any)
	insertFields["title"] = n.Title
	if n.Description != nil {
		insertFields["description"] = *n.Description
	}

	query, args, err := psql.
		Insert("notes").
		SetMap(insertFields).
		Suffix("RETURNING id, title, description").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	
	err = r.tx.QueryRow(query, args...).Scan(
		&note.Id, &note.Title, &note.Description,
	)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &note, nil
}

func (r *NoteRepo) Read(id int64) (*dto.NoteFromDB, error) {
	const op = "notes.repository.Read"
	var note dto.NoteFromDB

	query := "SELECT id, title, description FROM notes WHERE id=$1"

	err := r.tx.QueryRow(query, id).Scan(
		&note.Id, &note.Title, &note.Description,
	)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &note, nil
}

func (r *NoteRepo) Update(n *dto.UpdateNoteDTO) (*dto.NoteFromDB, error) {
	const op = "notes.repository.Update"
	var note dto.NoteFromDB

	updatedFields := make(map[string]any)
	if n.Title != nil {
		updatedFields["title"] = *n.Title
	}
	if n.Description != nil {
		updatedFields["description"] = *n.Description
	}

	query, args, err := psql.
		Update("notes").
		SetMap(updatedFields).
		Where(sq.Eq{"id": n.Id}).
		Suffix("RETURNING id, title, description").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = r.tx.QueryRow(query, args...).Scan(
		&note.Id, &note.Title, &note.Description,
	)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &note, nil
}

func (r *NoteRepo) Delete(id int64) error {
	const op = "notes.repository.Delete"

	query := "DELETE FROM notes WHERE id=$1"

	_, err := r.tx.Exec(query, id)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}