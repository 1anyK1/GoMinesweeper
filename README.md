# Minesweeper TCP Server (Go)

Многопользовательский сетевой сапёр на Go с клиент-серверной архитектурой.

Проект был переписан с C на Go с упором на backend-разработку:

- TCP сервер
- CLI клиент
- несколько клиентов одновременно
- игровые сессии (rooms)
- слоистая архитектура
- Go concurrency
- Docker
- Makefile

---

# Цель проекта

Показать практические навыки разработки backend-сервисов на Go:

- работа с сетью (`net`)
- конкурентная обработка клиентов (`goroutines`)
- организация архитектуры проекта
- работа с состоянием приложения
- разделение логики по слоям
- конфигурация через env
- контейнеризация

---

# Архитектура проекта

```text
.
├── cmd
│   ├── client
│   │   └── client.go
│   └── server
│       └── server.go
├── internal
│   ├── app
│   │   └── server.go
│   ├── config
│   │   └── config.go
│   ├── domain
│   │   ├── cell.go
│   │   ├── game.go
│   │   ├── player.go
│   │   └── session.go
│   ├── logger
│   │   └── logger.go
│   ├── service
│   │   ├── game_service.go
│   │   └── session_service.go
│   ├── storage
│   │   └── memory
│   │       └── session_repo.go
│   └── transport
│       ├── http
│       │   └── handler.go
│       └── tcp
│           ├── client.go
│           ├── client_conn.go
│           ├── client_io.go
│           ├── handler.go
│           ├── protocol.go
│           ├── server.go
│           └── terminal.go
├── Dockerfile
├── Makefile
└── go.mod
```

---

# Используемый стек

- Go
- TCP sockets
- goroutines
- mutex
- log/slog
- Docker
- Makefile

---

# Возможности

## Сервер
- принимает несколько клиентов одновременно
- поддерживает игровые комнаты (sessions)
- хранит состояние игры
- обрабатывает команды игроков
- рассылает обновления участникам комнаты
## Клиент
- CLI интерфейс
- подключение к TCP серверу
- автоматическое обновление поля
- очистка консоли после каждого обновления


# Команды клиента
```
create room1      создать комнату
join room1        войти в комнату
sessions          список комнат
state             показать поле
open x y          открыть клетку
flag x y          поставить / снять флаг
reset             перезапустить игру
list              список подключенных игроков
help              помощь
```

# Запуск локально
## Запуск сервера
```bash
make run-server
```
## Запуск клиента
```bash
make run-client
```
Можно открыть несколько терминалов и подключить несколько клиентов.

# Сборка проекта
```bash
make build
```

Бинарные файлы будут созданы в директории:

```
bin/server
```
```
bin/client
```

# Docker

## Сборка Docker образа
```bash
make docker-build
```

## Запуск сервера в контейнере
```bash
make docker-run
```

После запуска сервер доступен по адресу:
```
localhost:8080
```

# Конфигурация

Используются переменные окружения:
```bash
TCP_ADDR=:8080
GAME_SIZE=8
MINES=10
```

# Возможные улучшения
- PostgreSQL для хранения статистики
- HTTP API (/healthz, /sessions)
- graceful shutdown
- unit tests
- matchmaking
- Web UI / WebSocket клиент
- observability / metrics

# История проекта

Изначально проект был реализован на C с использованием POSIX sockets.

Текущая версия полностью переработана на Go с упором на idiomatic Go backend подходы и современную архитектуру приложения.

Автор
```text
Бочкарёв С.Е
```

Pet-project для практики Go backend разработки.