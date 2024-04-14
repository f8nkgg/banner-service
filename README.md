# Banner Service Avito [Backend Application] 

## Содержание 
- [Запуск](#запуск)
- [Примеры запросов и ответов](#примеры-запросов-и-ответов)
- [Реализация](#реализация)
## Запуск
Запуск с помощью docker compose
```sh
$ make compose-up
```
Интеграционные тесты
```sh
$ make test
```
Обычный запуск сервера
```sh
$ make run
```
Запуск линтера
```sh
$ make lint
```
### Зависимости
- go 1.22.1
- docker & docker-compose
- [golangci-lint](https://github.com/golangci/golangci-lint) (для проверки кода)

Конфигурационный файл лежит по адресу config/config.yml

## Примеры запросов и ответов

### POST  http://localhost:8080/banner

Создание нового баннера

##### Запрос:
```
Content-Type: application/json
Token: admin_token

{
  "tag_ids": [4, 5, 6],
  "feature_id": 123,
  "content": {
    "title": "some_title",
    "text": "some_text",
    "url": "some_url"
  },
  "is_active": true
}

```
##### Ответ:
```
{
  "banner_id": 1
}
```

### GET http://localhost:8080/user_banner?tag_id=4&feature_id=123&use_last_revision=true

Получение баннера для пользователя

##### Запрос:
```
Content-Type: application/json
Token: user_token
```

##### Ответ:
```
{
  "text": "some_text",
  "title": "some_title",
  "url": "some_url"
}
```

### GET http://localhost:8080/banner?feature_id=123&offset=0

Получение всех баннеров c фильтрацией по фиче и/или тегу

##### Запрос:
```
Content-Type: application/json
Token: admin_token
```

##### Ответ:
```
[
  {
    "id": 1,
    "tag_ids": [
      4,
      5,
      6
    ],
    "feature_id": 123,
    "content": {
      "text": "some_text",
      "title": "some_title",
      "url": "some_url"
    },
    "is_active": true,
    "created_at": "2024-04-14T03:44:58.059736+04:00",
    "updated_at": "2024-04-14T03:44:58.059736+04:00"
  }
]
```
### PATCH http://localhost:8080/banner/1

Обновление содержимого баннера

##### Запрос:
```
Content-Type: application/json
Token: admin_token
```

### DELETE http://localhost:8080/banner/1

Удаление баннера по идентификатору

##### Запрос:
```
Content-Type: application/json
Token: admin_token
```

### GET http://localhost:8080/banner/history/1

Получение истории изменений баннера по идентификатору

##### Запрос:
```
Content-Type: application/json
Token: admin_token
```

##### Ответ:
```
[
  {
    "Index": 1,
    "Banner": {
      "id": 1,
      "tag_ids": [
        4,
        5,
        6
      ],
      "feature_id": 123,
      "content": {
        "text": "some_text",
        "title": "some_title",
        "url": "some_url"
      },
      "is_active": true,
      "created_at": "2024-04-14T04:56:16.643184+04:00",
      "updated_at": "2024-04-14T04:56:18.809631+04:00"
    }
  }
]
```
## Реализация

##### 1. Версионирование баннеров

Считая, что основное назначение сервиса - получение пользователями баннеров, а работа администратора в нем представляет 
малую часть, база данных была спроектирована с целью производительности при выдаче данных. Версии баннеров - полностью 
функционал администратора, и нет смысла увеличивать таблицу в 4 раза. Поэтому создаем таблицу истории (SCD Type 4) с 
триггером на изменение оригинальной таблицы. Для возврата старого значения можно отправить patch запрос с интересующей 
версией при получении истории.

##### 2. Кэш

Так как сказано адаптировать систему с допущением увеличения времени исполнения по редко запрашиваемым тегам и фичам, то 
реализовываем LFU cache.

<br>В следующей жизни сделать шардирование бд