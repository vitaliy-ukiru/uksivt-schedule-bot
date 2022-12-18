```
├── adapters
│   ├── groups - "Service" for get college groups
|   ├── storage/{domain} - Implementations storages for services
│   └── schedule - Abstract adapter pkg/schedule-api
├── chat - Chat's domain.
├── scheduler - Scheduler's domain.
└── delivery        - Transport layer.
    └── telegram    - Telegram BotAPI
        ├── controller.go      - Handler declaration, bindings.
        ├── cron.go            - Methods for cron handlers.
        ├── select_group.go    - Select groups handlers.
        ├── uksivt_schedule.go - Sending schedule.
        └── keyboards/         - Utilities for inline keyboard
```
