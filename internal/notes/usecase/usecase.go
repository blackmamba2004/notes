package usecase

import (
	"notes/internal/notes/dto"
)

type NoteUseCase interface{
	CreateNote(*dto.NoteRequest) (*dto.NoteResponse, error)
	ReadNote(*dto.NoteRequest) (*dto.NoteResponse, error)
	UpdateNote(*dto.NoteRequest) (*dto.NoteResponse, error)
	DeleteNote(*dto.NoteRequest) (string, error)
}
