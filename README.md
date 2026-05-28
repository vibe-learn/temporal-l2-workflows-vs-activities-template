        # temporal — Workflow vs Activity

        Homework-шаблон для урока **l2_workflows_vs_activities** (Workflow vs Activity) на платформе Vibe Learn.

        ## Что делать

        На go.temporal.io/sdk реализуй OrderWorkflow и три активности (ChargeCard, ReserveStock,
SendEmail). Задай ActivityOptions через workflow.WithActivityOptions: StartToCloseTimeout
и RetryPolicy с MaximumAttempts. Активности — заглушки с логом и искусственной флакой
(падают с вероятностью p, потом успех). Сделай долгую активность ReserveStock с
HeartbeatTimeout и RecordHeartbeat. Тесты на TestWorkflowEnvironment проверят: workflow
доходит до конца несмотря на флаки (ретраи отрабатывают); при превышении StartToClose
попытка таймаутится и ретраится; heartbeat фиксируется.

## Контекст (из transfer-задачи урока)

Ты переносишь на Temporal обработку загруженного видео: (1) скачать исходник из S3,
(2) перекодировать в 3 разрешения (это 15-25 минут CPU-работы), (3) залить результаты
обратно в S3, (4) записать метаданные в PostgreSQL, (5) отправить вебхук клиенту.
Перекодирование иногда зависает наглухо, а воркеры периодически переезжают при деплое.

**Вопрос:** спроектируй это на Temporal. Опиши:
(a) что станет workflow, а что — активностями, и почему именно так проходит граница;
(b) какие таймауты и почему ты поставишь активности перекодирования (4 таймаута);
(c) как heartbeat поможет с «зависает наглухо» и переездом воркеров при деплое.

## Recap из урока

- **Workflow** — детерминированный оркестратор: реплеится, без прямого I/O, без time.Now/rand/goroutine. **Activity** — единица побочных эффектов: ей можно всё, выполняется один раз за попытку.
- Граница workflow/activity — не стиль, а **следствие replay**: весь контакт с внешним миром идёт через активности, иначе replay разъедется.
- Активность ретраится автоматически с гарантией **at-least-once** → делай её идемпотентной.
- Вызов: `workflow.ExecuteActivity(ctx, Fn, args)` возвращает Future; `.Get(ctx, &res)` durably ждёт результата, не занимая воркер.
- **Четыре таймаута**: ScheduleToStart (очередь), StartToClose (одна попытка — задавай всегда), ScheduleToClose (всё с ретраями), Heartbeat (для долгих активностей с RecordHeartbeat).

        ## Как работать

        1. Платформа Vibe Learn создаёт копию этого репо в твоём GitHub-аккаунте по клику «Начать домашку» на странице урока (через GitHub `/generate`, codecrafters-pattern).
        2. Склонируй копию локально, реализуй TODO в `main.go` (workflow + активности), прогони тесты, запушь.
        3. CI (`.github/workflows/ci.yml`) запускает `go vet` + `go test ./...` на каждый push. Платформа слушает результат через webhook от GitHub Actions и обновляет статус домашки на странице урока.

        ## Локальное окружение

        - Go 1.22+
        - SDK: `go.temporal.io/sdk`
        - Docker + docker-compose — `docker compose up` поднимает Temporal dev server на `:7233` + Web UI на `:8233`. Адрес переопределяется через env `TEMPORAL_ADDRESS` (дефолт `localhost:7233`).
        - Юнит-тесты на `testsuite.TestWorkflowEnvironment` (активности замоканы) бегут в CI БЕЗ сервера; интеграционный тест включается через `TEMPORAL_INTEGRATION=1`.

        ## Запуск

        ```bash
        # Поднять локальный Temporal dev server + UI
        docker compose up -d
        # Web UI: http://localhost:8233

        # Прогнать тесты (юнит на TestWorkflowEnvironment — без сервера;
        # интеграционный включается через TEMPORAL_INTEGRATION=1)
        go test ./...
        TEMPORAL_INTEGRATION=1 go test ./...

        # Запустить воркер (регистрирует workflow + активности, слушает task queue)
        go run .
        ```

        ## Заметка автора

        Это baseline-шаблон, сгенерированный платформой. Бизнес-сущность задачи (что конкретно реализовать в `main.go`, какие тесты сделать строгими) расширяется по ходу итераций — параллельно с углублением теории урока.
