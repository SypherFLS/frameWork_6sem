package models 

import (
	"framew/internal/log"
)

type Item struct {
	Id string
	Name string 
	Price float64
}

func (i *Item) Validate() *log.Err { 
	if i.Name == "" {
		result := log.MakeError(422, "bad item validation", 1, "empty item's name")
		return result
	} else if i.Price <= 0  {
		result := log.MakeError(422, "bad user validation", 1, "invalid price")
		return result
	} 

	result := log.BackSucsess()
	return result
}

