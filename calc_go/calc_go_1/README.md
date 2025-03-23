# calc_go

Этот проект представляет собой асинхронный калькулятор, который позволяет вычислять арифметические выражения. Основная особенность системы заключается в том, что вычисления выполняются асинхронно и могут быть распределены между несколькими вычислительными узлами (агентами)
Система состоит из двух основных компонентов:

Оркестратор (Сервер):

1.Принимает арифметические выражения от пользователей.
2.Разбивает выражение на отдельные задачи (например, 2 + 2, 4 * 3).
3.Управляет очередью задач и распределяет их между агентами.
4.Сохраняет результаты вычислений и предоставляет их пользователям по запросу.

Агент (Вычислитель):
1.Получает задачи от оркестратора.
2.Выполняет арифметические операции (сложение, вычитание, умножение, деление).
3.Отправляет результаты обратно оркестратору.

Чтобы запустить программу:
1. Скопируйте репозиторий git clone https://github.com/Milanaaa0/calc_go.git
2. Запуск оркестратора:
Перейдите в директорию cmd/orchestrator и выполните команду:
go run main.go
 
Запуск агента:
Перейдите в директорию cmd/agent и выполните команду:
go run main.go

## Структура проекта
```calc.go_1/
├── cmd/
│ ├── orchestrator/
│ │ └── main.go
│ └── agent/
│ └── main.go
│
├── internal/
│ ├── agent/
│ │ ├── agent_test.go
│ │ └── agent.go
│ └── orchestrator/
│ ├── handler.go
│ └── service.go
│
├── pkg/
│ └── calculation/
│ ├── calculation_test.go
│ ├── calculation.go
│ └── errors.go
│
├── web/
│ ├── index.html
│ └── style.css
│
├── go.mod
├── go.sum
└── README.md
```
Примеры 
Отправьте POST-запрос на эндпоинт /api/v1/calculate с выражением для вычисления:

# Успешное вычисление(200)
curl --location 'localhost/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "2+2*2"
}'


Ответ:
{
  "result": 6
}

# Вычисление нескольких скобок
curl --location 'localhost/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "(2+2)*(3+1)"
}'

# Ошибка: Некорректное выражение(422)

curl --location 'localhost/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "2++2"
}'


Ответ:
{
  "error": "invalid expression"
}

# Ошибка: деление на ноль

curl --location 'localhost/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "10/0"
}'
Ответ:
{
  "error": "division by zero"
}

Чтобы запустить тесты выполните в терминале команду
go test -v
# Про веб сервис
написать апи у меня не получилось :(
Однако вы можете просто посмотреть на дизайн калькулятора, для этого перейдите в директорию calc_go_1/web.
В вс коде во вкладке Extensions(Ctrl+Shift+X) скачайте Live Server от Ritwick Dey для удобства запуска.
В правом нижнем углу появится кнопка Go Live, нажмите на нее
