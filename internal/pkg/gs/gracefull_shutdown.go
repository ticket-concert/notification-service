package joshu

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
)

type GracefulShutdown struct {
	Timeout        time.Duration
	GracefulPeriod time.Duration

	terminated bool
	closers    []Closer
	wg         *sync.WaitGroup
}

func (gs *GracefulShutdown) Enable(app *fiber.App) {
	gs.wg = new(sync.WaitGroup)
	gs.wg.Add(1)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-sig
		gs.terminated = true
		time.Sleep(gs.GracefulPeriod)

		if err := app.Shutdown(); err != nil {
			log.Printf("Server Shutdown Failed:%+v\n", err)
		}

		gs.wg.Done()
	}()
}

// func (gs *GracefulShutdown) Enable(server *http.Server) {
// 	gs.wg = new(sync.WaitGroup)
// 	gs.wg.Add(1)

// 	sig := make(chan os.Signal, 1)
// 	signal.Notify(sig, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

// 	go func() {
// 		<-sig
// 		gs.terminated = true
// 		time.Sleep(gs.GracefulPeriod)

// 		ctx, cancel := context.WithTimeout(context.Background(), gs.Timeout)
// 		defer cancel()

// 		if err := server.Shutdown(ctx); err != nil {
// 			log.Printf("Server Shutdown Failed:%+v\n", err)
// 		}

// 		gs.wg.Done()
// 	}()
// }

func (gs *GracefulShutdown) Register(c ...Closer) {
	gs.closers = append(gs.closers, c...)
}

func (gs *GracefulShutdown) Cleanup() {
	gs.wg.Wait()

	ctx, cancel := context.WithTimeout(context.Background(), gs.Timeout)
	defer cancel()

	for _, c := range gs.closers {
		gs.triggerClose(ctx, c)
	}
}

func (gs *GracefulShutdown) triggerClose(ctx context.Context, c Closer) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic and recover from: %+v\n", r)
		}
	}()

	err := c.Close(ctx)
	if err != nil {
		log.Printf("Failed to cleanup resource:%+v\n", err)
	}

}

func (gs *GracefulShutdown) LivenessCheck(c *fiber.Ctx) error {
	return c.SendString("ok")
}

func (gs *GracefulShutdown) ReadinessCheck(c *fiber.Ctx) error {
	if gs.terminated {
		return c.Status(503).SendString("graceful shutdown")
	} else {
		return c.SendString("ok")
	}
}

type Closer interface {
	Close(context.Context) error
}

type FnWithContextAndError func(context.Context) error

func (fn FnWithContextAndError) Close(ctx context.Context) error {
	return fn(ctx)
}

type FnWithContext func(context.Context)

func (fn FnWithContext) Close(ctx context.Context) error {
	fn(ctx)
	return nil
}

type FnWithError func() error

func (fn FnWithError) Close(ctx context.Context) error {
	return fn()
}

type Fn func()

func (fn Fn) Close(ctx context.Context) error {
	fn()
	return nil
}
