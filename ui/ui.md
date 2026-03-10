# Веб-сервис с конвейером обработки запросов

## 📋 Описание проекта

Веб-сервис на Go для управления каталогом товаров (items) с использованием архитектуры конвейера обработки запросов через Worker Pool. Все HTTP-запросы обрабатываются через систему цепочек заданий (pipeline), где каждая задача получает результат предыдущей.

## 🏗️ Архитектура

### Основные компоненты

1. **Worker Pool** — пул воркеров для параллельной обработки запросов
2. **Task Chain (Pipeline)** — цепочки заданий, где каждая задача получает результат предыдущей
3. **Middleware Chain** — цепочка middleware для обработки запросов:
   - **LoggingMiddleware** — логирование всех запросов и ответов
   - **RecoveryMiddleware** — обработка паник и исключений
   - **TimingMiddleware** — измерение времени выполнения запросов
4. **In-Memory Storage** — хранение данных в памяти процесса

### Схема обработки запроса

```
HTTP Request
    ↓
Middleware Chain (Logging → Recovery → Timing)
    ↓
Handler создает TaskChain
    ↓
TaskChain отправляется в Worker Pool
    ↓
Worker обрабатывает цепочку задач последовательно (pipeline)
    ↓
Результат возвращается через канал
    ↓
HTTP Response
```

## 🚀 Запуск проекта

### Требования

- Go 1.19 или выше
- Git

### Установка и запуск

```bash
# Клонирование репозитория
git clone <repository-url>
cd frWork/backend

# Установка зависимостей
go mod download

# Запуск сервера
go run cmd/main.go
```

Сервер запустится на порту `8080`. Проверить работу можно по адресу: `http://localhost:8080`

## 📡 API Endpoints

### GET `/api/items`
Получить список всех товаров.

**Пример запроса:**
```bash
curl http://localhost:8080/api/items
```

**Пример ответа:**
```json
[
  {
    "Id": "1",
    "Name": "Товар 1",
    "Price": 100.5
  },
  {
    "Id": "2",
    "Name": "Товар 2",
    "Price": 200.0
  }
]
```

### GET `/api/items/{id}`
Получить товар по ID.

**Пример запроса:**
```bash
curl http://localhost:8080/api/items/1
```

**Пример ответа (успех):**
```json
{
  "Id": "1",
  "Name": "Товар 1",
  "Price": 100.5
}
```

**Пример ответа (ошибка):**
```json
{
  "code": 404,
  "success": false,
  "message": "operation",
  "request_id": {
    "num": 1,
    "process": "item with id 999 not found"
  }
}
```

### POST `/api/items`
Создать новый товар.

**Пример запроса:**
```bash
curl -X POST http://localhost:8080/api/items \
  -H "Content-Type: application/json" \
  -d '{
    "Name": "Новый товар",
    "Price": 150.0
  }'
```

**Пример ответа (успех):**
```json
{
  "Id": "1",
  "Name": "Новый товар",
  "Price": 150.0
}
```

**Пример ответа (ошибка валидации):**
```json
{
  "code": 422,
  "success": false,
  "message": "bad item validation",
  "request_id": {
    "num": 1,
    "process": "empty item's name"
  }
}
```

## 🧪 Тестирование

Все тесты находятся на ветке `tests`. Для запуска тестов:

```bash
# Переключение на ветку с тестами
git checkout tests

# Запуск всех тестов
go test ./...

# Запуск тестов конкретного пакета
go test ./internal/workerpool/...
go test ./internal/lib/...
go test ./internal/db/...
go test ./internal/models/...

# Запуск тестов с покрытием
go test -cover ./...

# Запуск тестов с подробным выводом
go test -v ./...
```

### Покрытие тестами

Проект полностью покрыт юнит-тестами:

