package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	var query string
	var header string

	app := &cli.App{
		Name:  "httpcli",
		Usage: "REST APIs",
		Commands: []*cli.Command{
			{
				Name:  "get",
				Usage: "todo-later",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "query",
						Usage:       "todo-later",
						Destination: &query,
					},
					&cli.StringFlag{
						Name:        "header",
						Usage:       "todo-later",
						Destination: &header,
					},
				},
				Action: func(cCtx *cli.Context) error {
					fmt.Println("get method:", cCtx.Args())
					return nil
				},
			},
			{
				Name:  "post",
				Usage: "todo-later",
				Action: func(cCtx *cli.Context) error {
					fmt.Println("post method:", cCtx.Args())
					return nil
				},
			},
			{
				Name:  "put",
				Usage: "todo-later",
				Action: func(cCtx *cli.Context) error {
					fmt.Println("put method:", cCtx.Args())
					return nil
				},
			},
			{
				Name:  "delete",
				Usage: "todo-later",
				Action: func(cCtx *cli.Context) error {
					fmt.Println("delete method:", cCtx.Args())
					return nil
				},
			},
		},
		Action: func(*cli.Context) error {
			fmt.Println("Hello friend!")
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
