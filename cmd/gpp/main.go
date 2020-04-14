package main

import (
	"fmt"
	"log"
	"os"

	"go-playstore-publisher/playpublisher"

	"github.com/urfave/cli/v2"
)

var apkFilePath string
var serviceAccountFilePath string
var packageNameID string

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
		Flags:    getFlags(),
		Name:     "go-play-publisher",
		Usage:    "Go - PlayStore Publisher",
		Action: func(c *cli.Context) error {
			fmt.Println("c :::", c)
			return nil
		},
	}
}

func getCommands() []*cli.Command {
	return []*cli.Command{
		{
			Name:   "list",
			Usage:  "List upload APK into the application in the Play Store",
			Action: actionListApk,
		},
		{
			Flags:  getUploadFlags(),
			Name:   "upload",
			Usage:  "Upload APK binary to the PlayStore",
			Action: actionUploadApk,
		},
	}
}

func getUploadFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Destination: &apkFilePath,
			Name:        "apkFile",
			Required:    true,
			Usage:       "The path to the APK to upload",
		},
		&cli.StringFlag{
			Destination: &packageNameID,
			Name:        "packageNameID",
			Required:    true,
			Usage:       "The package name ID",
		},
	}
}

func getFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Destination: &serviceAccountFilePath,
			Name:        "serviceAccountFile",
			Required:    true,
			Usage:       "The Play publisher service account file",
		},
	}
}

func actionListApk(c *cli.Context) error {
	client, err := playpublisher.NewClient(serviceAccountFilePath)
	if err != nil {
		return err
	}

	return client.ListService.List(c.Args().First())
}

func actionUploadApk(c *cli.Context) error {
	client, err := playpublisher.NewClient(serviceAccountFilePath)
	if err != nil {
		return err
	}

	file, err := os.Open(apkFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	return client.UploadService.Upload(packageNameID, file, "alpha")
}
