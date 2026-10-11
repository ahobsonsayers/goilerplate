package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/urfave/cli/v3"

	"github.com/ahobsonsayers/goilerplate/worker"
)

var rootCmd = &cli.Command{
	Name:  "goilerplate",
	Usage: "goilerplate cli",
	Commands: []*cli.Command{
		{
			Name:  "greet",
			Usage: "print a greeting",
			Arguments: []cli.Argument{
				&cli.StringArg{
					Name: "name",
				},
			},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "lang",
					Value:   "english",
					Usage:   "language for the greeting (english|spanish)",
					Aliases: []string{"l"},
				},
				&cli.BoolFlag{
					Name:    "shout",
					Usage:   "uppercase the greeting",
					Aliases: []string{"s"},
				},
			},
			Action: greet,
		},
		{
			Name:  "worker",
			Usage: "run the example worker until interrupted",
			Flags: []cli.Flag{
				&cli.DurationFlag{
					Name:  "period",
					Value: time.Second,
					Usage: "polling period for the worker",
				},
			},
			Action: runWorker,
		},
	},
	Action: run,
}

func run(_ context.Context, cmd *cli.Command) error {
	return cli.ShowRootCommandHelp(cmd)
}

func greet(_ context.Context, cmd *cli.Command) error {
	name := cmd.StringArg("name")
	if name == "" {
		name = "world"
	}

	greeting := "Hello"
	if cmd.String("lang") == "spanish" {
		greeting = "Hola"
	}

	msg := fmt.Sprintf("%s %s", greeting, name)
	if cmd.Bool("shout") {
		msg += "!"
	}

	fmt.Println(msg)

	return nil
}

func runWorker(ctx context.Context, cmd *cli.Command) error {
	config := worker.Config{
		PollingPeriod: cmd.Duration("period"),
	}

	w, err := worker.NewWorker(config)
	if err != nil {
		return err
	}

	// Register signal to stop worker
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(signalChan)
	go func() {
		<-signalChan
		w.Stop()
	}()

	// Create cancellable context
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Print worker result and error channels
	go func() {
		for result := range w.Results() {
			fmt.Printf("Current Time: %v\n", result)
		}
	}()
	go func() {
		for err := range w.Errors() {
			log.Printf("Worker error: %v", err)
		}
	}()

	// Run worker (in a goroutine)
	workerErrChan := make(chan error, 1)
	go func() {
		workerErrChan <- w.Start(ctx)
	}()

	// Wait for worker to stop
	w.WaitUntilStopped()

	workerErr := <-workerErrChan
	if errors.Is(workerErr, context.Canceled) {
		// Stopped cleanly via signal
		return nil
	}

	return workerErr
}

func main() {
	err := rootCmd.Run(context.Background(), os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
