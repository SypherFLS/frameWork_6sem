package db

import (
	"testing"
)

func TestConteiner_InitStorage(t *testing.T) {
	container := &Conteiner{}
	result := container.InitStorage()

	if result == nil {
		t.Fatal("InitStorage вернул nil")
	}

	if result.Items == nil {
		t.Error("Items не должен быть nil")
	}

	if len(result.Items) != 0 {
		t.Errorf("Ожидался пустой слайс, получено: %d элементов", len(result.Items))
	}
}

func TestConteiner_GetAllItems_Empty(t *testing.T) {
	container := &Conteiner{
		Items: make([]Item, 0),
	}

	items := container.GetAllItems()

	if items == nil {
		t.Error("GetAllItems не должен вернуть nil")
	}

	if len(items) != 0 {
		t.Errorf("Ожидалось 0 элементов, получено: %d", len(items))
	}
}

func TestConteiner_GetAllItems_WithItems(t *testing.T) {
	container := &Conteiner{
		Items: []Item{
			{Id: "1", Name: "Item 1", Price: 10.0},
			{Id: "2", Name: "Item 2", Price: 20.0},
			{Id: "3", Name: "Item 3", Price: 30.0},
		},
	}

	items := container.GetAllItems()

	if len(items) != 3 {
		t.Errorf("Ожидалось 3 элемента, получено: %d", len(items))
	}

	if items[0].Id != "1" {
		t.Errorf("Ожидался ID '1', получено: %s", items[0].Id)
	}

	if items[1].Name != "Item 2" {
		t.Errorf("Ожидалось имя 'Item 2', получено: %s", items[1].Name)
	}
}

func TestConteiner_GetItemById_Found(t *testing.T) {
	container := &Conteiner{
		Items: []Item{
			{Id: "1", Name: "Item 1", Price: 10.0},
			{Id: "2", Name: "Item 2", Price: 20.0},
		},
	}

	item, err := container.GetItemById("1")

	if err == nil {
		t.Fatal("GetItemById должен вернуть результат")
	}

	if !err.IsDone {
		t.Error("IsDone должен быть true для найденного элемента")
	}

	if err.Id != 200 {
		t.Errorf("Ожидался код 200, получено: %d", err.Id)
	}

	if item.Id != "1" {
		t.Errorf("Ожидался ID '1', получено: %s", item.Id)
	}

	if item.Name != "Item 1" {
		t.Errorf("Ожидалось имя 'Item 1', получено: %s", item.Name)
	}
}

func TestConteiner_GetItemById_NotFound(t *testing.T) {
	container := &Conteiner{
		Items: []Item{
			{Id: "1", Name: "Item 1", Price: 10.0},
		},
	}

	item, err := container.GetItemById("999")

	if err == nil {
		t.Fatal("GetItemById должен вернуть ошибку")
	}

	if err.IsDone {
		t.Error("IsDone должен быть false для не найденного элемента")
	}

	if err.Id != 404 {
		t.Errorf("Ожидался код 404, получено: %d", err.Id)
	}

	if item.Id != "" {
		t.Errorf("Ожидался пустой ID, получено: %s", item.Id)
	}
}

func TestConteiner_GetItemById_EmptyContainer(t *testing.T) {
	container := &Conteiner{
		Items: make([]Item, 0),
	}

	item, err := container.GetItemById("1")

	if err == nil {
		t.Fatal("GetItemById должен вернуть ошибку")
	}

	if err.IsDone {
		t.Error("IsDone должен быть false")
	}

	if err.Id != 404 {
		t.Errorf("Ожидался код 404, получено: %d", err.Id)
	}

	if item.Id != "" {
		t.Errorf("Ожидался пустой ID, получено: %s", item.Id)
	}
}

func TestConteiner_AddNyItem_Success(t *testing.T) {
	container := &Conteiner{
		Items: make([]Item, 0),
	}

	err := container.AddNyItem("New Item", 15.5)

	if err == nil {
		t.Fatal("AddNyItem должен вернуть результат")
	}

	if !err.IsDone {
		t.Error("IsDone должен быть true для успешного добавления")
	}

	if err.Id != 201 {
		t.Errorf("Ожидался код 201, получено: %d", err.Id)
	}

	if len(container.Items) != 1 {
		t.Errorf("Ожидался 1 элемент, получено: %d", len(container.Items))
	}

	if container.Items[0].Name != "New Item" {
		t.Errorf("Ожидалось имя 'New Item', получено: %s", container.Items[0].Name)
	}

	if container.Items[0].Price != 15.5 {
		t.Errorf("Ожидалась цена 15.5, получено: %f", container.Items[0].Price)
	}

	if container.Items[0].Id != "1" {
		t.Errorf("Ожидался ID '1', получено: %s", container.Items[0].Id)
	}
}

