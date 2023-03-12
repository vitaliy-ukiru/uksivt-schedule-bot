```
├── adapters
│   ├── groups - "Сервис" для получения групп по критериям
|   ├── dao/{domain} - Реализации стораджей для сервисов
│   └── schedule - Надстройка над pkg/schedule-api
├── config    - Cтруктура конфига
├── delivery  - Слой транспорта
│    └── telegram - Слой для Telegram BOT API
│       ├── controller - Слой контроллеров для бота
│       │   ├── cron   - Контроллеры для управления крон-задачи (создание редактирование
│       │   │   ├── controller.go      - Стуктура контроллеров создания и изменения
│       │   │   ├── handlers_{module}.go - Хэндлеры контроллера {module}CronHandler
│       │   │   ├── handlers_edit.go   - Хэндлеры контроллера для изменения крон
│       │   │   ├── keyboard.go        - Общие объекты для клавиатур (инлайн)
│       │   │   ├── keyboard_{module}.go - Объекты для клавиатур контроллера {module}CronHandler
│       │   │   └── list.go            - Хэндлер для вывода списка крон в чате
│       │   ├── group - Контроллеры для выбора группы
│       │   │   ├── commands.go     - Хэндлеры на команды
│       │   │   ├── controller.go   - Структура
│       │   │   ├── keyboard.go     - Объекты клавиатуры
│       │   │   └── select_group.go - Хэндлеры для выбора группы
│       │   └── schedule
│       │       ├── controller.go     - Структура
│       │       ├── cron_scheduler.go - Обработчик крон задач
│       │       ├── explorer.go       - Просмотр раписания
│       │       └── keyboard.go       - Инлайн кнопки
│       ├── controller.go                - Стуктура, объединяет все контроллеры в один
│       └── controllers_pkg_shortcuts.go - Декораторы для конструкторов для более удобного создания
└── domain - Слой доменной области
    ├── chat
    ├── group     
    └── scheduler 

```
