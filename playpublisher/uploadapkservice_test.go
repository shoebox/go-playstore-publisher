package playpublisher

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/api/androidpublisher/v3"
)

var (
	subject UploadApkService
)

var hm *HelperMock

func setupMock() {
	c := new(Client)
	hm = new(HelperMock)
	c.Helper = hm

	subject = UploadApkService{client: c}
}

func TestCreateEdit(t *testing.T) {
	t.Run("Should handle errors", func(t *testing.T) {
		// setup:
		setupMock()
		subject.packageNameID = "co.test.testing"

		hm.On("createEdit", subject.packageNameID).
			Return(nil, fmt.Errorf("Test error"))

		// when:
		err := subject.createEdit()

		// then:
		assert.EqualError(t, err, "Failed to create edit (Error : Test error)")

		// and:
		hm.AssertExpectations(t)
	})

	t.Run("Should apply the edit ID", func(t *testing.T) {
		// setup:
		setupMock()
		subject.packageNameID = "co.test.testing"
		hm.On("createEdit", subject.packageNameID).
			Return("123456XYZ", nil)

		// when:
		err := subject.createEdit()

		// then:
		assert.NoError(t, err)
		assert.EqualValues(t, "123456XYZ", subject.editID)

		// and:
		hm.AssertExpectations(t)
	})
}

func TestUploadBinary(t *testing.T) {
	subject.editID = "123456XYZ"
	reader := strings.NewReader("Hello world")

	t.Run("Binary upload should handle errors", func(t *testing.T) {
		// setup:
		setupMock()
		subject.packageNameID = "co.test.testing"

		hm.On("initiateUpload", reader, subject.packageNameID, subject.editID, "application/vnd.android.package-archive").
			Return(nil, fmt.Errorf("Test error"))

		edit, err := subject.uploadBinary(subject.packageNameID, reader)
		assert.EqualError(t, err, "Test error")
		assert.Nil(t, edit)
	})

	t.Run("Binary upload should return apk", func(t *testing.T) {
		// setup:
		setupMock()
		subject.packageNameID = "co.test.testing"

		a := &androidpublisher.Apk{VersionCode: 1234}

		hm.On("initiateUpload", reader, subject.packageNameID, subject.editID, "application/vnd.android.package-archive").
			Return(a, nil)

		apk, err := subject.uploadBinary(subject.packageNameID, reader)
		assert.EqualValues(t, apk, apk)
		assert.NoError(t, err)
	})
}
