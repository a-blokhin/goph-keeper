# GophKeeper - Менеджер паролей

GophKeeper - это защищённая клиент-серверная система управления паролями, позволяющая пользователям надёжно и безопасно хранить логины, пароли, бинарные данные (лимит 10МБ) и другую приватную информацию.

## Возможности

- **Защищённое хранение**: Все данные шифруются с использованием AES-256-GCM
- **Множество типов данных**:
  - Пары логин/пароль
  - Текстовые данные
  - Бинарные данные (до 10МБ)
  - Информация о банковских картах
- **Синхронизация**: Мульти-клиентская синхронизация с разрешением конфликтов
- **Аутентификация**: JWT-аутентификация с 24-часовым сроком действия токена
- **Кроссплатформенный CLI**: Доступен для Windows, Linux и macOS
- **gRPC-общение**: Весь обмен данными через gRPC с поддержкой TLS

## Архитектура

Проект следует принципам чистой архитектуры с чётким разделением ответственности:

```
.
├── api/proto/          # Определения gRPC протокола
├── cmd/
│   ├── client/         # CLI клиент
│   └── server/         # gRPC сервер
├── internal/
│   ├── client/         # Реализация клиента
│   ├── crypto/         # Шифрование/дешифрование
│   ├── jwt/            # JWT аутентификация
│   ├── model/          # Модели данных
│   ├── repository/     # Слой доступа к данным
│   ├── server/         # Обработчики gRPC сервера
│   ├── usecase/        # Бизнес-логика
│   └── di/             # Внедрение зависимостей
├── migrations/         # Миграции базы данных
├── tests/e2e/          # E2E тесты
└── docker-compose.yml  # Конфигурация Docker
```

## Требования

- Go 1.21 или выше
- PostgreSQL 15 или выше
- Компилятор Protocol Buffers (protoc) для разработки

## Сборка

### Сборка сервера

```bash
go build -o bin/server ./cmd/server
```

### Сборка клиента

```bash
go build -o goph_keeper ./cmd/client
```

### Сборка для разных платформ

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o goph_keeper-linux-amd64 ./cmd/client
GOOS=linux GOARCH=amd64 go build -o bin/server-linux-amd64 ./cmd/server

# Windows
GOOS=windows GOARCH=amd64 go build -o goph_keeper-windows-amd64.exe ./cmd/client
GOOS=windows GOARCH=amd64 go build -o bin/server-windows-amd64.exe ./cmd/server

# macOS
GOOS=darwin GOARCH=amd64 go build -o goph_keeper-darwin-amd64 ./cmd/client
GOOS=darwin GOARCH=amd64 go build -o bin/server-darwin-amd64 ./cmd/server
```

## Запуск через Docker Compose

Простейший способ запуска GophKeeper - использование Docker Compose:

```bash
docker-compose up -d
```

Это запустит:
- Базу данных PostgreSQL на порту 5432
- gRPC сервер на порту 50051

### Локальный запуск

Для локального запуска вам понадобятся:

1. **Запуск PostgreSQL**:

```bash
docker run -d \
  --name gophkeeper-postgres \
  -e POSTGRES_DB=gophkeeper \
  -e POSTGRES_USER=gophkeeper \
  -e POSTGRES_PASSWORD=gophkeeper_password \
  -p 5432:5432 \
  postgres:15-alpine
```

2. **Сборка клиента и сервера**:

```bash
go build -o goph_keeper ./cmd/client
go build -o bin/server ./cmd/server
```

3. **Добавление клиента в PATH** (опционально):

```bash
# Для Linux/macOS
sudo mv goph_keeper /usr/local/bin/
```

4. **Запуск сервера**:

```bash
./bin/server \
  -address :50051 \
  -dsn "host=localhost port=5432 user=gophkeeper password=gophkeeper_password dbname=gophkeeper sslmode=disable" \
  -jwt-secret "your-secret-key"
