# Как работает приложение

Проект представляет собой небольшой сервис авторизации на Go. Он позволяет зарегистрировать пользователя по email и паролю, а затем авторизовать его и получить JWT-токен.

## Основные возможности

Сервис предоставляет два HTTP-эндпоинта:

- `POST /sign-up` - регистрация пользователя;
- `POST /sign-in` - авторизация пользователя и выдача JWT-токена.

## Технологии

В проекте используются:

- Go 1.23.2 - основной язык приложения;
- Gin - HTTP-фреймворк для роутинга и обработки запросов;
- PostgreSQL - база данных;
- GORM - ORM для работы с PostgreSQL;
- Viper - чтение YAML-конфига;
- godotenv - загрузка переменных окружения из `.env`;
- go-playground/validator - валидация входных данных;
- Argon2id - хэширование паролей;
- golang-jwt/jwt/v5 - генерация и парсинг JWT;
- Docker и Docker Compose - контейнеризация приложения и базы данных.

## Структура проекта

```text
cmd/
  main.go
config/
  config.go
  config.yaml
controllers/
  signin_controller.go
  signup_controller.go
internal/
  database/
    psql.go
  entities/
    user/
      user.go
      user_signin.go
      user_signup.go
  repository/
    user_repository.go
  usecases/
    signin_usecase.go
    signup_usecase.go
pkg/
  hash/
    hash.go
  jwt/
    jwt.go
Dockerfile
docker-compose.yaml
.env.example
```

## Общий поток запуска

При старте приложение выполняет следующие действия:

1. Загружает конфигурацию через `config.LoadConfig()`.
2. Читает переменные окружения из `.env`.
3. Читает файл `config/config.yaml`.
4. Подставляет значения переменных окружения в YAML-конфиг.
5. Подключается к PostgreSQL через GORM.
6. Выполняет `AutoMigrate` для таблицы пользователей.
7. Создает repository, usecase и controller-слои.
8. Регистрирует HTTP-маршруты.
9. Запускает Gin-сервер на порту из конфига.

## Конфигурация

Файл `.env` хранит реальные значения переменных:

```env
DB_HOST=db
DB_PORT=5432
DB_EXTERNAL_PORT=5436
DB_USERNAME=postgres
DB_NAME=go_db
DB_SSLMODE=disable
DB_PASSWORD=135267984

PORT=8080

SECRET_KEY=mysecret
```

Файл `config/config.yaml` хранит структуру конфигурации:

```yaml
postgres:
  db_host: "${DB_HOST}"
  db_port: "${DB_PORT}"
  db_external_port: "${DB_EXTERNAL_PORT}"
  db_username: "${DB_USERNAME}"
  db_name: "${DB_NAME}"
  db_sslmode: "${DB_SSLMODE}"
  db_password: "${DB_PASSWORD}"

server:
  port: "${PORT}"
```

Код в `config/config.go` загружает `.env`, затем читает `config.yaml`, проходит по всем ключам Viper и заменяет строки вида `${DB_HOST}` на реальные значения из окружения.

## База данных

Подключение к PostgreSQL находится в `internal/database/psql.go`.

Строка подключения собирается из конфига:

```text
host=<host> port=<port> user=<user> dbname=<db_name> sslmode=<sslmode> password=<password>
```

После подключения вызывается:

```go
db.AutoMigrate(&entities.User{})
```

Это значит, что GORM автоматически создает или обновляет таблицу пользователей под структуру `User`.

## Модель пользователя

Основная сущность находится в `internal/entities/user/user.go`:

```go
type User struct {
    ID       int64  `json:"id" gorm:"primaryKey:autoIncrement"`
    Email    string `json:"email" gorm:"not null;uniqueIndex"`
    Password string `json:"password" gorm:"not null"`
}
```

Поля:

- `ID` - первичный ключ;
- `Email` - email пользователя, обязательный и уникальный;
- `Password` - хэш пароля, обязательное поле.

## Регистрация

Регистрация проходит через эндпоинт:

```text
POST /sign-up
```

Тело запроса:

