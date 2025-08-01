package dto

import (
	"reflect"
	"strings"
)

type NoteRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
}

type NoteDTO struct {
	Title       *string
	Description *string
}

func (n *NoteDTO) Fields() (count int) {
	if n.Title != nil {
		count++
	}
	if n.Description != nil {
		count++
	}
	return
}

// func (n *NoteDTO) toDollar() {

// }

func (n *NoteDTO) ToCortage() (dbFields string) {
	t := reflect.TypeOf(*n)
	for i := 0; i < t.NumField(); i++ {
		dbFields = dbFields + strings.ToLower(t.Field(i).Name) + ", "
	}
	return strings.TrimRight(dbFields, ", ")
}

type NoteFromDB struct {
	Id          *int64
	Title       *string
	Description *string
}

type NoteResponse struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}
