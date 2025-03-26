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
1) для запуска бд выбрала docker compose
`docker compose -f deployments/docker-compose.yaml up`
2) затем нужно выполнить миграции
`make migration-up`
3) а теперь можно запустить сам проект
`go run ./cmd/calendar --config=./configs/config.yaml`
Если всё успешно, то будет такой вывод:
```
[INFO] 2025-03-24 03:01:06 calendar is running...
[{8 test title 2025-03-24 03:01:06.746575 +0000 UTC 2025-03-24 06:01:06.746575 +0000 UTC testing}]
```