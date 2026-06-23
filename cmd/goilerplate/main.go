package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v3"
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
			Name:  "count",
			Usage: "count to n",
			Flags: []cli.Flag{
				&cli.IntFlag{
					Name:    "n",
					Value:   3,
					Usage:   "number to count to",
					Aliases: []string{"n"},
				},
			},
			Action: count,
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

func count(_ context.Context, cmd *cli.Command) error {
	n := cmd.Int("n")
	for i := 1; i <= n; i++ {
		fmt.Println(i)
	}
	return nil
}

func main() {
	err := rootCmd.Run(context.Background(), os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
