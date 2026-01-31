# DeliveryFlow — учебный проект на NATS

DeliveryFlow — мини-система оформления и доставки заказов. Это цельный проект с несколькими сервисами и инфраструктурой, где естественно используются разные паттерны NATS.

## Архитектура

### Сервисы
- `api-gateway` — HTTP API для клиентов (создать заказ, получить статус).
- `identity` — аутентификация (JWT), управление пользователями.
- `order` — создание заказа и публикация доменных событий.
- `inventory` — резервирование товара.
- `payment` — списание/отмена платежа.
- `delivery` — назначение курьера и статусы доставки.
- `review` — отзывы к товарам и рейтинги.
- `notification` — отправка уведомлений (email/push).
- `audit` — запись истории событий.

### Инфраструктура
- **NATS Core** (Cluster 3 nodes) — шина событий и RPC.
- **JetStream** — сохранение истории доменных событий.
- **PostgreSQL** — состояние заказов, платежей и пользователей.
- **MongoDB** — хранение отзывов (документная модель).
- **Redis** — хранение сессий и временных токенов.

### Зоны взаимодействия
1. **Event Bus (NATS Core)**  
   Обмен доменными событиями между сервисами.

2. **Queue Groups (Tasks)**  
   Фоновые задачи уведомлений и обработки очередей.

3. **Request/Reply (RPC)**  
   Синхронные запросы статусов, проверки наличия и валидация токенов.

4. **JetStream (History)**  
   История заказов и возможность реплея событий.

## Диаграмма потоков

```mermaid
flowchart LR
  subgraph ClientZone[Client Zone]
    C[client]
  end

  subgraph NATSZone[NATS Zone]
    direction TB
    subgraph NATSCluster[Core Cluster]
      direction LR
      NATS[NATS Core 1] --- NATS2[NATS Core 2] --- NATS3[NATS Core 3]
    end
    JS[JetStream]
  end

  subgraph Services[Services]
    API[api-gateway]
    ID[identity]
    O[order]
    I[inventory]
    P[payment]
    D[delivery]
    R[review]
    N[notification]
    A[audit]
  end

  C -->|HTTP auth/order/review| API
  
  %% Auth Flow
  API -->|request rpc.auth.login| NATS
  NATS -->|request| ID
  ID -->|reply token| NATS
  ID -->|publish events.user_registered| NATS
  
  %% Order Flow
  API -->|request rpc.create_order| NATS
  NATS -->|request| O
  O -->|reply order_id| NATS

  O -->|publish events.order_created| NATS
  NATS -->|subscribe events.order_created| I
  I -->|publish events.inventory_reserved| NATS
  NATS -->|subscribe events.inventory_reserved| P
  P -->|publish events.payment_captured| NATS
  NATS -->|subscribe events.payment_captured| D
  D -->|publish events.delivery_assigned| NATS
  
  %% Review Flow
  API -->|request rpc.create_review| NATS
  NATS -->|request| R
  R -->|reply ok| NATS
  R -->|publish events.review_created| NATS

  %% Notifications & Audit
  NATS -->|queue group jobs.notify_user| N
  NATS -->|stream events.*| JS
  JS -->|deliver historical events| A
```

## Потоки сообщений

### Pub/Sub (доменные события)
- `identity` → `events.user_registered` → `notification` (Email "Welcome")
- `order` → `events.order_created` → `inventory`
- `inventory` → `events.inventory_reserved` → `payment`
- `payment` → `events.payment_captured` → `delivery`
- `review` → `events.review_created` → `notification` (Email "Thank you for review")

### Queue Group (фоновая обработка)
- `events.*` → `jobs.notify_user` → `notification` (несколько воркеров)

### Request/Reply
- `api-gateway` → `rpc.auth.login` → `identity` → токен
- `api-gateway` → `rpc.auth.verify` → `identity` → ok/error
- `api-gateway` → `rpc.create_order` → `order` → ответ
- `api-gateway` → `rpc.get_order_status` → `order` → ответ
- `api-gateway` → `rpc.create_review` → `review` → ответ

### JetStream
- `events.*` → JetStream → `audit` → сохранение истории

## NATS темы (subjects)

### Events
- `events.user_registered`
- `events.order_created`
- `events.inventory_reserved`
- `events.payment_captured`
- `events.delivery_assigned`
- `events.review_created`

### Jobs
- `jobs.notify_user`

### RPC
- `rpc.auth.login`
- `rpc.auth.verify`
- `rpc.create_order`
- `rpc.get_order_status`
- `rpc.create_review`
