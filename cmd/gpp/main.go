package main

import (
	"fmt"
	"log"

	"go-playstore-publisher/playpublisher"
)

func main() {
	client, err := playpublisher.NewClient("file.json")
	if err != nil {
		log.Fatal(err)
	}

	client.ListService.List("co.massive.bein.beta")

	err = client.UploadService.Upload("co.massive.bein.beta", "test.apk", "alpha")
	fmt.Println(err)
}
