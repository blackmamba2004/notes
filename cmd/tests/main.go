package main

import (
	"fmt"
	"notes/internal/notes/dto"
	"notes/internal/utils"
)

func main() {
	n := dto.NoteDTO{
		Title: utils.ToPtr("My book"), 
		Description: utils.ToPtr("My book"),
	}
	fmt.Println(n.Description)
	fmt.Println(n.Fields())

	n2 := dto.NoteDTO{
		Description: utils.ToPtr("My book"),
	}
	fmt.Println(n.Description)
	fmt.Println(n2.Fields())

	fmt.Println(n2.ToCortage())
}
