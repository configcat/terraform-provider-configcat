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

//
func (client *Client) GetAuthContext() context.Context {
	return context.WithValue(context.Background(), sw.ContextBasicAuth, sw.BasicAuth{
		UserName: client.basicAuthUsername,
		Password: client.basicAuthPassword,
	})
}

//
func (client *Client) GetMe() (sw.MeModel, error) {
	model, _, err := client.apiClient.MeApi.GetMe(client.GetAuthContext())
	return model, err
}

func (client *Client) GetProducts() ([]sw.ProductModel, error) {
	model, _, err := client.apiClient.ProductsApi.GetProducts(client.GetAuthContext())
	return model, err
}

func (client *Client) GetConfigs(productId string) ([]sw.ConfigModel, error) {
	model, _, err := client.apiClient.ConfigsApi.GetConfigs(client.GetAuthContext(), productId)
	return model, err
}

func (client *Client) GetEnvironments(productId string) ([]sw.EnvironmentModel, error) {
	model, _, err := client.apiClient.EnvironmentsApi.GetEnvironments(client.GetAuthContext(), productId)
	return model, err
}

func (client *Client) CreateEnvironment(productId, environmentName string) (sw.EnvironmentModel, error) {
	body := sw.CreateEnvironmentModel{}
	body.Name = environmentName
	model, _, err := client.apiClient.EnvironmentsApi.CreateEnvironment(
		client.GetAuthContext(),
		body,
		productId)
	return model, err
}

func (client *Client) UpdateEnvironment(environmentId, environmentName string) (sw.EnvironmentModel, error) {
	body := sw.UpdateEnvironmentModel{}
	body.Name = environmentName
	model, _, err := client.apiClient.EnvironmentsApi.UpdateEnvironment(
		client.GetAuthContext(),
		body,
		environmentId)
	return model, err
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
	meModel, err := client.GetMe()
	if err != nil {
		return nil, err
	}
	log.Printf("ConfigCat provider authorized with %s %s", meModel.Email, meModel.FullName)

	return client, nil
}
