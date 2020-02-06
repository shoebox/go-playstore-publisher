package playpublisher

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/mock"
	"google.golang.org/api/androidpublisher/v3"
)

const (
	fakeEmail      = "fake@email.com"
	fakePrivateKey = "PrivateKey"
)

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("Invalid payload")
}

var client *Client

type mockAndroidPublisherService struct {
	mock.Mock
}

func (m *mockAndroidPublisherService) SendChargeNotification(packageNameid string,
	appedit *androidpublisher.AppEdit) *androidpublisher.EditsInsertCall {
	return nil
}

func setup() error {
	return nil
}

func teardown() {
}

func TestServiceAccountResolutio(t *testing.T) {
	t.Run("When the service account file is badly formed should return an error", func(t *testing.T) {
		c, err := resolveServiceAccount(strings.NewReader("toto"))
		assert.Error(t, err, "Invald token file payload")
		assert.Nil(t, c)
	})

	t.Run("When the reader input is invald should return an error", func(t *testing.T) {
		c, err := resolveServiceAccount(errReader(0))
		assert.EqualError(t, err, "Invalid payload")
		assert.Nil(t, c)
	})

	t.Run("Should return the JWT token", func(t *testing.T) {
		json := fmt.Sprintf("{\"client_email\":\"%v\", \"private_key\":\"%v\"}", fakeEmail, fakePrivateKey)

		c, err := resolveServiceAccount(strings.NewReader(json))
		assert.NotNil(t, c)
		assert.NoError(t, err)
	})

	t.Run("JWT Token parsing failing should report error", func(t *testing.T) {
		json := fmt.Sprintf("{\"client_email\":\"%v\"}", fakeEmail)

		c, err := resolveServiceAccount(strings.NewReader(json))
		assert.Nil(t, c)
		assert.Error(t, err, "Invalid token file payload")
	})
}

func TestInputServiceAccount(t *testing.T) {
	token := Token{}

	t.Run("When token is invalid or empty, the jwt token should return an error", func(t *testing.T) {
		jwt, err := tokenToJwt(token)
		assert.Error(t, err, "Invalid token file payload")
		assert.Nil(t, jwt)
	})

	t.Run("When the token is valid, the jwt token should be valid", func(t *testing.T) {
		token.Email = fakeEmail
		token.PrivateKey = fakePrivateKey

		jwt, err := tokenToJwt(token)
		assert.NoError(t, err)
		assert.NotNil(t, jwt)

		assert.Equal(t, jwt.Email, fakeEmail)
		assert.Equal(t, jwt.PrivateKey, []byte(fakePrivateKey))
	})
}

func TestInitClient(t *testing.T) {
	h := http.Client{}

	t.Run("Create new client should work without error", func(t *testing.T) {
		client, err := initClient(&h)
		assert.NoError(t, err)
		assert.NotNil(t, client)
	})

	t.Run("Create new client with invalid file should throw an error", func(t *testing.T) {
		client, err := NewClient("toto.json")
		assert.Error(t, err, "pen toto.json: no such file or directory")
		assert.Nil(t, client)
	})
}
