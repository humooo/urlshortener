# URL Shortener Service

Сервис для сокращения ссылок, написанный на чистом Go (Standard Library).

## Особенности
- Реализован REST API (`POST /shorten`, `GET /go/{id}`).
- **Thread-safe**: Использование `sync.Mutex` для безопасного доступа к данным из горутин.
- **In-memory storage**: Высокая скорость работы (хранение данных в памяти).
- **Idempotency**: Повторная отправка одной и той же ссылки возвращает существующий ключ, не создавая дубликатов.

## Запуск
```bash
go run main.go -port 8080
```

## Примеры использования

### Сократить ссылку
```bash
curl -X POST http://localhost:8080/shorten -d '{"url":"https://avito.ru"}'
# Ответ: {"url":"https://avito.ru", "key":"AbCd12"}
```

### Перейти по ссылке
Откройте в браузере: `http://localhost:8080/go/AbCd12`