```

## Использование CLI

### Регистрация нового пользователя

```bash
goph_keeper register user@example.com password123
```

### Вход в систему

```bash
goph_keeper login user@example.com password123
```

### Управление парами логин/пароль

```bash
# Создание пары логин/пароль
goph_keeper credential create "Мой сайт" "myuser" "mypassword" "website.com"

# Просмотр списка
goph_keeper credential list

# Обновление
goph_keeper credential update <id> "Обновлённый сайт" "newuser" "newpassword" "updated meta"

# Удаление
goph_keeper credential delete <id>
```

### Управление текстовыми данными

```bash
# Создание текстовых данных
goph_keeper text create "Моя заметка" "Это моя секретная заметка" "личное"

# Просмотр списка
goph_keeper text list

# Обновление
goph_keeper text update <id> "Обновлённая заметка" "Обновлённое содержимое" "updated meta"

# Удаление
goph_keeper text delete <id>
```

### Управление бинарными данными

```bash
# Создание бинарных данных
goph_keeper binary create "Мой файл" /path/to/file.bin "документы"

# Просмотр списка
goph_keeper binary list

# Обновление
goph_keeper binary update <id> "Обновлённый файл" /path/to/newfile.bin "updated meta"

# Удаление
goph_keeper binary delete <id>
```

### Управление банковскими картами

```bash
# Создание карты
goph_keeper card create "Моя карта" "1234567890123456" "Иван Иванов" "12/25" "123" "личное"

# Просмотр списка
goph_keeper card list

# Обновление
goph_keeper card update <id> "Обновлённая карта" "9876543210987654" "Петр Петров" "06/26" "456" "updated meta"

# Удаление
goph_keeper card delete <id>
```

### Синхронизация данных

```bash
goph_keeper sync
```

### Просмотр версии

```bash
goph_keeper version
```

## Опции клиента

```bash
goph_keeper -server <address> -tls -tls-cert <path> <command> [args]
```

- `-server`: Адрес сервера (по умолчанию: localhost:50051)
- `-tls`: Включить TLS
- `-tls-cert`: Путь к файлу сертификата TLS (для тестирования)
- `-encryption-key`: 32-байтный ключ шифрования для клиентского шифрования (по умолчанию: предоставляется)

## Тестирование

### Запуск E2E тестов

E2E тесты проверяют полный рабочий процесс через CLI клиент. Тесты автоматически создают тестового пользователя, выполняют все операции и удаляют созданные данные после завершения.

**Требования для запуска E2E тестов:**

1. Запущенный сервер и база данных (через Docker Compose или локально)
2. Собранный клиент `goph_keeper` в корневой директории проекта

**Запуск E2E тестов:**

```bash
# Сначала запустите сервер и базу данных
docker-compose up -d

# Соберите клиент
go build -o goph_keeper ./cmd/client

# Запустите E2E тесты
go test ./tests/e2e/... -v
```

E2E тесты выполняют следующие сценарии:
- Регистрация нового пользователя
- Вход в систему
- Создание, обновление и удаление всех типов данных
- Синхронизация данных
- Очистка созданных данных после завершения тестов

## Безопасность

- Все пароли хешируются с использованием PBKDF2 с 100,000 итераций
- Данные шифруются с использованием AES-256-GCM
- JWT токены истекают через 24 часа
- Поддерживается TLS для защищённого обмена данными
- Доступно клиентское шифрование для дополнительной безопасности

## Разрешение конфликтов

Когда несколько клиентов пытаются обновить одни и те же данные, система использует оптимистическую блокировку на основе версий:

- Если версия клиента совпадает с версией сервера, обновление выполняется успешно
- Если версии различаются, обновление завершается с ошибкой конфликта версий
- Клиент должен получить последние данные перед повторной попыткой обновления

## Схема базы данных

Система использует PostgreSQL со следующими таблицами:

- `users`: Учётные записи пользователей
- `credentials`: Пары логин/пароль
- `text_data`: Текстовые заметки
- `binary_data`: Бинарные файлы
- `cards`: Информация о банковских картах

## Разработка

### Генерация gRPC кода

```bash
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    api/proto/keeper.proto
```
