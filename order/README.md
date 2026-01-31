# Order Service

Order — сервис управления заказами в проекте DeliveryFlow. Он принимает запросы на создание заказа, хранит состояние, публикует доменные события и отвечает на RPC-запросы.

## Цели сервиса
- создать заказ и зафиксировать статус
- хранить историю статусов
- публиковать события для следующих сервисов
- отвечать на запросы статуса по RPC

## Зоны взаимодействия
- **NATS Core**: RPC и публикация событий
- **JetStream**: запись событий в историю (через отдельный `audit` сервис)
- **PostgreSQL** (опционально): хранение заказов

## Основные функции (MVP)

### RPC
- `rpc.create_order` — создать заказ
- `rpc.get_order` — получить детали заказа
- `rpc.get_order_status` — получить статус заказа

### События
- `events.order_created` — заказ создан
- `events.order_cancelled` — заказ отменен (опционально)

## Статусы заказа (минимум)
- `new`
- `reserved`
- `paid`
- `delivery_assigned`
- `cancelled`

## Контракт данных (минимум)

### Order
- `id`
- `user_id`
- `items[]` (sku, qty)
- `total_amount`
- `status`
- `created_at`

## План реализации

### Этап 1: каркас
- минимальный `main.go`
- подключение к NATS
- обработка `rpc.create_order` с мок-репозиторием

### Этап 2: репозиторий
- in-memory реализация
- интерфейс `OrderRepository`

### Этап 3: события
- публикация `events.order_created`
- отдельный файл subjects/констант

### Этап 4: чтение статуса
- `rpc.get_order_status`
- `rpc.get_order`

### Этап 5: БД (опционально)
- Postgres репозиторий

## Структура папок (рекомендация)

```
order/
  cmd/
    order/
      main.go
  internal/
    app/
      app.go
    handler/
      nats_rpc.go
    service/
      order_service.go
    domain/
      order.go
      status.go
    repository/
      order_repo.go
      memory/
        order_repo.go
    messaging/
      publisher.go
      subjects.go
    config/
      config.go
```

## Заметки
- HTTP не обязателен: сервис может жить как RPC + события.
- Для тестов можно использовать in-memory репозиторий.
