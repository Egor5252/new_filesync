# new_filesync

Лёгкий демонстрационный клиент-серверный синхронизатор файлов на Go (gRPC + protobuf).

## Краткое описание

- Сервер предоставляет gRPC API (`SyncService`) для получения списка файлов, загрузки и скачивания файлов.
- Клиент — CLI на основе `cobra` с командой `pull`, которая получает список файлов с сервера и скачивает отсутствующие.
- Файловая логика (сканирование, запись/чтение, вычисление SHA1) реализована в пакете `inretnal/fs`.

## Что есть в проекте (ключевые файлы)

- `go.mod` — зависимости (gRPC, protobuf, cobra).
- `cmd/server/main.go` — точка входа сервера.
- `cmd/client/main.go` — точка входа клиента (CLI).
- `cmd/client/cli` — реализация CLI (команда `pull`).
- `config/config.json` — конфигурация (порт сервера, `main_path`).
- `proto/sync.proto` — protobuf описание API; сгенерированные файлы `proto/sync.pb.go` и `proto/sync_grpc.pb.go` уже присутствуют.
- `inretnal/server/sync_server.go` — реализация серверных RPC.
- `inretnal/client/sync_client.go` — утилиты клиента для загрузки/скачивания файлов.
- `inretnal/fs` — сканирование директорий, хеширование файлов, запись файлов.

## Требования

- Go 1.24 (в `go.mod` указана версия 1.24.4).
- protoc (если требуется регенерация protobuf).
- (опционально) инструменты генерации для Go:
  - `protoc-gen-go`
  - `protoc-gen-go-grpc`

## Конфигурация

Файл конфигурации: `config/config.json` (пример):

```json
{
  "main_path": "C:\\Users\\...\\папка",
  "server_port": 8000
}
```

`main_path` — корневая папка, где будут храниться файлы клиента.
`server_port` — порт, на котором слушает сервер.

## Сборка и запуск (PowerShell)

Запуск сервера (без сборки бинарника):

```powershell
# из корня репозитория
go run .\cmd\server\
```

Запуск клиентского CLI (пример вывода доступных команд):

```powershell
go run .\cmd\client\
```

Запуск команды `pull` (которую предоставляет CLI) — скачивает отсутствующие файлы из сервера в `main_path`:

```powershell
go run .\cmd\client\ pull
```

Сборка бинарников:

```powershell
go build -o .\bin\server.exe ./cmd/server
go build -o .\bin\client.exe ./cmd/client
# затем запуск
.\bin\server.exe
.\bin\client.exe pull
```

## Regenerate protobuf (опционально)

Если вы правите `proto/sync.proto` и хотите регенерировать Go-файлы, выполните (PowerShell):

```powershell
# установите плагины (один раз)
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.8; go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0

# затем из корня проекта
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto\sync.proto
```

После этого должны появиться/обновиться файлы `proto/sync.pb.go` и `proto/sync_grpc.pb.go`.

## Структура проекта (кратко)

```
cmd/            # точки входа для server и client
config/         # конфиг
inretnal/       # внутренняя логика (fs, client, server) — содержит scan/write
proto/          # protobuf описание и сгенерированные Go файлы
```

## Основные сценарии использования

- Запустить сервер: библиотека слушает порт из `config/config.json`.
- На клиенте выполнить `pull`: клиент подключится к серверу, запросит `ListFiles`, сравнит с локальным `main_path` и скачает отсутствующие файлы.

## Полезные замечания

- В проекте уже присутствуют сгенерированные protobuf Go-файлы, поэтому не обязательно иметь protoc для простого запуска.
- Путь к основным директориям и поведение функций можно посмотреть в `inretnal/fs` и `inretnal/server`.
- Логи и ошибки выводятся в стандартный вывод.