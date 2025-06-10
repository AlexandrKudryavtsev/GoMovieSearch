# GoMovieSearch

### Languages:

- English
- Русский

---

### English 🇬🇧

Required:

1) Docker
2) docker-compose
3) Make

Main commands:

1) `make help` - get help on commands
2) `make compose-up-prod` - start app (for production)
3) `make compose-up-dev` - start only Postgres, RabbitMQ and other
4) `make run-app` - run the application itself (for development), used with `make compose-up-dev`
5) `make install-deps` - install tools (for development)
6) `make migrate-create name="migration_name"` - create migration
7) `make generate-docs` - generate Swagger-documentation
8) `make integration-test` - run integration-test
9) `make mock` - update mocks_test.go file
10) `make test` - start all tests

Documentation: /api/docs/index.html

Structure

```
├── Dockerfile
├── Makefile
├── .env               # for frequently changed or private variables
├── README.md
├── cmd
│   └── app
│       └── main.go    # app launch point 
├── config
│   ├── config.go      # parsing configs
│   └── config.yml     # for other variables
├── docker-compose.yml
├── go.mod
├── go.sum
├── internal           # only the business logic of app
│   ├── app
│   │   ├── app.go     # connects all parts, launching app
│   │   └── migrate.go
│   ├── controller     # all kinds of controllers for HTTP, GRPC, RabbitMQ, etc.
│   ├── entity         # app entities
│   └── usecase        # controller processing
├── migrations         # migrations files
└── pkg                # only auxiliary code
    ├── httpserver
    ├── logger
    └── postgres
```

---

### Русский 🇷🇺

Требуется:

1) Docker
2) docker-compose
3) Make

Основные команды:

1) `make help` - получить справку о командах
2) `make compose-up-prod` - запустить сервис (для продакшена)
3) `make compose-up-dev` - запустить только Postgres, RabbitMQ и т.п.
4) `make run-app` -  запусктить само приложение (для разработки), используется вместе с `make compose-up-dev`
5) `make install-deps` - установить инструменты (для разработки)
6) `make migrate-create name="migration_name"` - создать миграцию
7) `make generate-docs` - сгенерировать Swagger-документацию. 
8) `make integration-test` - запустить интеграционные тесты 
9) `make mock` - обновить mocks_test.go файл
10) `make test` - запустить все тесты

Документация: /api/docs/index.html

Структура:

```
├── Dockerfile
├── Makefile
├── .env               # для частоизменяемых или приватных переменных
├── README.md
├── cmd
│   └── app
│       └── main.go    # точка запуска приложения 
├── config
│   ├── config.go      # парсинг конфигов
│   └── config.yml     # файл для остальных переменных
├── docker-compose.yml
├── go.mod
├── go.sum
├── internal           # только бизнес-логику приложения
│   ├── app
│   │   ├── app.go     # сборка всех частей, запуск приложения
│   │   └── migrate.go
│   ├── controller     # всевозможные контроллеры для HTTP, GRPC, RabbitMQ и т.д.
│   ├── entity         # сущности приложения
│   └── usecase        # обработка контроллеров
├── migrations         # миграции
└── pkg                # только вспомогательный код
    ├── httpserver
    ├── logger
    └── postgres
```
