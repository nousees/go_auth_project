# Как работают Dockerfile и docker-compose.yaml

Этот файл подробно объясняет, что делают `Dockerfile` и `docker-compose.yaml` в проекте.

## Dockerfile

`Dockerfile` описывает, как собрать Docker-образ Go-приложения.

В проекте используется multi-stage build: сначала приложение собирается в Go-образе, затем готовый бинарный файл копируется в легкий Alpine-образ.

Полный файл:

```dockerfile
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.* ./
COPY .env ./
COPY config/config.yaml ./config/

RUN go mod download

COPY . .

RUN go build -o main ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/.env .
COPY --from=builder /app/config/config.yaml ./config/

EXPOSE 8080

CMD ["./main"]
```

## Разбор Dockerfile построчно

### `FROM golang:1.23-alpine AS builder`

Берется официальный Docker-образ Go версии `1.23` на базе Alpine Linux.

`AS builder` задает имя стадии сборки. Эта стадия нужна только для компиляции приложения.

### `WORKDIR /app`

Устанавливает рабочую директорию внутри контейнера:

```text
/app
```

Все следующие команды будут выполняться относительно этой папки.

### `COPY go.* ./`

Копирует в контейнер все файлы, которые начинаются с `go.`:

- `go.mod`
- `go.sum`

Эти файлы нужны для установки Go-зависимостей.

### `COPY .env ./`

Копирует файл `.env` в контейнер.

Важно: если файла `.env` нет в корне проекта, сборка образа завершится ошибкой.

Перед сборкой нужно создать `.env`:

```bash
cp .env.example .env
```

Или в PowerShell:

```powershell
Copy-Item .env.example .env
```

### `COPY config/config.yaml ./config/`

Копирует YAML-конфиг приложения в папку:

```text
/app/config/
```

Приложение ожидает конфиг именно по пути:

```text
config/config.yaml
```

### `RUN go mod download`

Скачивает зависимости проекта из `go.mod` и `go.sum`.

Эта команда вынесена до `COPY . .`, чтобы Docker мог кэшировать слой с зависимостями. Если код меняется, но `go.mod` и `go.sum` не меняются, зависимости не скачиваются заново.

### `COPY . .`

Копирует весь проект в контейнер в текущую рабочую директорию `/app`.

После этой команды внутри контейнера будут доступны:

- `cmd/`
- `config/`
- `controllers/`
- `internal/`
- `pkg/`
- `go.mod`
- `go.sum`
- остальные файлы проекта.

### `RUN go build -o main ./cmd/main.go`

Компилирует Go-приложение.

`-o main` означает, что результат сборки будет записан в файл:

```text
main
```

`./cmd/main.go` - точка входа приложения.

### `FROM alpine:latest`

Начинается вторая стадия сборки.

Берется чистый легкий образ Alpine Linux. В нем уже не будет Go-компилятора и исходного кода, если их явно не скопировать.

Это уменьшает размер итогового Docker-образа.

### `WORKDIR /app`

Создает и устанавливает рабочую директорию `/app` во втором контейнерном слое.

### `COPY --from=builder /app/main .`

Копирует скомпилированный бинарный файл `main` из первой стадии `builder` в финальный образ.

Источник:

```text
/app/main
```

Назначение:

```text
/app/main
```

### `COPY --from=builder /app/.env .`

Копирует `.env` из стадии сборки в финальный образ.

Приложение читает `.env` при запуске через:

```go
godotenv.Load("./.env")
```

### `COPY --from=builder /app/config/config.yaml ./config/`

Копирует `config.yaml` из стадии сборки в финальный образ.

В итоге файл оказывается здесь:

```text
/app/config/config.yaml
```

### `EXPOSE 8080`

Документирует, что контейнерное приложение слушает порт `8080`.

Важно: `EXPOSE` сам по себе не публикует порт на хост-машину. Публикация порта делается в `docker-compose.yaml` через `ports`.

### `CMD ["./main"]`

Команда, которая запускается при старте контейнера.

В данном случае запускается скомпилированное Go-приложение:

```text
./main
```

## docker-compose.yaml

`docker-compose.yaml` описывает не один контейнер, а весь набор сервисов проекта:

