package main

import (
	"log"
	"os"

	"go-playstore-publisher/playpublisher"

	"github.com/urfave/cli/v2"
)

func main() {
	app := initCli()
	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}

func initCli() *cli.App {
	return &cli.App{
		Commands: getCommands(),
		Name:     "go-play-publisher",
		Usage:    "Go - PlayStore Publisher",
		Flags:    getFlags(),
		Action: func(c *cli.Context) error {
			return nil
		},
	}
}

func getCommands() []*cli.Command {
	return []*cli.Command{
		{
			Name:  "list",
			Usage: "List upload APK into the application in the Play Store",
			Action: func(c *cli.Context) error {
				client, err := getClient(c)
				if err != nil {
					return err
				}
				return client.ListService.List(c.Args().First())
			},
		},
		{
			Name:  "upload",
			Usage: "Upload APK binary to the PlayStore",
			Action: func(c *cli.Context) error {
				client, err := getClient(c)
				if err != nil {
				}

				args := c.Args()

				return client.UploadService.Upload(args.First(),
					args.Get(1),
					args.Get(2))
			},
		},
	}
}

func getFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     "serviceAccountFile",
			Required: true,
			Usage:    "The Play publisher service account",
		},
	}
}

func getClient(c *cli.Context) (*playpublisher.Client, error) {
	return playpublisher.NewClient(c.String("serviceAccountFile"))
}

/*
app.Commands = []*cli.Command{
		&cli.Command{
			Name: "upload",
			Flags: []cli.Flag{
				&cli.PathFlag{Name: "file",
					Aliases:  []string{"f"},
					Required: true},
			},
			Action: executeUpload,
		},
	}


}
*/
