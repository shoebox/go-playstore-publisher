package playpublisher

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/androidpublisher/v3"
)

// Token structure
type Token struct {
	Email      string `json:"client_email,omitempty"`
	PrivateKey string `json:"private_key,omitempty"`
}

// Client structure
type Client struct {
	service       *androidpublisher.Service
	ListService   *ListApkService
	UploadService *UploadApkService
}

// NewClient create a new instance of the client for the provided APIKey
func NewClient(serviceAccountFile string) (*Client, error) {
	file, err := os.Open(serviceAccountFile)
	if err != nil {
		return nil, err
	}

	sa, err := resolveServiceAccount(file)
	if err != nil {
		return nil, err
	}

	//
	service, err := androidpublisher.New(sa)
	if err != nil {
		return nil, err
	}

	client := &Client{}
	client.service = service
	client.ListService = &ListApkService{client: client, editService: service.Edits}
	client.UploadService = &UploadApkService{client: client}

	return client, nil
}

// ResolveServiceAccount will resolve the http client for the designate service account file
func resolveServiceAccount(reader io.Reader) (*http.Client, error) {
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	// Unmarshal the token file
	var token Token
	json.Unmarshal(b, &token)

	//
	jwt, err := tokenToJwt(token)
	if err != nil {
		return nil, err
	}

	return jwt.Client(oauth2.NoContext), nil
}

func tokenToJwt(token Token) (*jwt.Config, error) {
	if token.Email == "" || token.PrivateKey == "" {
		return nil, fmt.Errorf("Invalid token file payload")
	}

	return &jwt.Config{
		Email:      token.Email,
		PrivateKey: []byte(token.PrivateKey),
		Scopes: []string{
			androidpublisher.AndroidpublisherScope,
		},
		TokenURL: google.JWTTokenURL,
	}, nil
}
