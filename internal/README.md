```
├── adapters
│   ├── groups - "Service" for get college groups
|   ├── dao/{domain} - Implementations storages for services
│   └── schedule - Abstract adapter pkg/schedule-api
├── config    - Structure of config
├── delivery  - Transport layer.
│    └── telegram - Layer of Telegram BOT API
│        ├── controller - Layer of controllers for bot
│        │   ├── cron   - Controller for manage crons-job (create, edit)
│        │   │   ├── controller.go        - Structure of controller CreateCronHandler & EditCronHandler
│        │   │   ├── handlers_{module}.go - Handlers for {module}CronHandler
│        │   │   ├── keyboard.go          - Common objects for inline markups (инлайн)
│        │   │   ├── keyboard_{module}.go - Objects for inline markups for {module}CronHandler
│        │   │   └── list.go              - Handler for output list of crons in chats
│        │   ├── group - Controllers for select group
│        │   │   ├── commands.go     - Handler on commands
│        │   │   ├── controller.go   - Structure of controller
│        │   │   ├── keyboard.go     - Objects of inline markup
│        │   │   └── select_group.go - Handlers for select_groups
│        │   └── schedule
│        │       ├── controller.go     - Strucutre
│        │       ├── cron_scheduler.go - Handler of crons job
│        │       ├── explorer.go       - Viewer of lessons
│        │       └── keyboard.go       - Inline markup's objects
│        ├── controller.go                - Structure, uniting all controller in one object
│        └── controllers_pkg_shortcuts.go - Decorators of constrcutor (for more comfortable init controllers from packages)
└── domain - Domain layer
    ├── chat
    ├── group     
    └── scheduler 
```
