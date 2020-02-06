package playpublisher

import (
	"fmt"
	"io"

	"google.golang.org/api/androidpublisher/v3"

	logr "github.com/sirupsen/logrus"
)

// UploadApkService definition
type UploadApkService struct {
	client        *Client
	packageNameID string
	editID        string
}

type UploadConfig struct {
	PackageName string
	FilePath    string
	Track       string
}

// Upload allows to list all the APK for the target bundle identifier into the PlayStore account
func (s *UploadApkService) Upload(id string,
	reader io.Reader,
	trackName string) error {

	s.packageNameID = id

	logr.WithFields(logr.Fields{
		"Package name ID": s.packageNameID,
		"Track name":      trackName,
	}).Info("Upload")

	// Create edit
	err := s.createEdit()
	if err != nil {
		return err
	}

	// Upload APK
	apk, err := s.uploadBinary(s.packageNameID, reader)
	if err != nil {
		return err
	}

	// Create release to track
	release := s.createReleaseToTrack("", apk.VersionCode)

	// Resolve track
	track, err := s.resolveTrack(s.packageNameID, trackName)

	// Release APK to track
	err = s.releaseApkToTrack(apk, track, release)
	if err != nil {
		return err
	}

	return s.validateAndCommitEdit()
}

func (s *UploadApkService) createEdit() error {
	logr.WithField("Package name ID", s.packageNameID).Info("Creating edit for package")

	// Create edit and handle error if any
	edit, err := s.client.Helper.createEdit(s.packageNameID)
	if err != nil {
		return fmt.Errorf("Failed to create edit (Error : %v)", err)
	}

	// Apply the editID for furthrer usage
	s.editID = edit.Id

	logr.WithField("Edit ID", s.editID).Info("Edit created")

	return nil
}

func (s *UploadApkService) uploadBinary(packageNameID string, reader io.Reader) (*androidpublisher.Apk, error) {
	logr.Info("Uploading binary")
	return s.client.Helper.initiateUpload(reader,
		packageNameID,
		s.editID,
		"application/vnd.android.package-archive")
}

func (s *UploadApkService) createReleaseToTrack(changeLog string, versionCode int64) *androidpublisher.TrackRelease {
	logr.WithField("Version code", versionCode).Info("Creating release")
	// TODO: Changelog support
	return s.client.Helper.createReleaseToTrack("", versionCode, 0)
}

func (s *UploadApkService) resolveTrack(packageNameID string,
	trackName string) (*androidpublisher.Track, error) {

	logr.Info("Resolving track", trackName)

	track, err := s.client.Helper.resolveTrackName(packageNameID, s.editID, trackName)
	if err != nil {
		err = fmt.Errorf("Failed to resolve track '%v' (Error: %v)", trackName, err)
	}

	return track, err
}

func (s *UploadApkService) releaseApkToTrack(apk *androidpublisher.Apk,
	track *androidpublisher.Track,
	release *androidpublisher.TrackRelease) error {
	logr.Info("Release to track", track)
	return s.client.Helper.releaseApkToTrack(s.packageNameID, apk, s.editID, track, release)
}

func (s *UploadApkService) validateAndCommitEdit() error {
	logr.Info("Validating commit edit", s.editID)
	// Validate the edit
	err := s.client.Helper.validateEdit(s.packageNameID, s.editID)
	if err != nil {
		return err
	}

	// Commit the edit
	return s.client.Helper.commitEdit(s.packageNameID, s.editID)
}
