// Package main is the temporal lesson `l2_workflows_vs_activities` homework scaffold for Vibe Learn.
//
// Задача: OrderWorkflow с ActivityOptions (StartToCloseTimeout, RetryPolicy) и heartbeat у долгой ReserveStock.
// Реализуй workflow и активности ниже — сигнатуры и тестовая поверхность
// фиксированы; CI (.github/workflows/ci.yml) гоняет `go vet` и `go test ./...`.
// Подробности и критерии приёмки — в README.md.
//
// SDK: go.temporal.io/sdk (worker + workflow + activity).
// Воркер подключается к Temporal по TEMPORAL_ADDRESS (дефолт localhost:7233 —
// совпадает с docker-compose.yml) и слушает task queue из TaskQueue().
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

// ----- config -----

// envOr returns the env var for `key` if set, else `fallback`.
func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// TemporalAddress — адрес Temporal frontend. Дефолт совпадает с docker-compose.yml.
func TemporalAddress() string {
	return envOr("TEMPORAL_ADDRESS", "localhost:7233")
}

// TaskQueue — очередь задач, которую слушает воркер этого урока.
func TaskQueue() string {
	return envOr("TEMPORAL_TASK_QUEUE", "lesson-l2_workflows_vs_activities-tq")
}

// ----- Workflow: OrderWorkflow -----
//
// Оркеструет активности ниже. Тело — TODO: добавь ExecuteActivity-шаги,
// ActivityOptions (StartToCloseTimeout, RetryPolicy) и обработку ошибок
// согласно README.md. Должно оставаться ДЕТЕРМИНИРОВАННЫМ (никаких
// time.Now/rand/итераций по map — используй workflow.Now/SideEffect).
func OrderWorkflow(ctx workflow.Context) error {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 30 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	logger := workflow.GetLogger(ctx)
	logger.Info("OrderWorkflow started", "taskQueue", TaskQueue())

	// TODO #1: вызови активность ChargeCard через workflow.ExecuteActivity.
	// var chargecardRes string
	// if err := workflow.ExecuteActivity(ctx, ChargeCard).Get(ctx, &chargecardRes); err != nil {
	// 	return err
	// }
	// TODO #2: вызови активность ReserveStock через workflow.ExecuteActivity.
	// var reservestockRes string
	// if err := workflow.ExecuteActivity(ctx, ReserveStock).Get(ctx, &reservestockRes); err != nil {
	// 	return err
	// }
	// TODO #3: вызови активность SendEmail через workflow.ExecuteActivity.
	// var sendemailRes string
	// if err := workflow.ExecuteActivity(ctx, SendEmail).Get(ctx, &sendemailRes); err != nil {
	// 	return err
	// }

	return nil
}

// ----- Activity #1: ChargeCard -----
//
// списать деньги; флакает с вероятностью p, потом успех
func ChargeCard(ctx context.Context) (string, error) {
	// TODO: implement
	return "", fmt.Errorf("ChargeCard: not implemented")
}

// ----- Activity #2: ReserveStock -----
//
// долгая активность: RecordHeartbeat в цикле, HeartbeatTimeout в опциях
func ReserveStock(ctx context.Context) (string, error) {
	// TODO: implement
	return "", fmt.Errorf("ReserveStock: not implemented")
}

// ----- Activity #3: SendEmail -----
//
// отправить письмо-подтверждение
func SendEmail(ctx context.Context) (string, error) {
	// TODO: implement
	return "", fmt.Errorf("SendEmail: not implemented")
}

// ----- main entry: register worker + run with graceful shutdown -----

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.Printf("Vibe Learn — temporal lesson %s scaffold up", "l2_workflows_vs_activities")
	log.Printf("temporal address: %s  task queue: %s", TemporalAddress(), TaskQueue())
	log.Printf("Реализуй workflow и активности, затем `go test ./...`. README.md содержит задачу.")

	c, err := client.Dial(client.Options{HostPort: TemporalAddress()})
	if err != nil {
		log.Fatalf("unable to create Temporal client (is `docker compose up -d` running?): %v", err)
	}
	defer c.Close()

	w := worker.New(c, TaskQueue(), worker.Options{})
	w.RegisterWorkflow(OrderWorkflow)
	w.RegisterActivity(ChargeCard)
	w.RegisterActivity(ReserveStock)
	w.RegisterActivity(SendEmail)

	// Graceful shutdown so `go run .` is interactive — worker.InterruptCh()
	// stops the worker on Ctrl-C / SIGTERM.
	if err := w.Run(worker.InterruptCh()); err != nil {
		log.Fatalf("worker stopped with error: %v", err)
	}
}
