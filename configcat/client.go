package configcat

import (
	"context"
	"fmt"

	configcatpublicapi "github.com/configcat/configcat-publicapi-go-client"
	sw "github.com/configcat/configcat-publicapi-go-client"
)

type Client struct {
	basePath          string
	basicAuthUsername string
	basicAuthPassword string
	apiClient         *sw.APIClient
	authEmail         string
	authFullName      string
}

func (client *Client) GetAuthContext() context.Context {
	return context.WithValue(context.Background(), sw.ContextBasicAuth, sw.BasicAuth{
		UserName: client.basicAuthUsername,
		Password: client.basicAuthPassword,
	})
}

func (client *Client) GetMe() (sw.MeModel, error) {
	model, response, err := client.apiClient.MeApi.GetMe(client.GetAuthContext())
	defer response.Body.Close()
	return model, handleAPIError(err)
}

func (client *Client) GetOrganizations() ([]sw.OrganizationModel, error) {
	model, response, err := client.apiClient.OrganizationsApi.GetOrganizations(client.GetAuthContext())
	defer response.Body.Close()
	return model, handleAPIError(err)
}

func NewClient(basePath, basicAuthUsername, basicAuthPassword string) (*Client, error) {
	configuration := configcatpublicapi.NewConfiguration()
	configuration.BasePath = basePath
	configuration.UserAgent = "terraform-provider-configcat/1.0.3"
	apiClient := configcatpublicapi.NewAPIClient(configuration)

	client := &Client{
		basePath:          basePath,
		basicAuthUsername: basicAuthUsername,
		basicAuthPassword: basicAuthPassword,
		apiClient:         apiClient,
	}

	meModel, err := client.GetMe()
	if err != nil {
		return nil, err
	}
	client.authEmail = meModel.Email
	client.authFullName = meModel.FullName

	return client, nil
}

func handleAPIError(err error) error {
	if err == nil {
		return nil
	}
	if swaggerErr, ok := err.(sw.GenericSwaggerError); ok {
		if swaggerErr.Error() == "404 Not Found" {
			return NotFoundError{
				error: swaggerErr.Error(),
				body:  string(swaggerErr.Body()),
			}
		}
		return fmt.Errorf("%s: %s", swaggerErr.Error(), string(swaggerErr.Body()))
	}
	return err
}

// NotFoundError Provides access to the body, error in case of 404 Not Found.
type NotFoundError struct {
	error string
	body  string
}

// Error returns non-empty string if there was an error.
func (e NotFoundError) Error() string {
	return e.error
}

// Body returns non-empty string if there was an error.
func (e NotFoundError) Body() string {
	return e.body
}
