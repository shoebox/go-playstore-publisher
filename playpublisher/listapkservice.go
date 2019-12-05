package playpublisher

import (
	"fmt"

	"google.golang.org/api/androidpublisher/v3"
)

// ListApkService definition
type ListApkService struct {
	editService *androidpublisher.EditsService
	client      *Client
}

// List allows to list all the APK for the target bundle identifier into the PlayStore account
func (s *ListApkService) List(packageName string) error {
	fmt.Println("Resolving APK list for package : ", packageName)

	edit, err := s.editService.Insert(packageName, nil).Do()
	if err != nil {
		return err
	}

	response, err := s.editService.Apks.List(packageName, edit.Id).Do()
	if err != nil {
		return err
	}

	for _, v := range response.Apks {
		fmt.Printf("\tVersion: %v - Binary sha1: %v\n", v.VersionCode, v.Binary.Sha1)
	}

	return nil
}
