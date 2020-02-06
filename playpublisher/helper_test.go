package playpublisher

import (
	"io"

	"github.com/stretchr/testify/mock"
	"google.golang.org/api/androidpublisher/v3"
)

type HelperMock struct {
	mock.Mock
}

func (m *HelperMock) commitEdit(packageNameID string, editID string) error {
	args := m.Called(packageNameID, editID)
	if args.Error(0) != nil {
		return args.Error(0)
	}

	return nil
}

func (m *HelperMock) createEdit(packageNameID string) (*androidpublisher.AppEdit, error) {
	args := m.Called(packageNameID)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	value := args.String(0)

	return &androidpublisher.AppEdit{Id: value}, nil
}

func (m *HelperMock) createReleaseToTrack(changeLog string, versionCode int64, userFraction float64) *androidpublisher.TrackRelease {
	args := m.Called(changeLog, versionCode, userFraction)
	return args.Get(0).(*androidpublisher.TrackRelease)
}

func (m *HelperMock) initiateUpload(reader io.Reader, packageNameID string, editID string, contentType string) (*androidpublisher.Apk, error) {
	args := m.Called(reader, packageNameID, editID, contentType)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*androidpublisher.Apk), nil
}

func (m *HelperMock) insertEdit(packageNameID string) (*string, error) {
	args := m.Called(packageNameID)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	s := args.String(0)
	return &s, nil
}

func (m *HelperMock) listApk(packageNameID string, editID string) ([]*androidpublisher.Apk, error) {
	args := m.Called(packageNameID, editID)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*androidpublisher.Apk), nil
}

func (m *HelperMock) releaseApkToTrack(packageNameID string, apk *androidpublisher.Apk, editID string, track *androidpublisher.Track, release *androidpublisher.TrackRelease) error {
	args := m.Called(packageNameID, apk, editID, track, release)
	return args.Error(0)
}

func (m *HelperMock) resolveTrackName(packageNameID string, editID string, trackName string) (*androidpublisher.Track, error) {
	args := m.Called(packageNameID, editID, trackName)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*androidpublisher.Track), nil
}

func (m *HelperMock) validateEdit(packageNameID string, editID string) error {
	args := m.Called(packageNameID, editID)
	return args.Error(0)
}
