package playpublisher

import (
	"fmt"
	"io"
	"os"

	"google.golang.org/api/androidpublisher/v3"
	"google.golang.org/api/googleapi"
)

// UploadApkService definition
type UploadApkService struct {
	client *Client
}

// Upload allows to list all the APK for the target bundle identifier into the PlayStore account
func (s *UploadApkService) Upload(packageName string, path string, track string) error {
	fmt.Println("Uploading", path)

	// Create edit
	edit, err := s.createEdit(packageName)
	if err != nil {
		return err
	}

	// Upload APK
	apk, err := s.uploadApk(packageName, edit, path)
	if err != nil {
		return err
	}

	// Release APK to track
	err = s.releaseApkToTrack(packageName, apk, []int64{apk.VersionCode}, track, edit)
	if err != nil {
		return err
	}

	// Validate edit
	return s.validateAndCommtEdit(packageName, edit)
}

func validate(packageName string, path string, track string) error {
	return nil
}

func (s *UploadApkService) createEdit(packageName string) (*androidpublisher.AppEdit, error) {
	edit, err := s.client.service.Edits.Insert(packageName, nil).Do()
	if err != nil {
		return nil, fmt.Errorf("Failed to create the edit (Error: %v)", err)
	}

	fmt.Printf("\tEdit: %v created\n", edit.Id)

	return edit, err
}

func (s *UploadApkService) uploadApk(packageName string, edit *androidpublisher.AppEdit, path string) (*androidpublisher.Apk, error) {
	fmt.Println("\tUploading APK")

	// Opening file
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return s.initateUpload(file, packageName, edit.Id)
}

func (s *UploadApkService) initateUpload(reader io.Reader, packageName string, edit string) (*androidpublisher.Apk, error) {
	// Initiate the upload
	upload := s.client.service.Edits.Apks.Upload(packageName, edit)

	// Upload the binary to edit
	apk, err := upload.Media(reader, googleapi.ContentType("application/vnd.android.package-archive")).Do()
	if err != nil {
		return nil, fmt.Errorf("Falied to upload media (Error: %v)", err)
	}

	if apk != nil {
		fmt.Printf("\tAPK Version: %v - Binary sha1: %v added to edit: %v\n",
			apk.VersionCode,
			apk.Binary.Sha1,
			edit)
	}

	return apk, nil
}

func (s *UploadApkService) releaseApkToTrack(packageName string,
	apk *androidpublisher.Apk,
	versionCodes []int64,
	trackName string,
	edit *androidpublisher.AppEdit) error {

	service := androidpublisher.NewEditsTracksService(s.client.service)

	// Resolve all available tracks for the package
	tracks, err := getAllTracks(packageName, service, edit)
	if err != nil {
		return fmt.Errorf("Failed to retrieve the tracks of package : `%v`. (Error : %v)", packageName, err)
	}

	// Resolve the track by name if existing
	track, err := resolveTrack(trackName, tracks)
	if err != nil {
		return fmt.Errorf("Failed to resolve the track : %v (Error : %v)", trackName, err)
	}

	// Create the track release
	release, err := createTrackRelease("", versionCodes, 0)
	if err != nil {
		return fmt.Errorf("Failed release the track (Error : %v)", err)
	}

	// And assign the release to the track
	track.Releases = []*androidpublisher.TrackRelease{release}

	t, err := service.Update(packageName, edit.Id, trackName, track).Do()
	if err != nil {
		return fmt.Errorf("Failed to udate the track (Error: %v)", err)
	}

	fmt.Println("\tEdit track to:", t.Track)
	return nil
}

// createTrackRelease returns a release object with the given version codes and adds the listing information.
func createTrackRelease(whatsNewsDir string, versionCodes googleapi.Int64s, userFraction float64) (*androidpublisher.TrackRelease, error) {
	newRelease := &androidpublisher.TrackRelease{
		VersionCodes: versionCodes,
		Status:       "completed",
	}

	if userFraction != 0 {
		fmt.Println("Stage rollout to ", userFraction, "% of users")
		newRelease.Status = "inProgress"
		newRelease.UserFraction = userFraction
	}

	fmt.Println("\tRelease version codes are: ", newRelease.VersionCodes)

	// Should handle release notes here

	return newRelease, nil
}

func getAllTracks(packageName string,
	editTrackService *androidpublisher.EditsTracksService,
	appEdit *androidpublisher.AppEdit) ([]*androidpublisher.Track, error) {

	response, err := editTrackService.List(packageName, appEdit.Id).Do()
	if err != nil {
		return nil, err
	}

	return response.Tracks, nil
}

func resolveTrack(name string, allTracks []*androidpublisher.Track) (*androidpublisher.Track, error) {
	for _, track := range allTracks {
		if track.Track == name {
			fmt.Printf("\nTrack found, name '%s'\n", name)
			return track, nil
		}
	}

	return nil, fmt.Errorf("could not find track with name %s", name)
}

func (s *UploadApkService) validateAndCommtEdit(packageName string, edit *androidpublisher.AppEdit) error {
	fmt.Println("\tValidating edit")
	service := androidpublisher.NewEditsService(s.client.service)

	// Validate
	_, err := service.Validate(packageName, edit.Id).Do()
	if err != nil {
		return fmt.Errorf("Failed to validate release. Error: %v", err)
	}

	// Commit
	/*
		_, err = service.Commit(packageName, edit.Id).Do()
		if err != nil {
			return fmt.Errorf("Failed to commit release. Error : %v", err)
		}
	*/

	return nil
}
