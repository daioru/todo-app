# 📝 Todo App

Простое и мощное API для управления задачами, написанное на **Go** с использованием **Gin, PostgreSQL, Docker и Swagger**.

---

## 🚀 Функционал

- Регистрация и аутентификация пользователей (JWT)
- Создание, обновление, удаление задач
- Фильтрация задач по пользователю
- Хранение данных в PostgreSQL
- Гибкая система миграций с Goose
- Развёртывание с Docker Compose

---

## 🏰 Технологии

- **Язык**: Go (Gin, sqlx, squirrel, godotenv)
- **База данных**: PostgreSQL
- **Аутентификация**: JWT
- **Миграции**: Goose
- **Логирование**: Zerolog
- **Контейнеризация**: Docker, Docker Compose

---

## 🛠 Установка и запуск

### 🔹 1. Клонирование репозитория
```sh
git clone https://github.com/daioru/todo-app.git
cd todo-app
```

### 🔹 2. Создание `.env` файла  
Создайте файл `.env` в корневой папке проекта:
```sh
JWTSECRET=your_secret_key
```

### 🔹 3. Запуск с Docker
```sh
docker-compose up --build
```
После успешного запуска приложение будет доступно на `http://localhost:8080`.

---

## 📌 API Эндпоинты
Список всех эндпоинтов доступен в **Swagger**:
```
http://localhost:8080/swagger/index.html
```

### 🔹 Регистрация пользователя
```http
POST /api/auth/register
```
**Тело запроса (JSON)**:
```json
{
  "username": "testuser",
  "password": "securepassword"
}
```

### 🔹 Авторизация
```http
POST /api/auth/login
```
**Тело запроса (JSON)**:
```json
{
  "username": "testuser",
  "password": "securepassword"
}
```
**Ответ**:
в качестве ответа устанавливается авторизационный Cookie сроком на 72 часа

### 🔹 Создание задачи (требует Cookie)
```http
POST /api/tasks
```
**Тело запроса (JSON)**:
```json
{
  "title": "Buy groceries",
  "description": "Milk, eggs, bread",
  "status": "pending"
}
```

**Ответ**:
```json
{
  "id": 0,
  "user_id": 0,
  "title": "string",
  "description": "string",
  "status": "string",
  "created_at": "string"
}
```

### 🔹 Получение всех задач пользователя (требует Cookie)
```http
GET /api/tasks
```

**Ответ**:
```json
[
  {
    "id": 0,
    "user_id": 0,
    "title": "string",
    "description": "string",
    "status": "string",
    "created_at": "string"
  }
]
```

### 🔹 Редактирование задачи (требует Cookie)
```http
PUT /api/tasks/{id}
```
**Тело запроса (JSON)**:
обновляются только указанные поля

```json
{
  "title": "Buy groceries",
  "description": "Milk, eggs, bread",
  "status": "pending"
}
```

### 🔹 Удаление задачи (требует Cookie)
```http
DELETE /api/tasks/{id}
```

---

## 🛠 TODO

- [ ] Реализовать фильтрацию задач по статусу
- [ ] Unit тесты

---