- ✅ **workerpool/** — тесты для пула воркеров, задач, цепочек и обработки запросов
- ✅ **lib/** — тесты для middleware, логирования и обработки ошибок
- ✅ **db/** — тесты для работы с хранилищем данных
- ✅ **models/** — тесты для моделей данных

## 📁 Структура проекта

```
backend/
├── cmd/
│   └── main.go              # Точка входа приложения
├── internal/
│   ├── db/                   # Хранилище данных
│   │   ├── conteiner.go     # Контейнер для хранения items
│   │   └── item.go          # Модель Item и валидация
│   ├── lib/                  # Библиотека утилит
│   │   ├── error.go         # Обработка ошибок
│   │   ├── log.go           # Логирование
│   │   └── middleware.go    # HTTP middleware
│   ├── models/               # Модели данных
│   │   └── user.go          # Модель User
│   └── workerpool/           # Worker Pool система
│       ├── pool.go           # Пул воркеров
│       ├── task.go           # Задачи и цепочки
│       ├── worker.go         # Воркер
│       └── request.go        # Обработка HTTP запросов
├── logs/                     # Логи приложения
└── go.mod                    # Зависимости Go
```

## 🔧 Особенности реализации

### Worker Pool с Pipeline

Каждый HTTP-запрос обрабатывается через цепочку задач (pipeline):

**GET /api/items:**
```
ParseRequestTask → GetAllItemsTask → WriteResponseTask
```

**POST /api/items:**
```
ParseRequestTask → ValidateItemTask → AddItemTask → WriteResponseTask
```

**GET /api/items/{id}:**
```
ParseRequestTask → GetItemByIdTask → WriteResponseTask
```

### Параллельная обработка

- Пул из 5 воркеров обрабатывает запросы параллельно
- Каждый воркер обрабатывает цепочку задач последовательно
- Данные передаются между задачами в цепочке (pipeline pattern)

### Валидация данных

Товар (Item) валидируется по следующим правилам:
- Имя не должно быть пустым
- Цена должна быть положительным числом (> 0)
- Длина имени не должна превышать 100000 символов

### Обработка ошибок

Единый формат ошибок:
```json
{
  "code": 422,
  "success": false,
  "message": "bad item validation",
  "request_id": {
    "num": 1,
    "process": "empty item's name"
  }
}
```

### Логирование

Все запросы логируются в файлы в директории `logs/`:
- Формат: `app-YYYY-MM-DD.log`
- Содержит: время, метод, путь, статус, длительность, request ID

## 📝 Примеры использования

### Создание товара

```bash
curl -X POST http://localhost:8080/api/items \
  -H "Content-Type: application/json" \
  -d '{"Name": "Ноутбук", "Price": 50000.0}'
```

### Получение всех товаров

```bash
curl http://localhost:8080/api/items
```

### Получение товара по ID

```bash
curl http://localhost:8080/api/items/1
```

### Ошибка валидации

```bash
curl -X POST http://localhost:8080/api/items \
  -H "Content-Type: application/json" \
  -d '{"Name": "", "Price": 100}'
```

Ответ:
```json
{
  "code": 422,
  "success": false,
  "message": "bad item validation",
  "request_id": {
    "num": 1,
    "process": "empty item's name"
  }
}
```

## 🔍 Мониторинг и логи

Логи сохраняются в директории `backend/logs/`:
- Каждый день создается новый файл лога
- Формат: `app-YYYY-MM-DD.log`
- Содержит информацию о всех запросах, ошибках и операциях

Пример записи в логе:
```
[2026-02-24 15:30:45] [INFO] ОПЕРАЦИЯ | REQUEST_START | Method: GET | Path: /api/items | RequestID: req-1234567890-1
[2026-02-24 15:30:45] [INFO] ОПЕРАЦИЯ | REQUEST_END | Method: GET | Path: /api/items | Status: 200 | Duration: 2.5ms | RequestID: req-1234567890-1
```

## 🛡️ Безопасность и производительность

- **Thread-safe операции**: Использование мьютексов для безопасной работы с хранилищем
- **Recovery от паник**: Все паники обрабатываются и возвращаются как HTTP 500 ошибки
- **Параллельная обработка**: Worker pool позволяет обрабатывать несколько запросов одновременно
- **Валидация входных данных**: Все данные валидируются перед обработкой

## 📚 Технологии

- **Go 1.19+** — язык программирования
- **net/http** — HTTP сервер
- **Worker Pool Pattern** — паттерн для параллельной обработки
- **Pipeline Pattern** — паттерн для цепочек задач
- **Middleware Pattern** — паттерн для обработки запросов

## 👨‍💻 Разработка

### Добавление нового endpoint

1. Создать handler функцию в `cmd/main.go`
2. Создать задачи обработки в `internal/workerpool/request.go`
3. Добавить маршрут с middleware в `main()`
4. Добавить тесты в соответствующий `*_test.go` файл

### Добавление новой задачи в pipeline

1. Создать функцию задачи в `internal/workerpool/request.go`:
```go
func NewTaskFunction(data any, mu *sync.Mutex) (any, error) {
    ctx := data.(*RequestContext)
    // Ваша логика
    return ctx, nil
}
```

2. Добавить задачу в цепочку в handler:
```go
tasks := []*workerpool.Task{
    workerpool.NewTask(workerpool.NewTaskFunction, ctx),
    // другие задачи
}
```

## 📄 Лицензия

Проект создан в рамках учебного задания.

## 🔗 Полезные ссылки

- [Go Documentation](https://go.dev/doc/)
- [net/http Package](https://pkg.go.dev/net/http)
- [Worker Pool Pattern](https://gobyexample.com/worker-pools)

---

**Примечание**: Тесты находятся на ветке `tests`. Для запуска тестов переключитесь на эту ветку командой `git checkout tests`.
