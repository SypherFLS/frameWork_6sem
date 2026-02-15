package db

import (
	"framew/internal/lib"
	"strconv"
)

type Conteiner struct {
	Items []Item
}

func (c *Conteiner) InitStorage() *Conteiner{
	return &Conteiner{
		Items: make([]Item, 0),
	}
}

func (c *Conteiner) GetAllItems() []Item {
	return c.Items
}

func (c *Conteiner) GetItemById(num string) (Item, *lib.Err) {
	var result Item

	for i := 0; i < len(c.Items); i++ {
		if c.Items[i].Id == num {
			result = c.Items[i]
			err := lib.BackSucsess(200, "operation", 4, "founded succesfuly")
			return result, err
		}
	}

	err := lib.MakeError(404, "operation", 1, "item with id "+num+" not found")
	return result, err
}

func (c *Conteiner) AddNyItem(nick string, pr float64) *lib.Err {
	if nick == "" {
		result := lib.MakeError(422, "bad item validation", 1, "empty item's name")
		return result
	} else if pr <= 0 {
		result := lib.MakeError(422, "bad user validation", 1, "invalid price")
		return result
	}

	item :=  Item{
		Id:    strconv.Itoa(len(c.Items) + 1),
		Name:  nick,
		Price: pr,
	}
	c.Items = append(c.Items, item)
	result := lib.BackSucsess(201, "operation", 4, "item created")
	return result
}