- `auth_service` - Go-приложение;
- `db` - PostgreSQL;
- `app_network` - внутренняя сеть между контейнерами;
- `pgdata` - volume для данных PostgreSQL.

Полный файл:

```yaml
services:
  auth_service:
    build:
      context: .
      dockerfile: Dockerfile
    ports:
      - "${PORT}:${PORT}"
    volumes:
      - ./config/config.yaml:/app/config/config.yaml
    depends_on:
      - db
    restart: on-failure
    networks:
      - app_network

  db:
    image: postgres:latest
    environment:
      - POSTGRES_USER=${DB_USERNAME}
      - POSTGRES_PASSWORD=${DB_PASSWORD}
      - POSTGRES_DB=${DB_NAME}
    ports:
      - "${DB_EXTERNAL_PORT}:${DB_PORT}"
    volumes:
      - ./pgdata:/var/lib/postgresql/data
    networks:
      - app_network

volumes:
  pgdata:

networks:
  app_network:
    driver: bridge
```

## Разбор docker-compose.yaml построчно

### `services:`

Начинает описание сервисов.

Сервис в Docker Compose - это описание контейнера: из какого образа он создается, какие порты открывает, какие переменные окружения использует и так далее.

### `auth_service:`

Имя сервиса с Go-приложением.

Внутри Docker Compose другие контейнеры могут обращаться к этому сервису по имени `auth_service`.

### `build:`

Указывает, что образ для `auth_service` нужно не скачать, а собрать локально.

### `context: .`

Контекст сборки - текущая папка проекта.

Docker получит доступ к файлам проекта из этой директории.

### `dockerfile: Dockerfile`

Указывает, какой Dockerfile использовать для сборки.

### `ports:`

Описывает проброс портов между хост-машиной и контейнером.

### `- "${PORT}:${PORT}"`

Пробрасывает порт приложения.

Формат:

```text
порт_на_хосте:порт_в_контейнере
```

Если в `.env` указано:

```env
PORT=8080
```

то получится:

```yaml
ports:
  - "8080:8080"
```

После этого приложение доступно на хосте по адресу:

```text
http://localhost:8080
```

### `volumes:`

Описывает монтирование файлов или папок.

### `- ./config/config.yaml:/app/config/config.yaml`

Монтирует локальный файл:

```text
./config/config.yaml
```

в контейнер по пути:

```text
/app/config/config.yaml
```

Это позволяет менять конфиг на хосте без пересборки образа.

### `depends_on:`

Описывает зависимость одного сервиса от другого.

### `- db`

Говорит Docker Compose сначала создать и запустить сервис `db`, а затем `auth_service`.

Важно: `depends_on` не гарантирует, что PostgreSQL уже полностью готов принимать подключения. Он гарантирует только порядок запуска контейнеров.

В этом проекте для `auth_service` указано `restart: on-failure`, поэтому если приложение стартует раньше базы и упадет, Docker попробует перезапустить его.

### `restart: on-failure`

Перезапускает контейнер, если приложение завершилось с ошибкой.

Это полезно при старте вместе с PostgreSQL: база может подниматься чуть дольше приложения.

### `networks:`

Описывает, к каким сетям подключить контейнер.

### `- app_network`

Подключает `auth_service` к сети `app_network`.

В этой сети приложение может обращаться к PostgreSQL по имени сервиса:

```text
db
```

### `db:`

Имя сервиса PostgreSQL.

Для других контейнеров в этой сети имя `db` работает как hostname.

Поэтому для запуска приложения в Docker Compose нужно указывать:

```env
DB_HOST=db
```

### `image: postgres:latest`

Указывает Docker-образ для базы данных.

`postgres:latest` означает последнюю доступную версию PostgreSQL из Docker Hub.

Для production лучше фиксировать конкретную версию, например:

```yaml
image: postgres:16
```

Так поведение будет более предсказуемым.

### `environment:`

Передает переменные окружения внутрь контейнера PostgreSQL.

Официальный образ PostgreSQL использует эти переменные при первом запуске для создания пользователя, пароля и базы данных.

### `- POSTGRES_USER=${DB_USERNAME}`

Создает пользователя PostgreSQL со значением из `.env`.

Например:

```env
DB_USERNAME=postgres
```

