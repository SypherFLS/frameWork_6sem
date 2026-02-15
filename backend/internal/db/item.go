package db

import (
	"framew/internal/lib"
)

type Item struct {
	Id string
	Name string 
	Price float64
}

func (i *Item) Validate() *lib.Err { 
	if i.Name == "" {
		result := lib.MakeError(422, "bad item validation", 1, "empty item's name")
		return result
	} else if i.Price <= 0  {
		result := lib.MakeError(422, "bad user validation", 1, "invalid price")
		return result
	} 

	result := lib.BackSucsess(200, "validation", 2, "item validated")
	return result
}

