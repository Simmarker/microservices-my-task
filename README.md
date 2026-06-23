# Домашняя работа #1

![GitHub Classroom Workflow](../../workflows/GitHub%20Classroom%20Workflow/badge.svg?branch=master)

## Microservices

### Формулировка

В рамках домашней работы требуется реализовать систему _Агрегатора Билетов в Кино_, состоящую из
трех взаимодействующих друг с другом сервисов: Ticket Service, Film Service, Cinema Service.

Ниже описано публичное API, оно описывает методы, предоставляемые пользователю. Внутреннюю
коммуникацию между сервисами вы проектируете и реализуете сами.

* Ticket Service – 8080;
* Film Service – 8070;
* Cinema Service – 8060;

Полное описание публичного API приведено в спецификации [OpenAPI](openapi.yml).

##### Просмотр списка всех фильмов, которые идут в кино

```http request
GET {{filmsUrl}}/api/v1/films?page=1&size=10
Accept: application/json
```

##### Просмотр списка всех кинотеатров

```http request
GET {{cinemaUrl}}/api/v1/cinema?page=1&size=10
Accept: application/json
```

##### Просмотр афиши выбранного кинотеатра

```http request
GET {{cinemaUrl}}/api/v1/cinema/{{cinemaUid}}/films
Accept: application/json
```

##### Купить билет в кино

Запрос приходит на Ticket Service, в Cinema Service выполняется
проверка `film_session.total_seats` > `film_session.booked_seats`, если все успешно, то `film_session.booked_seats`
увеличивается на 1 и в Ticket Service создается ticket со статусом `BOOKED`.

```http request
POST {{ticketUrl}}/api/v1/tickets/cinema/{{cinemaUid}}/films/{{filmUid}}
Content-Type: application/json
Accept: application/json
X-User-Name: {{username}}

{
   "date": "2024-01-01T08:00:00",
   "row": 10,
   "seat": 15
}
```

В ответ приходит:

`201 Created` и ссылка на билет в заголовке `Location: /api/v1/tickets/{ticketUid}`

##### Просмотр информации о билете

```http request
GET {{ticketUrl}}/api/v1/tickets/{{ticketUid}}
Accept: application/json
X-User-Name: {{username}}
```

В ответ приходит:

```json
{
  "ticketUid": "<ticketUid>",
  "status": "BOOKED",
  "date": "2024-01-01T08:00:00",
  "seat": 15,
  "row": 10
}
```

##### Вернуть билет, если до сеанса осталось больше 1 часа

В Ticket Service выполняется поиск билета по `ticketUid`, проверяется что до начала сеанса больше 1 часа, если условие
выполняется, то билет помечается `CANCELED`, но не удаляется. Если до сеанса меньше 1 часа, то
возвращается `409 Conflict`.

```http request
DELETE {{ticketUrl}}/api/v1/tickets/{{ticketUid}}
Accept: application/json
X-User-Name: {{username}}
```

В ответ приходит:

`204 No Content` с пустым телом ответа. Билет помечается как `CANCELED`, но не удаляется.

#### Структура Базы Данных

Ниже приведено _примерное_ описание структуры баз данных каждого сервиса, вы можете менять структуру таблиц, если
считаете что она будет лучше ложиться на реализацию.

##### Film Service

```sql
CREATE TABLE film
(
    id       SERIAL PRIMARY KEY,
    film_uid uuid          NOT NULL,
    name     VARCHAR(255)  NOT NULL,
    rating   NUMERIC(8, 2) NOT NULL DEFAULT 10
        CHECK ( rating BETWEEN 0 AND 10 ),
    director VARCHAR(255),
    producer VARCHAR(255),
    genre    VARCHAR(255)  NOT NULL
);

CREATE UNIQUE INDEX udx_film_uid ON film (film_uid);
```

##### Cinema Service

