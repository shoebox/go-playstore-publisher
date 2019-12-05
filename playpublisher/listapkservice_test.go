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

func TestTotototo(t *testing.T) {
	setup()
	defer teardown()

	// create an instance of our test object
	testObj := new(ListApkServiceMock)
	testObj.On("Insert", "com.package.name", nil).Return(true, nil)

	// fmt.Println(">>>>", m)

	/*
		service := ListApkService{
		}

		service.List("com.toto.tutu")
	*/

	/*
		mux.HandleFunc("/api/v1/feeds/temperature/data",
			func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "POST")
				fmt.Fprint(w, `{"id":1, "value":"67.112"}`)
			},
		)

		assert := assert.New(t)

		dp := &Data{}
		datapoint, response, err := client.Data.Create(dp)

		assert.NotNil(err)
		assert.Nil(datapoint)
		assert.Nil(response)

		assert.Equal(err.Error(), "CurrentFeed must be set")
	*/
}
