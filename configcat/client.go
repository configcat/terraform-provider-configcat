package configcat

import (
	"context"
	"fmt"

	configcatpublicapi "github.com/configcat/configcat-publicapi-go-client"
	sw "github.com/configcat/configcat-publicapi-go-client"
)

//
type Client struct {
	basePath          string
	basicAuthUsername string
	basicAuthPassword string
	apiClient         *sw.APIClient
	authEmail         string
	authFullName      string
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
	if err != nil {
		err = handleAPIError(err)
	}
	defer response.Body.Close()
	return model, err
}

func (client *Client) GetProducts() ([]sw.ProductModel, error) {
	model, response, err := client.apiClient.ProductsApi.GetProducts(client.GetAuthContext())
	if err != nil {
		err = handleAPIError(err)
	}
	defer response.Body.Close()
	return model, err
}

func (client *Client) GetConfigs(productId string) ([]sw.ConfigModel, error) {
	model, response, err := client.apiClient.ConfigsApi.GetConfigs(client.GetAuthContext(), productId)
	if err != nil {
		err = handleAPIError(err)
	}
	defer response.Body.Close()
	return model, err
}

func (client *Client) GetEnvironments(productId string) ([]sw.EnvironmentModel, error) {
	model, response, err := client.apiClient.EnvironmentsApi.GetEnvironments(client.GetAuthContext(), productId)
	if err != nil {
		err = handleAPIError(err)
	}
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
	if err != nil {
		err = handleAPIError(err)
	}
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
	if err != nil {
		err = handleAPIError(err)
	}
	defer response.Body.Close()
	return model, err
}

func (client *Client) GetSetting(settingId int32) (sw.SettingModel, error) {
	model, response, err := client.apiClient.FeatureFlagsSettingsApi.GetSetting(client.GetAuthContext(), settingId)
	if err != nil {
		err = handleAPIError(err)
	}
	defer response.Body.Close()
	return model, err
}

func (client *Client) GetSettings(configId string) ([]sw.SettingModel, error) {
	model, response, err := client.apiClient.FeatureFlagsSettingsApi.GetSettings(client.GetAuthContext(), configId)
	if err != nil {
		err = handleAPIError(err)
	}
	defer response.Body.Close()
	return model, err
}

func (client *Client) CreateSetting(configId string, body sw.CreateSettingModel) (sw.SettingModel, error) {
	model, response, err := client.apiClient.FeatureFlagsSettingsApi.CreateSetting(client.GetAuthContext(), body, configId)
	if err != nil {
		err = handleAPIError(err)
	}
	defer response.Body.Close()
	return model, err
}

func (client *Client) UpdateSetting(settingId int32, body []sw.Operation) (sw.SettingModel, error) {
	model, response, err := client.apiClient.FeatureFlagsSettingsApi.UpdateSetting(client.GetAuthContext(), body, settingId)
	if err != nil {
		err = handleAPIError(err)
	}
	defer response.Body.Close()
	return model, err
}

func (client *Client) DeleteSetting(settingId int32) error {
	response, err := client.apiClient.FeatureFlagsSettingsApi.DeleteSetting(client.GetAuthContext(), settingId)
	if err != nil {
		err = handleAPIError(err)
	}
	defer response.Body.Close()
	return err
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
		return fmt.Errorf("%s: %s", swaggerErr.Error(), string(swaggerErr.Body()))
	}
	return err
}