```sql
CREATE TABLE cinema
(
    id         SERIAL PRIMARY KEY,
    cinema_uid uuid NOT NULL,
    name       VARCHAR(255),
    address    VARCHAR(255)
);

CREATE UNIQUE INDEX udx_cinema_uid ON cinema (cinema_uid);

CREATE TABLE film_session
(
    id           SERIAL PRIMARY KEY,
    session_uid  uuid      NOT NULL,
    film_uid     uuid      NOT NULL,
    total_seats  INT       NOT NULL,
    booked_seats INT       NOT NULL DEFAULT 0
        CHECK ( booked_seats < total_seats ),
    date         TIMESTAMP NOT NULL,
    cinema_id    INT
        CONSTRAINT fk_film_session_cinema_id REFERENCES cinema (id)
);

CREATE UNIQUE INDEX udx_film_session_session_uid ON film_session (session_uid);
```

##### Ticket Service

```sql
CREATE TABLE tickets
(
    id          SERIAL PRIMARY KEY,
    ticket_uid  uuid        NOT NULL,
    film_uid    uuid        NOT NULL,
    session_uid uuid        NOT NULL,
    user_name   VARCHAR(80) NOT NULL,
    date        TIMESTAMP   NOT NULL,
    status      VARCHAR(20) NOT NULL
        CHECK ( status IN ('BOOKED', 'CANCELED') )
);

CREATE UNIQUE INDEX udx_tickets_ticket_uid ON tickets (ticket_uid);
```

#### Тестовые данные

Для успешного прохождения тестов в базе в таблицах должны быть данные:

##### Film Service

```yaml
films:
  - id: 1
    film_uid: "049161bb-badd-4fa8-9d90-87c9a82b0668"
    name: "Terminator 2 Judgment day"
    rating: 8.6
    director: "James Cameron"
    producer: "James Cameron"
    genre: "Sci-Fi"
```

##### Cinema Service

```yaml
cinema:
  - id: 1
    cinemaUid: "06cc4ba3-ee97-4d29-a814-c40588290d17",
    name: "Кинотеатр Москва",
    address: "Ереван, улица Хачатура Абовяна, 18"
```

```yaml
film_session:
  - id: 1
    cinema_id: 1
    film_uid: "049161bb-badd-4fa8-9d90-87c9a82b0668"
    date: "2024-01-01T08:00:00"
    total_seats: 5000
    booked_seats: 0
```

### Требования

1. Каждый сервис имеет свое собственное хранилище, если оно ему нужно. Для локальной разработки можно использовать
   Postgres 15 в [docker](docker-compose.yml), для этого нужно запустить `docker compose up postgres -d`, поднимется
   контейнер, будет создан пользователь `program`:`test` и 3 БД: `films`, `tickets`, `cinema`.
2. Для межсервисного взаимодействия использовать HTTP (придерживаться RESTful).
3. На каждом сервисе сделать специальный endpoint `GET /manage/health`, отдающий 200 ОК, он будет использоваться для
   проверки доступности сервиса (в [Github Actions](.github/workflows/classroom.yml) в скрипте проверки готовности всех
   сервисов [wait-script.sh](scripts/wait-script.sh)).
   ```shell
   ./scripts/wait-for.sh -t 120 "http://localhost:$port/manage/health" -- echo "Host localhost:$port is active"
   ```
4. Код хранить на Github, для сборки использовать Github Actions.
5. Каждый сервис должен быть завернут в docker.
6. В [build.yml](.github/workflows/classroom.yml) дописать шаги на сборку.
7. Интеграционные тесты можно проверить локально, для этого нужно импортировать в Postman
   коллекцию [collection.json](postman/collection.json)] и environment [local-env.json](postman/local-env.json).

### Пояснения

Для тестирования _все сервисы_ поднимаются через docker compose и для них запускаются тесты.

### Прием задания

1. При получении задания у вас создается _копия_ этого репозитория для вашего пользователя.
2. После того как все тесты успешно завершатся, в Github Classroom на Dashboard будет отмечено успешное выполнение
   тестов.
