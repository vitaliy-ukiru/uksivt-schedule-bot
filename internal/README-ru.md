```
├── adapters
│   ├── groups - "Сервис" для получения групп по критериям
|   ├── storage/{domain} - Реализации стораджей для сервисов
│   └── schedule - Надстройка над pkg/schedule-api
├── chat - Домейн чата. Хранит его модель, сервис
├── scheduler - Домейн для шедулера.
└── delivery        - Слой транспорта.
    └── telegram    - Слой для телеграм бота
        ├── controller.go      - Структура хэндлера
        ├── cron.go            - Методы для работы с кронами
        ├── select_group.go    - Меню выбора группы
        ├── uksivt_schedule.go - Отправка расписания
        └── keyboards/         - Вспомогательные функции для генерации инлайн клавиату
```