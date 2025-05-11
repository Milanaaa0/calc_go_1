Асинхронный распределённый калькулятор, который позволяет вычислять арифметические выражения.
Основная особенность - вычисления могут выполняться распределённо между несколькими вычислительными узлами (агентами), а пользователи взаимодействуют с оркестратором через HTTP API.

Компоненты системы:
Оркестратор (сервер)

Принимает арифметические выражения от пользователей.

Управляет очередью задач и распределяет их между агентами через gRPC.

Сохраняет результаты вычислений и предоставляет их пользователям.

Реализует HTTP API с регистрацией, авторизацией и историей вычислений.

Агент (вычислитель)

Получает задачи от оркестратора по gRPC.

Выполняет арифметические операции (сложение, вычитание, умножение, деление).

Отправляет результаты обратно оркестратору.

Структура проекта

calc_go_1/
├── calculator/
│   └── calculatorapi/proto/calculator/        # Сгенерированные proto файлы gRPC
├── cmd/
│   ├── orchestrator/
│   │   ├── main.go                            # Запуск HTTP + gRPC сервера
│   │   └── agent/
│   │       └── main.go                        # Запуск gRPC агента
├── internal/
│   ├── agent/
│   │   ├── agent.go
│   │   └── agent_test.go
│   ├── orchestrator/
│   │   ├── handler.go                         # HTTP обработчики
│   │   └── grpc_server.go                     # Реализация gRPC сервера
│   ├── storage/
│   │   └── sqlite.go                          # Работа с SQLite
│   └── user/
│       └── user.go                            # Регистрация, логин, JWT
├── pkg/
│   └── calculation/
│       ├── calculation.go                     # Логика вычислений
│       └── calculation_test.go
├── web/
│   ├── index.html
│   └── style.css
├── go.mod
├── go.sum
└── README.md

Перед началом работы с проектом необходимо убедиться, что на вашем компьютере установлены следующие инструменты и компоненты:

- **Go (версия 1.22.0 или выше)**: Язык программирования Go необходим для разработки и запуска серверной части приложения. Вы можете скачать его с официального сайта [golang.org](https://golang.org/dl/).
- **PostgreSQL**: Система управления базами данных, которая используется в проекте для хранения и обработки данных. Установка PostgreSQL доступна на официальном сайте [postgresql.org](https://www.postgresql.org/download/).
- **Git**: Система контроля версий, необходимая для клонирования проекта из репозитория на GitHub. Скачать Git можно с сайта [git-scm.com](https://git-scm.com/downloads).

Убедитесь, что все эти компоненты установлены и правильно настроены на вашем компьютере перед продолжением работы с проектом.

Клонирование репозитория
bash
```git clone https://github.com/Milanaaa0/calc_go.git```
```cd calc_go_1```
Запуск оркестратора (HTTP + gRPC сервер)
bash
```go run ./cmd/orchestrator/...```
Запуск агента (gRPC клиент)
В отдельном терминале:

bash
```go run ./cmd/orchestrator/agent/...```
Примеры использования с curl
Успешное вычисление выражения
bash
curl --location 'http://localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer <JWT_TOKEN>' \
--data '{
  "expression": "2+2*2"
}'
Ответ:

json
{
  "result": 6
}
Вычисление со скобками
bash
curl --location 'http://localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer <JWT_TOKEN>' \
--data '{
  "expression": "(2+2)*(3+1)"
}'
Ошибка: некорректное выражение
bash
curl --location 'http://localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer <JWT_TOKEN>' \
--data '{
  "expression": "2++2"
}'
Ответ:

json
{
  "error": "invalid expression"
}
Ошибка: деление на ноль
bash
curl --location 'http://localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer <JWT_TOKEN>' \
--data '{
  "expression": "10/0"
}'
Ответ:

json
{
  "error": "division by zero"
}
Регистрация и логин
Регистрация пользователя
bash
curl --location 'http://localhost:8080/api/v1/register' \
--header 'Content-Type: application/json' \
--data '{
  "login": "user1",
  "password": "pass123"
}'
Логин (получение JWT)
bash
curl --location 'http://localhost:8080/api/v1/login' \
--header 'Content-Type: application/json' \
--data '{
  "login": "user1",
  "password": "pass123"
}'
Ответ:

json
{
  "token": "<JWT_TOKEN>"
}
История вычислений
Получить историю вычислений текущего пользователя:

bash
curl --location 'http://localhost:8080/api/v1/history' \
--header 'Authorization: Bearer <JWT_TOKEN>'
Запуск тестов
Выполните команду:

bash
go test -v ./...
Про веб-сервис
API для веб-сервиса пока не реализовано, но вы можете посмотреть дизайн калькулятора в папке web/.
Для удобного запуска используйте расширение Live Server в VS Code:

Установите расширение Live Server

Откройте web/index.html

Нажмите кнопку Go Live в правом нижнем углу VS Code