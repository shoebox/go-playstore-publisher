package playpublisher

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
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
