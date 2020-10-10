package configcat

import (
	"context"
	"log"

	configcatpublicapi "github.com/configcat/configcat-publicapi-go-client"
	sw "github.com/configcat/configcat-publicapi-go-client"
)

//
type Client struct {
	basePath          string
	basicAuthUsername string
	basicAuthPassword string
	apiClient         *sw.APIClient
}

// test
func (a *Client) GetAuthContext() context.Context {
	return context.WithValue(context.Background(), sw.ContextBasicAuth, sw.BasicAuth{
		UserName: a.basicAuthUsername,
		Password: a.basicAuthPassword,
	})
}

//
func NewClient(basePath, basicAuthUsername, basicAuthPassword string) (*Client, error) {
	configuration := configcatpublicapi.NewConfiguration()
	configuration.BasePath = basePath
	configuration.UserAgent = "terraform-provider-configcat/1.0.0"
	apiClient := configcatpublicapi.NewAPIClient(configuration)

	client := &Client{
		basePath:          basePath,
		basicAuthUsername: basicAuthUsername,
		basicAuthPassword: basicAuthPassword,
		apiClient:         apiClient,
	}

	log.Printf("Validating ConfigCat configuration")
	meModel, _, err := client.apiClient.MeApi.GetMe(client.GetAuthContext())

	if err != nil {
		return nil, err
	}
	log.Printf("ConfigCat provider authorized with %s %s", meModel.Email, meModel.FullName)

	return client, nil
}
