#### Результатом выполнения следующих домашних заданий является сервис «Календарь»:
- [Домашнее задание №12 «Заготовка сервиса Календарь»](./docs/12_README.md)
- [Домашнее задание №13 «Внешние API от Календаря»](./docs/13_README.md)
- [Домашнее задание №14 «Кроликизация Календаря»](./docs/14_README.md)
- [Домашнее задание №15 «Докеризация и интеграционное тестирование Календаря»](./docs/15_README.md)

#### Ветки при выполнении
- `hw12_calendar` (от `master`) -> Merge Request в `master`
- `hw13_calendar` (от `hw12_calendar`) -> Merge Request в `hw12_calendar` (если уже вмержена, то в `master`)
- `hw14_calendar` (от `hw13_calendar`) -> Merge Request в `hw13_calendar` (если уже вмержена, то в `master`)
- `hw15_calendar` (от `hw14_calendar`) -> Merge Request в `hw14_calendar` (если уже вмержена, то в `master`)
- `hw16_calendar` (от `hw15_calendar`) -> Merge Request в `hw15_calendar` (если уже вмержена, то в `master`)


**Домашнее задание не принимается, если не принято ДЗ, предшествующее ему.**


Как запускать (все команды выполняем в директории домашнего задания `hw12_13_14_15_16_calendar`): 
1) для запуска сервисов(postgresql, rabbitmq) выбрала docker compose
`docker compose -f deployments/docker-compose.yaml up`
2) затем нужно выполнить миграции
`make migration-up`
3) запустить календарь
`go run ./cmd/calendar --config=./configs/config.yaml`
Если всё успешно, то будет такой вывод:
```
[INFO] 2025-03-24 03:01:06 calendar is running...
```
4) запустить scheduler
`go run ./cmd/scheduler --config=./configs/config_scheduler.yaml`
Если всё успешно, то будет такой вывод:
```
[INFO] 2025-05-11 13:49:27 scheduler is running...
```
4) запустить sender
`go run ./cmd/sender --config=./configs/config_sender.yaml`
Если всё успешно, то будет такой вывод:
```
[INFO] 2025-05-11 13:50:48 sender is running...
```
Если в хранилище будут подходящие данные, то будет такой вывод:
```
[INFO] 2025-05-11 13:50:48  SEND {"ID":2,"title":"test title","startDate":{},"endDate":{},"description":"testing"}
[INFO] 2025-05-11 13:50:48  SEND {"ID":4,"title":"test title","startDate":{},"endDate":{},"description":"testing"}
[INFO] 2025-05-11 13:50:48  SEND {"ID":6,"title":"test title","startDate":{},"endDate":{},"description":"testing"}
[INFO] 2025-05-11 13:50:48  SEND {"ID":7,"title":"test title","startDate":{},"endDate":{},"description":"testing"}
```