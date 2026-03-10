package db

import (
	"testing"
)

func TestItem_Validate_ValidItem(t *testing.T) {
	item := &Item{
		Id:    "1",
		Name:  "Valid Item",
		Price: 10.5,
	}

	err := item.Validate()

	if err == nil {
		t.Fatal("Validate должен вернуть результат")
	}

	if !err.IsDone {
		t.Error("IsDone должен быть true для валидного item")
	}

	if err.Id != 200 {
		t.Errorf("Ожидался код 200, получено: %d", err.Id)
	}
}

func TestItem_Validate_EmptyName(t *testing.T) {
	item := &Item{
		Id:    "1",
		Name:  "",
		Price: 10.5,
	}

	err := item.Validate()

	if err == nil {
		t.Fatal("Validate должен вернуть ошибку для пустого имени")
	}

	if err.IsDone {
		t.Error("IsDone должен быть false для ошибки")
	}

	if err.Id != 422 {
		t.Errorf("Ожидался код 422, получено: %d", err.Id)
	}
}

func TestItem_Validate_InvalidPrice_Zero(t *testing.T) {
	item := &Item{
		Id:    "1",
		Name:  "Test Item",
		Price: 0,
	}

	err := item.Validate()

	if err == nil {
		t.Fatal("Validate должен вернуть ошибку для нулевой цены")
	}

	if err.IsDone {
		t.Error("IsDone должен быть false для ошибки")
	}

	if err.Id != 422 {
		t.Errorf("Ожидался код 422, получено: %d", err.Id)
	}
}

func TestItem_Validate_InvalidPrice_Negative(t *testing.T) {
	item := &Item{
		Id:    "1",
		Name:  "Test Item",
		Price: -5.0,
	}

	err := item.Validate()

	if err == nil {
		t.Fatal("Validate должен вернуть ошибку для отрицательной цены")
	}

	if err.IsDone {
		t.Error("IsDone должен быть false для ошибки")
	}

	if err.Id != 422 {
		t.Errorf("Ожидался код 422, получено: %d", err.Id)
	}
}

func TestItem_Validate_NameTooLong(t *testing.T) {
	longName := make([]byte, 100001)
	for i := range longName {
		longName[i] = 'a'
	}

	item := &Item{
		Id:    "1",
		Name:  string(longName),
		Price: 10.5,
	}

	err := item.Validate()

	if err == nil {
		t.Fatal("Validate должен вернуть ошибку для слишком длинного имени")
	}

	if err.IsDone {
		t.Error("IsDone должен быть false для ошибки")
	}

	if err.Id != 400 {
		t.Errorf("Ожидался код 400, получено: %d", err.Id)
	}
}

func TestItem_Validate_EdgeCases(t *testing.T) {
	testCases := []struct {
		name     string
		item     *Item
		expected int
		isDone   bool
	}{
		{
			name:     "Минимальная валидная цена",
			item:     &Item{Id: "1", Name: "Item", Price: 0.01},
			expected: 200,
			isDone:   true,
		},
		{
			name:     "Максимальная длина имени",
			item:     &Item{Id: "1", Name: string(make([]byte, 100000)), Price: 10.0},
			expected: 200,
			isDone:   true,
		},
		{
			name:     "Пустое имя и нулевая цена",
			item:     &Item{Id: "1", Name: "", Price: 0},
			expected: 422,
			isDone:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.item.Validate()

			if err == nil {
				t.Fatal("Validate должен вернуть результат")
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
