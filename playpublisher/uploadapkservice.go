package playpublisher

import (
	"fmt"
	"io"
)

// UploadApkService definition
type UploadApkService struct {
	client *Client
}

type UploadConfig struct {
	PackageName string
	FilePath    string
	Track       string
}

// Upload allows to list all the APK for the target bundle identifier into the PlayStore account
func (s *UploadApkService) Upload(packageNameID string,
	reader io.Reader,
	trackName string) error {

	fmt.Println("Uploading", packageNameID, trackName)

	// Create edit
	edit, err := s.client.Helper.createEdit(packageNameID)
	if err != nil {
		return err
	}

	// Upload APK
	apk, err := s.client.Helper.initiateUpload(reader,
		packageNameID,
		edit.Id,
		"application/vnd.android.package-archive")

	if err != nil {
		return err
	}

	//
	release := s.client.Helper.createReleaseToTrack("", apk.VersionCode, 0)

	// Validate track name
	track, err := s.client.Helper.resolveTrackName(packageNameID, edit.Id, trackName)
	if err != nil {
		return fmt.Errorf("Failed to resolve track '%v' (Error: %v)", trackName, err)
	}

	// Release APK to track
	fmt.Println("Assigning release to track")
	err = s.client.Helper.releaseApkToTrack(packageNameID, apk, edit.Id, track, release)
	if err != nil {
		return err
	}

	// Validate the edit
	err = s.client.Helper.validateEdit(packageNameID, edit.Id)
	fmt.Println(err)
	if err != nil {
		return err
	}

	// Commit the edit
	return s.client.Helper.commitEdit(packageNameID, edit.Id)
}
