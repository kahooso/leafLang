
# 📘 leafLang — приложение для изучения иностранных слов

**leafLang** — это веб-приложение для изучения иностранных слов с использованием метода **интервального повторения** (Spaced Repetition). Приложение разработано на языке Go с использованием фреймворка Gin и базой данных PostgreSQL. Оно подходит для студентов, школьников, самоучек и всех, кто хочет эффективно запоминать иностранные слова.

---

## 🔧 Возможности

- 📘 Ведение персонального словаря
- 🧠 Интервальное повторение слов
- 🧾 Практики с учётом статистики
- 📊 Отслеживание прогресса

---

## ⚙️ Технологии

| Компонент            | Технология                         |
|----------------------|------------------------------------|
| Язык программирования | Go (1.20+)                         |
| Веб-фреймворк         | Gin                                |
| ORM-библиотека        | GORM                               |
| База данных           | PostgreSQL                         |
| Аутентификация        | JWT (в cookie)                     |
| Шаблоны               | Go Templates (HTML)                |
| Клиентская часть      | HTML, CSS, JavaScript              |
| Контейнеризация       | Docker, Docker Compose             |
| Архитектура           | MVC (монолит)                      |

---

## 📁 Структура проекта

```plaintext
leafLang/
├── cmd/                 # main.go (точка входа)
├── internal/
│   ├── config/          # Загрузка конфигурации из .env
│   ├── database/        # Подключение к БД, миграции
│   ├── handlers/        # Обработчики HTTP-запросов
│   ├── models/          # Структуры БД и бизнес-логика
│   └── routes/          # Маршруты API и HTML
├── pkg/
│   └── middleware/      # JWT-аутентификация, логгирование
├── templates/           # HTML-шаблоны
├── static/              # CSS, JS, изображения
├── sqlbackup.sql        # Скрипт создания и наполнения БД
├── .env                 # Конфигурация среды
├── docker-compose.yml   # Запуск через Docker
└── README.md
```

---

## 🚀 Запуск проекта

### 🐳 Запуск с Docker

1. Установите Docker и Docker Compose.
2. Поместите `sqlbackup.sql` рядом с `docker-compose.yml`.
3. Выполните команду:

```bash
docker-compose up --build
```

4. Приложение будет доступно на [http://localhost:8080](http://localhost:8080)

5. Если БД не инициализировалась автоматически:

```bash
docker exec -i leaflang-db psql -U postgres -d leaflang_db < sqlbackup.sql
```

### 💻 Ручной запуск (без Docker)

1. Установите Go и PostgreSQL.
2. Создайте базу данных:

```bash
createdb leaflang_db
psql -U postgres -d leaflang_db < sqlbackup.sql
```

3. Создайте файл `.env`:

```env
DB_USER=postgres
DB_PASS=your_password
DB_NAME=leaflang_db
DB_HOST=localhost
DB_PORT=5432
JWT_SECRET=your_secret_key
```

4. Установите зависимости и запустите проект:

```bash
go mod tidy
cd cmd/leaflang
go run main.go
```

Или соберите бинарный файл:

```bash
go build -o leaflang ./cmd/leaflang
./leaflang
```

---

## 📌 Основные маршруты API

| Метод | URI                    | Назначение                          |
|-------|------------------------|-------------------------------------|
| POST  | `/auth/register`       | Регистрация пользователя            |
| POST  | `/auth/login`          | Авторизация                         |
| POST  | `/auth/logout`         | Выход из системы                    |
| GET   | `/words`               | Получить список слов                |
| POST  | `/words`               | Добавить слово                      |
| PUT   | `/words/:id`           | Обновить слово                      |
| DELETE| `/words/:id`           | Удалить слово                       |
| GET   | `/practice`            | Получить задания на повторение      |
| POST  | `/practice/submit`     | Отправить результаты практики       |
| GET   | `/profile`             | Получить информацию о пользователе |
| GET   | `/stats`               | Статистика изучения                 |

---

## 📥 Импорт базы данных

Файл `sqlbackup.sql` содержит SQL-скрипт, который:

- Создаёт структуру таблиц
- Заполняет БД начальными данными

### Пример импорта вручную:

```bash
psql -U postgres -d leaflang_db < sqlbackup.sql
```

Или внутри Docker-контейнера:

```bash
docker exec -i leaflang-db psql -U postgres -d leaflang_db < sqlbackup.sql
```

---

## 🖼 Интерфейс

Приложение использует Go Templates и серверный рендеринг.  
Интерфейс включает:

- форму регистрации/входа
- список слов
- карточки для практики
- страницу статистики
