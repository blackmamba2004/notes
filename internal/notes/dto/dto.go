package dto

type NoteRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
}

type CreateNoteDTO struct {
	Title       string
	Description *string
}

type UpdateNoteDTO struct {
	Id          int64
	Title       *string
	Description *string
}

type NoteFromDB struct {
	Id          int64
	Title       string
	Description *string
}

type NoteResponse struct {
	Id          int64   `json:"id"`
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
}