func TestConteiner_AddNyItem_MultipleItems(t *testing.T) {
	container := &Conteiner{
		Items: make([]Item, 0),
	}

	container.AddNyItem("Item 1", 10.0)
	container.AddNyItem("Item 2", 20.0)
	container.AddNyItem("Item 3", 30.0)

	if len(container.Items) != 3 {
		t.Errorf("Ожидалось 3 элемента, получено: %d", len(container.Items))
	}

	if container.Items[0].Id != "1" {
		t.Errorf("Ожидался ID '1', получено: %s", container.Items[0].Id)
	}

	if container.Items[1].Id != "2" {
		t.Errorf("Ожидался ID '2', получено: %s", container.Items[1].Id)
	}

	if container.Items[2].Id != "3" {
		t.Errorf("Ожидался ID '3', получено: %s", container.Items[2].Id)
	}
}

func TestConteiner_AddNyItem_EmptyName(t *testing.T) {
	container := &Conteiner{
		Items: make([]Item, 0),
	}

	err := container.AddNyItem("", 10.0)

	if err == nil {
		t.Fatal("AddNyItem должен вернуть ошибку")
	}

	if err.IsDone {
		t.Error("IsDone должен быть false для ошибки")
	}

	if err.Id != 422 {
		t.Errorf("Ожидался код 422, получено: %d", err.Id)
	}

	if len(container.Items) != 0 {
		t.Errorf("Элемент не должен быть добавлен, получено: %d элементов", len(container.Items))
	}
}

func TestConteiner_AddNyItem_InvalidPrice_Zero(t *testing.T) {
	container := &Conteiner{
		Items: make([]Item, 0),
	}

	err := container.AddNyItem("Item", 0)

	if err == nil {
		t.Fatal("AddNyItem должен вернуть ошибку")
	}

	if err.IsDone {
		t.Error("IsDone должен быть false для ошибки")
	}

	if err.Id != 422 {
		t.Errorf("Ожидался код 422, получено: %d", err.Id)
	}

	if len(container.Items) != 0 {
		t.Errorf("Элемент не должен быть добавлен, получено: %d элементов", len(container.Items))
	}
}

func TestConteiner_AddNyItem_InvalidPrice_Negative(t *testing.T) {
	container := &Conteiner{
		Items: make([]Item, 0),
	}

	err := container.AddNyItem("Item", -5.0)

	if err == nil {
		t.Fatal("AddNyItem должен вернуть ошибку")
	}

	if err.IsDone {
		t.Error("IsDone должен быть false для ошибки")
	}

	if err.Id != 422 {
		t.Errorf("Ожидался код 422, получено: %d", err.Id)
	}

	if len(container.Items) != 0 {
		t.Errorf("Элемент не должен быть добавлен, получено: %d элементов", len(container.Items))
	}
}

func TestConteiner_AddNyItem_EdgeCases(t *testing.T) {
	testCases := []struct {
		name     string
		itemName string
		price    float64
		expected int
		isDone   bool
	}{
		{
			name:     "Минимальная валидная цена",
			itemName: "Item",
			price:    0.01,
			expected: 201,
			isDone:   true,
		},
		{
			name:     "Большая цена",
			itemName: "Expensive Item",
			price:    999999.99,
			expected: 201,
			isDone:   true,
		},
		{
			name:     "Пустое имя и нулевая цена",
			itemName: "",
			price:    0,
			expected: 422,
			isDone:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			container := &Conteiner{
				Items: make([]Item, 0),
			}

			err := container.AddNyItem(tc.itemName, tc.price)

			if err == nil {
				t.Fatal("AddNyItem должен вернуть результат")
			}

			if err.Id != tc.expected {
				t.Errorf("Ожидался код %d, получено: %d", tc.expected, err.Id)
			}

			if err.IsDone != tc.isDone {
				t.Errorf("Ожидался IsDone %v, получено: %v", tc.isDone, err.IsDone)
			}
		})
	}
}

func TestConteiner_Integration(t *testing.T) {
	container := &Conteiner{
		Items: make([]Item, 0),
	}

	err1 := container.AddNyItem("Item 1", 10.0)
	if err1 == nil || !err1.IsDone {
		t.Fatal("Не удалось добавить первый элемент")
	}

	err2 := container.AddNyItem("Item 2", 20.0)
	if err2 == nil || !err2.IsDone {
		t.Fatal("Не удалось добавить второй элемент")
	}

	items := container.GetAllItems()
	if len(items) != 2 {
		t.Errorf("Ожидалось 2 элемента, получено: %d", len(items))
	}

	item, err3 := container.GetItemById("1")
	if err3 == nil || !err3.IsDone {
		t.Fatal("Не удалось найти элемент с ID '1'")
	}

	if item.Name != "Item 1" {
		t.Errorf("Ожидалось имя 'Item 1', получено: %s", item.Name)
	}

	_, err4 := container.GetItemById("999")
	if err4 == nil || err4.IsDone {
		t.Fatal("Не должен был найти элемент с ID '999'")
	}
}
