# Гайд по запуску проекта

Этот файл описывает, как запустить сервис авторизации локально и через Docker Compose.

## Что нужно установить

Для запуска через Docker:

- Docker
- Docker Compose

Для запуска без Docker:

- Go 1.23+
- PostgreSQL

## Быстрый запуск через Docker Compose

### 1. Склонируйте проект

```bash
git clone https://github.com/nousees/go_auth_project.git
cd go_auth_project
```

Если проект уже открыт локально, просто перейдите в папку проекта.

### 2. Создайте файл `.env`

```bash
cp .env.example .env
```

В PowerShell можно выполнить:

```powershell
Copy-Item .env.example .env
```

### 3. Проверьте переменные окружения

Пример `.env` для запуска через Docker Compose:

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

Важно: при запуске через Docker Compose значение `DB_HOST` должно быть `db`, потому что `db` - это имя сервиса PostgreSQL в `docker-compose.yaml`.

Если указать `localhost`, приложение внутри контейнера будет искать базу данных внутри контейнера приложения, а не в контейнере PostgreSQL.

### 4. Запустите проект

```bash
docker compose up --build
```

Команда:

- соберет Docker-образ приложения;
- скачает образ PostgreSQL;
- запустит контейнер с базой данных;
- запустит контейнер с Go-приложением;
- пробросит порт приложения наружу.

После запуска API будет доступно по адресу:

```text
http://localhost:8080
```

Если в `.env` указан другой `PORT`, используйте его.

### 5. Остановка проекта

Остановить контейнеры:

```bash
docker compose down
```

Остановить контейнеры и удалить данные PostgreSQL:

```bash
docker compose down -v
```

В этом проекте данные PostgreSQL также монтируются в локальную папку `pgdata`, поэтому при необходимости ее можно удалить вручную после остановки контейнеров.

## Проверка API

### Регистрация пользователя

Запрос:

```bash
curl -X POST http://localhost:8080/sign-up \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"user@example.com\",\"password\":\"password123\"}"
```

PowerShell-вариант:

```powershell
Invoke-RestMethod `
  -Method Post `
  -Uri "http://localhost:8080/sign-up" `
  -ContentType "application/json" `
  -Body '{"email":"user@example.com","password":"password123"}'
```

Успешный ответ:

```json
{
  "message": "registration successfully",
  "status": "success"
}
```

### Авторизация пользователя

Запрос:

```bash
curl -X POST http://localhost:8080/sign-in \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"user@example.com\",\"password\":\"password123\"}"
```

PowerShell-вариант:

```powershell
Invoke-RestMethod `
  -Method Post `
  -Uri "http://localhost:8080/sign-in" `
  -ContentType "application/json" `
  -Body '{"email":"user@example.com","password":"password123"}'
```

Успешный ответ:

```json
{
  "message": "authorization successfully",
  "status": "success",
  "token": "jwt_token_here"
}
```

## Запуск без Docker

### 1. Поднимите PostgreSQL

Можно использовать локально установленный PostgreSQL или запустить только базу через Docker:

```bash
docker compose up db
```

### 2. Настройте `.env`

Для локального запуска приложения без контейнера значение `DB_HOST` обычно должно быть:

```env
DB_HOST=localhost
```

Если база запущена через Docker Compose и проброшена наружу на порт `5436`, можно использовать:

```env
DB_HOST=localhost
DB_PORT=5436
```

### 3. Установите зависимости

```bash
go mod download
```

### 4. Запустите приложение

```bash
go run ./cmd/main.go
```

Приложение прочитает `.env`, загрузит `config/config.yaml`, подключится к PostgreSQL и запустит HTTP-сервер.

## Частые проблемы

### Dockerfile не находит `.env`

В `Dockerfile` есть строка:

```dockerfile
COPY .env ./
```

Поэтому перед `docker compose up --build` файл `.env` должен существовать в корне проекта.

### Приложение не подключается к базе в Docker

Проверьте, что в `.env` указано:

```env
DB_HOST=db
```

Также можно посмотреть логи:

```bash
docker compose logs auth_service
docker compose logs db
```

### Порт уже занят

Измените `PORT` или `DB_EXTERNAL_PORT` в `.env`.

Например:

```env
PORT=8081
DB_EXTERNAL_PORT=5437
```

После изменения перезапустите проект:

```bash
docker compose down
docker compose up --build
```