### `- POSTGRES_PASSWORD=${DB_PASSWORD}`

Задает пароль пользователя PostgreSQL.

Например:

```env
DB_PASSWORD=135267984
```

### `- POSTGRES_DB=${DB_NAME}`

Создает базу данных с именем из `.env`.

Например:

```env
DB_NAME=go_db
```

### `ports:` для `db`

Описывает проброс порта PostgreSQL наружу.

### `- "${DB_EXTERNAL_PORT}:${DB_PORT}"`

Формат:

```text
порт_на_хосте:порт_в_контейнере
```

Если в `.env` указано:

```env
DB_EXTERNAL_PORT=5436
DB_PORT=5432
```

то получится:

```yaml
ports:
  - "5436:5432"
```

Контейнер PostgreSQL слушает `5432`, а с хост-машины к нему можно подключиться через `localhost:5436`.

### `volumes:` для `db`

Описывает, где хранить данные PostgreSQL.

### `- ./pgdata:/var/lib/postgresql/data`

Монтирует локальную папку:

```text
./pgdata
```

в папку данных PostgreSQL внутри контейнера:

```text
/var/lib/postgresql/data
```

Благодаря этому данные базы сохраняются после остановки контейнера.

### `networks:` для `db`

Подключает PostgreSQL к общей сети.

### `- app_network`

Подключает `db` к сети `app_network`, чтобы `auth_service` мог подключиться к базе.

### `volumes:`

В нижней части файла объявляются named volumes Docker Compose.

### `pgdata:`

Объявляет volume с именем `pgdata`.

В текущем compose-файле фактически используется bind mount `./pgdata:/var/lib/postgresql/data`, а не named volume `pgdata`.

То есть эта секция объявлена, но напрямую в сервисе `db` не используется.

### `networks:`

Объявляет сети Docker Compose.

### `app_network:`

Имя сети, к которой подключаются оба сервиса.

### `driver: bridge`

Используется стандартный Docker bridge-драйвер.

Bridge-сеть позволяет контейнерам общаться друг с другом по имени сервиса внутри одного Docker Compose проекта.

## Как контейнеры общаются между собой

После запуска Docker Compose создает сеть `app_network`.

В этой сети:

- контейнер приложения называется `auth_service`;
- контейнер PostgreSQL называется `db`;
- приложение подключается к базе по hostname `db`;
- PostgreSQL внутри контейнера слушает порт `5432`.

Поэтому строка подключения приложения должна использовать:

```env
DB_HOST=db
DB_PORT=5432
```

## Как порты видны снаружи

Приложение:

```env
PORT=8080
```

Compose пробрасывает:

```text
localhost:8080 -> auth_service:8080
```

PostgreSQL:

```env
DB_EXTERNAL_PORT=5436
DB_PORT=5432
```

Compose пробрасывает:

```text
localhost:5436 -> db:5432
```

## Полезные команды

Собрать и запустить проект:

```bash
docker compose up --build
```

Запустить в фоне:

```bash
docker compose up --build -d
```

Посмотреть логи:

```bash
docker compose logs
```

Посмотреть логи приложения:

```bash
docker compose logs auth_service
```

Посмотреть логи базы:

```bash
docker compose logs db
```

Остановить контейнеры:

```bash
docker compose down
```

Остановить контейнеры и удалить Docker volumes:

```bash
docker compose down -v
```

Пересобрать образ без кэша:

```bash
docker compose build --no-cache
```

## Замечания по текущей Docker-настройке

В текущем `Dockerfile` `.env` копируется внутрь образа. Для учебного проекта это удобно, но для production лучше передавать переменные окружения на этапе запуска контейнера, а не хранить секреты внутри Docker-образа.

Также в `docker-compose.yaml` объявлен named volume `pgdata`, но сервис `db` использует bind mount `./pgdata:/var/lib/postgresql/data`. Это рабочий вариант, просто объявленный named volume сейчас лишний.

Для более стабильной production-конфигурации можно:

- заменить `postgres:latest` на конкретную версию PostgreSQL;
- добавить `healthcheck` для базы данных;
- передавать `.env` в контейнер через `env_file`;
- не копировать `.env` внутрь Docker-образа;
- добавить `.dockerignore`, чтобы не копировать лишние файлы в образ.

