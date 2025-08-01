package interfaces

import (
	"notes/internal/notes/dto"
)

type NoteRepo interface {
	Create(dto.NoteDTO) (dto.NoteDTO)
}