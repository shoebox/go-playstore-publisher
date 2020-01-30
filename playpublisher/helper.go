package playpublisher

import (
	"fmt"
	"io"

	"google.golang.org/api/androidpublisher/v3"
	"google.golang.org/api/googleapi"
)

type Helper interface {
	createEdit(packageName string) (string, error)
	insertEdit(packageNameID string) (*string, error)
	listApk(packageNameID string, editID string) ([]*androidpublisher.Apk, error)
}

type AndroidPublisherHelper struct {
	service          *androidpublisher.Service
	editService      *androidpublisher.EditsService
	editTrackService *androidpublisher.EditsTracksService
}

func NewHelper(service *androidpublisher.Service) *AndroidPublisherHelper {
	res := AndroidPublisherHelper{
		service:          service,
		editService:      androidpublisher.NewEditsService(service),
		editTrackService: androidpublisher.NewEditsTracksService(service),
	}

	return &res
}

func (h *AndroidPublisherHelper) createEdit(packageNameId string) (*androidpublisher.AppEdit, error) {
	edit, err := h.service.Edits.Insert(packageNameId, nil).Do()
	if err != nil {
		return nil, fmt.Errorf("Failed to create the edit (Error: %v)", err)
	}

	return edit, nil
}

func (h *AndroidPublisherHelper) insertEdit(packageNameId string) (*string, error) {
	edit, err := h.createEdit(packageNameId)
	if err != nil {
		return nil, err
	}
	return &edit.Id, nil
}

func (h *AndroidPublisherHelper) listApk(packageNameID string,
	editID *string) ([]*androidpublisher.Apk, error) {
	call, err := h.service.Edits.Apks.List(packageNameID, *editID).Do()
	if err != nil {
		return nil, err
	}

	return call.Apks, nil
}

func (h *AndroidPublisherHelper) initiateUpload(reader io.Reader,
	packageNameID string,
	editID string,
	contentType string) (*androidpublisher.Apk, error) {
	edit := h.service.Edits.Apks.Upload(packageNameID, editID)

	apk, err := edit.Media(reader, googleapi.ContentType(contentType)).Do()
	if err != nil {
		return nil, fmt.Errorf("Falied to initiate upload (Error: %v)", err)
	}

	return apk, nil
}

func (h *AndroidPublisherHelper) resolveTrackName(packageNameID string,
	editID string,
	trackName string) (*androidpublisher.Track, error) {

	response, err := h.editTrackService.List(packageNameID, editID).Do()
	if err != nil {
		return nil, err
	}

	for _, track := range response.Tracks {
		if track.Track == trackName {
			return track, nil
		}
	}

	return nil, fmt.Errorf("could not find track with name %s", trackName)
}

func (h *AndroidPublisherHelper) createReleaseToTrack(changeLog string,
	versionCode int64,
	userFraction float64) *androidpublisher.TrackRelease {

	newRelease := &androidpublisher.TrackRelease{
		VersionCodes: []int64{versionCode},
		Status:       "completed",
	}

	if userFraction != 0 {
		fmt.Println("Stage rollout to ", userFraction, "% of users")
		newRelease.Status = "inProgress"
		newRelease.UserFraction = userFraction
	}

	fmt.Println("\tRelease version codes are: ", newRelease.VersionCodes)

	return newRelease
}

func (h *AndroidPublisherHelper) releaseApkToTrack(packageNameID string,
	apk *androidpublisher.Apk,
	editID string,
	track *androidpublisher.Track,
	release *androidpublisher.TrackRelease) error {

	// And assign the release to the track
	track.Releases = []*androidpublisher.TrackRelease{release}

	_, err := h.editTrackService.Update(packageNameID, editID, track.Track, track).Do()
	if err != nil {
		return fmt.Errorf("Failed to update the track (Error: %v)", err)
	}

	return nil
}

func (h *AndroidPublisherHelper) validateEdit(packageNameID string, editID string) error {
	_, err := h.editService.Validate(packageNameID, editID).Do()
	if err != nil {
		return fmt.Errorf("Failed to validate release. Error: %v", err)
	}

	return nil
}

func (h *AndroidPublisherHelper) commitEdit(packageNameID string, editID string) error {

	appEdit, err := h.editService.Commit(packageNameID, editID).Do()
	fmt.Println(appEdit, err)

	fmt.Println("code", appEdit)
	fmt.Println("code", appEdit.HTTPStatusCode)

	return err
}