```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

Входная структура находится в `internal/entities/user/user_signup.go`:

```go
type SignUpInput struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
}
```

Что происходит при регистрации:

1. Gin читает JSON из тела запроса.
2. Validator проверяет, что email заполнен и похож на email.
3. Validator проверяет, что пароль заполнен и содержит минимум 8 символов.
4. Usecase хэширует пароль через Argon2id.
5. Repository сохраняет пользователя в PostgreSQL.
6. Если email уже существует, база вернет ошибку уникальности.

Успешный ответ:

```json
{
  "status": "success",
  "message": "registration successfully"
}
```

## Авторизация

Авторизация проходит через эндпоинт:

```text
POST /sign-in
```

Тело запроса:

```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

Что происходит при авторизации:

1. Gin читает JSON из тела запроса.
2. Usecase ищет пользователя по email.
3. Если пользователь не найден, возвращается ошибка.
4. Пароль из запроса сравнивается с хэшем из базы.
5. Если пароль неверный, возвращается ошибка.
6. Если данные корректные, генерируется JWT-токен.

Успешный ответ:

```json
{
  "status": "success",
  "message": "authorization successfully",
  "token": "jwt_token_here"
}
```

## Хэширование паролей

Файл `pkg/hash/hash.go` отвечает за работу с паролями.

При регистрации вызывается:

```go
argon2id.CreateHash(password, argon2id.DefaultParams)
```

В базу данных сохраняется не исходный пароль, а Argon2id-хэш.

При авторизации вызывается:

```go
argon2id.ComparePasswordAndHash(plainPwd, hashedPwd)
```

Это безопаснее, чем хранить пароли в открытом виде.

## JWT

Файл `pkg/jwt/jwt.go` отвечает за генерацию и парсинг токенов.

При успешной авторизации создается JWT с такими claims:

```json
{
  "user_id": 1,
  "exp": 1710000000
}
```

Поле `user_id` хранит ID пользователя.

Поле `exp` задает срок жизни токена. Сейчас токен живет 24 часа:

```go
time.Now().Add(time.Hour * 24).Unix()
```

Токен подписывается алгоритмом `HS256` и секретом из переменной окружения:

```env
SECRET_KEY=mysecret
```

## Слои приложения

### Controller

Файлы:

- `controllers/signup_controller.go`
- `controllers/signin_controller.go`

Controller отвечает за HTTP:

- принять запрос;
- распарсить JSON;
- вернуть HTTP-ответ;
- передать бизнес-логику в usecase.

### Usecase

Файлы:

- `internal/usecases/signup_usecase.go`
- `internal/usecases/signin_usecase.go`

Usecase отвечает за бизнес-логику:

- хэширование пароля при регистрации;
- поиск пользователя при авторизации;
- проверка пароля;
- генерация JWT.

### Repository

Файл:

- `internal/repository/user_repository.go`

Repository отвечает за работу с базой данных:

- создать пользователя;
- найти пользователя по email.

### Entity

Папка:

- `internal/entities/user`

Entity описывает структуры данных:

- пользователь в базе;
- входные данные регистрации;
- входные данные авторизации.

## Ошибки и HTTP-статусы

Регистрация:

- `400 Bad Request` - некорректный JSON или ошибка валидации;
- `409 Conflict` - ошибка регистрации, например email уже существует;
- `200 OK` - пользователь успешно создан.

Авторизация:

- `400 Bad Request` - некорректный JSON;
- `401 Unauthorized` - пользователь не найден или пароль неверный;
- `200 OK` - авторизация успешна, токен выдан.

## Что можно улучшить

В текущем виде сервис уже выполняет базовую авторизацию, но для production можно добавить:

- middleware для проверки JWT на защищенных маршрутах;
- refresh-токены;
- logout и blacklist токенов;
- миграции через отдельный инструмент вместо `AutoMigrate`;
- более строгую политику паролей;
- rate limit на `/sign-in`;
- скрытие подробных ошибок авторизации;
- тесты для usecase, repository и handlers;
- healthcheck для Docker Compose.

