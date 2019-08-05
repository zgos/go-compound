package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"strconv"

	"github.com/postables/go-compound/client"
	"github.com/urfave/cli"
)

var (
	url string
)

func main() {
	app := newApp()
	if err := app.Run(os.Args); err != nil {
		log.Println("error running application: ", err.Error())
		os.Exit(1)
	}
}

func newApp() *cli.App {
	app := cli.NewApp()
	app.Name = "go-compound cli"
	app.Authors = []cli.Author{
		{
			Name:  "Alex Trottier (postables)",
			Email: "postables@rtradetechnologies.com",
		},
	}
	app.Flags = loadFlags()
	app.Commands = loadCommands()
	return app
}

func loadFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:        "api.url, au",
			Usage:       "the compound api url",
			Value:       "https://api.compound.finance/api/v2",
			Destination: &url,
		},
	}
}
func loadCommands() cli.Commands {
	return loadAccountCommands()
}

func loadAccountCommands() cli.Commands {
	return cli.Commands{
		cli.Command{
			Name:    "account",
			Usage:   "account management commands",
			Aliases: []string{"acct"},
			Subcommands: cli.Commands{
				cli.Command{
					Name:  "health",
					Usage: "get account health",
					Action: func(c *cli.Context) error {
						if c.String("eth.address") == "" {
							return errors.New("eth.address flag is empty")
						}
						cl := client.NewClient(url)
						resp, err := cl.GetAccount(c.String("eth.address"))
						if err != nil {
							return err
						}
						if len(resp.Accounts) == 0 {
							return errors.New("an unexpected error occurred")
						}
						health, err := strconv.ParseFloat(resp.Accounts[0].Health.Value, 64)
						if err != nil {
							return err
						}
						fmt.Println("account health: ", health)
						if health <= 1.0 {
							fmt.Printf("health is %v and at risk of liquidation\n", health)
						} else if health <= 1.2 {
							fmt.Printf("health is %v and nearing liquidation\n", health)
						}
						return nil
					},
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "eth.address",
							Usage: "the address to lookup",
						},
					},
				},
			},
		},
	}
}
