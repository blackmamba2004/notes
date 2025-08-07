package tests

import (
	"database/sql"
	"fmt"
	"notes/internal/notes/dto"
	"notes/internal/notes/repository"
	"notes/internal/storage/postgresql"
	"notes/internal/utils"
	"testing"
)

func initStorageAndTx() (*postgresql.Storage, *sql.Tx, error) {
	const op = "tests.repository_test.initStorageAndTx"
	postgresDSN := fmt.Sprintf(
		"%s://%s:%s@%s:%d/%s", "postgres", "postgres", "passwd", "localhost", 5441, "app",
	)

	storage, err := postgresql.New(postgresDSN)

	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}

	tx, err := storage.BeginTransaction()
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}
	return storage, tx, nil
}

func TestCreateUniqueNote(t *testing.T) {
	const op = "tests.repository_test.TestCreateUniqueNote"
	storage, tx, err := initStorageAndTx()

	if err != nil {
		t.Error(err)
	}

	defer storage.Close()
	defer tx.Rollback()

	noteRepo := repository.New(storage, tx)
	
	notesDTO := []*dto.CreateNoteDTO{
		{
			Title: "New Note",
			Description: utils.ToPtr("New Description"),
		},
		{
			Title: "",
			Description: utils.ToPtr("New Description"),
		},
		{
			Title: "About journal",
		},
	}
	for _, noteDTO := range notesDTO{
		t.Run("Create Note", func(t *testing.T) {
			createNoteFromDB, err := noteRepo.Create(noteDTO)
			if err != nil {
				err = fmt.Errorf("%s: %w", op, err)
				t.Error(err)
			}
			if createNoteFromDB.Title != noteDTO.Title {
				t.Errorf("Note's title is %s; want %s", createNoteFromDB.Title, noteDTO.Title)
			}
			if (createNoteFromDB.Description != nil) && (*createNoteFromDB.Description != *noteDTO.Description){
				t.Errorf("Note's description is %s; want %s", *createNoteFromDB.Description, *noteDTO.Description)
			}
		})
	}
}

func TestReadNote(t *testing.T) {
	const op = "tests.repository_test.TestReadNote"
	storage, tx, err := initStorageAndTx()

	if err != nil {
		t.Error(err)
	}

	defer storage.Close()
	defer tx.Rollback()

	noteRepo := repository.New(storage, tx)

	createNoteDTO := &dto.CreateNoteDTO{
		Title: "Any Note",
		Description: utils.ToPtr("some description"),
	}

	createNoteFromDB, err := noteRepo.Create(createNoteDTO)
	if err != nil {
		err = fmt.Errorf("%s: %w", op, err)
		t.Error(err)
	}

	readNote, err := noteRepo.Read(createNoteFromDB.Id)
	if err != nil {
		err = fmt.Errorf("%s: %w", op, err)
		t.Error(err)
	}

	if readNote.Id != createNoteFromDB.Id {
		t.Errorf("Note's id is %d; want %d", readNote.Id, createNoteFromDB.Id)
	}
}

func TestUpdateNote(t *testing.T) {
	const op = "tests.repository_test.TestUpdateNote"
	storage, tx, err := initStorageAndTx()

	if err != nil {
		t.Error(err)
	}

	defer storage.Close()
	defer tx.Rollback()

	noteRepo := repository.New(storage, tx)

	createNoteDTO := &dto.CreateNoteDTO{
		Title: "Any Note",
		Description: utils.ToPtr("some description"),
	}

	createNoteFromDB, err := noteRepo.Create(createNoteDTO)
	if err != nil {
		err = fmt.Errorf("%s: %w", op, err)
		t.Error(err)
	}

	updateNoteDTO := &dto.UpdateNoteDTO{
		Id: createNoteFromDB.Id,
		Title: utils.ToPtr("New Title"),
		Description: utils.ToPtr("New description"),
	}

	updateNoteFromDB, err := noteRepo.Update(updateNoteDTO)
	if err != nil {
		err = fmt.Errorf("%s: %w", op, err)
		t.Error(err)
	}

	if updateNoteFromDB.Id != createNoteFromDB.Id {
		t.Errorf("Note's id is %d; want %d", updateNoteFromDB.Id, createNoteFromDB.Id)
	}
	if updateNoteFromDB.Title != *updateNoteDTO.Title {
		t.Errorf("Note's title is %s; want %s", updateNoteFromDB.Title, *updateNoteDTO.Title)
	}
	if *updateNoteFromDB.Description != *updateNoteDTO.Description {
		t.Errorf("Note's description is %s; want %s", *updateNoteFromDB.Description, *updateNoteDTO.Description)
	}
}

func TestDeleteNote(t *testing.T) {
		const op = "tests.repository_test.TestDeleteNote"
	storage, tx, err := initStorageAndTx()

	if err != nil {
		t.Error(err)
	}

	defer storage.Close()
	defer tx.Rollback()

	noteRepo := repository.New(storage, tx)

	createNoteDTO := &dto.CreateNoteDTO{
		Title: "Any Note",
		Description: utils.ToPtr("some description"),
	}

	createNoteFromDB, err := noteRepo.Create(createNoteDTO)
	if err != nil {
		err = fmt.Errorf("%s: %w", op, err)
		t.Error(err)
	}

	err = noteRepo.Delete(createNoteFromDB.Id)
	if err != nil {
		err = fmt.Errorf("%s: %w", op, err)
		t.Error(err)
	}

	_, err = noteRepo.Read(createNoteFromDB.Id)
	if err != nil {
		err = fmt.Errorf("%s: %w", op, err)
		t.Log(err)
	}
}