package playpublisher

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/mock"
	"google.golang.org/api/androidpublisher/v3"
)

type listApkServiceMock struct {
	mock.Mock
}

func (m *listApkServiceMock) insertEdit(packageNameID string) (string, error) {
	return "", nil
}

func (m *listApkServiceMock) listApk(packageNameID string,
	editID *string) ([]*androidpublisher.Apk, error) {
	return nil, nil
}

func TestListing(t *testing.T) {
	service := new(listApkServiceMock)
	id, err := service.insertEdit("wffew")
	fmt.Println("test listing ", id, " ::: ", err)
}
