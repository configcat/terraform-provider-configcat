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
	model, response, err := client.apiClient.MeApi.GetMe(client.GetAuthContext())
	defer response.Body.Close()
	return model, err
}

func (client *Client) GetProducts() ([]sw.ProductModel, error) {
	model, response, err := client.apiClient.ProductsApi.GetProducts(client.GetAuthContext())
	defer response.Body.Close()
	return model, err
}

func (client *Client) GetConfigs(productId string) ([]sw.ConfigModel, error) {
	model, response, err := client.apiClient.ConfigsApi.GetConfigs(client.GetAuthContext(), productId)
	defer response.Body.Close()
	return model, err
}

func (client *Client) GetEnvironments(productId string) ([]sw.EnvironmentModel, error) {
	model, response, err := client.apiClient.EnvironmentsApi.GetEnvironments(client.GetAuthContext(), productId)
	defer response.Body.Close()
	return model, err
}

func (client *Client) CreateEnvironment(productId, environmentName string) (sw.EnvironmentModel, error) {
	body := sw.CreateEnvironmentModel{}
	body.Name = environmentName
	model, response, err := client.apiClient.EnvironmentsApi.CreateEnvironment(
		client.GetAuthContext(),
		body,
		productId)
	defer response.Body.Close()
	return model, err
}

func (client *Client) UpdateEnvironment(environmentId, environmentName string) (sw.EnvironmentModel, error) {
	body := sw.UpdateEnvironmentModel{}
	body.Name = environmentName
	model, response, err := client.apiClient.EnvironmentsApi.UpdateEnvironment(
		client.GetAuthContext(),
		body,
		environmentId)
	defer response.Body.Close()
	return model, err
}

func (client *Client) GetSetting(settingId int32) (sw.SettingModel, error) {
	model, response, err := client.apiClient.FeatureFlagsSettingsApi.GetSetting(client.GetAuthContext(), settingId)
	defer response.Body.Close()
	return model, err
}

func (client *Client) GetSettings(configId string) ([]sw.SettingModel, error) {
	model, response, err := client.apiClient.FeatureFlagsSettingsApi.GetSettings(client.GetAuthContext(), configId)
	defer response.Body.Close()
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
