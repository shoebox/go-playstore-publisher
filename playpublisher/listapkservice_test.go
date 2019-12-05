package playpublisher

import (
	"testing"

	"github.com/stretchr/testify/mock"
)

type Mocks struct {
	// add a Mock object instance
	mock.Mock
}

type ListApkServiceMock struct {
	mock.Mock
}

type Service struct {
	editService *ListApkServiceMock
	// client *MockClent
}

func TestPlaceHolder(t *testing.T) {
	setup()
	defer teardown()

	// create an instance of our test object
	testObj := new(ListApkServiceMock)
	testObj.On("Insert", "com.package.name", nil).Return(true, nil)
}
