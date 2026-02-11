package db

import (
	"framew/internal/log"
	"framew/internal/models"
	"strconv"
)

type Conteiner struct {
	Items []models.Item
}

func (c *Conteiner) GetAllItems() []models.Item {
	return c.Items
}

func (c *Conteiner) GetItemById(num string) (models.Item, *log.Err) {
	var result models.Item

	for i := 0; i < len(c.Items); i++ {
		if c.Items[i].Id == num {
			result = c.Items[i]
			err := log.BackSucsess()
			return result, err
		}
	}

	err := log.MakeError(404, "item not found", 1, "item with id "+num+" not found")
	return result, err
}

func (c *Conteiner) AddNyItem(nick string, pr float64) *log.Err {
	if nick == "" {
		result := log.MakeError(422, "bad item validation", 1, "empty item's name")
		return result
	} else if pr <= 0 {
		result := log.MakeError(422, "bad user validation", 1, "invalid price")
		return result
	}

	item := models.Item{
		Id:    strconv.Itoa(len(c.Items) + 1),
		Name:  nick,
		Price: pr,
	}
	c.Items = append(c.Items, item)
	result := log.BackSucsessCreate()
	return result
}
