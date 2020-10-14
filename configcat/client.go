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
	defer response.Body.Close()
	return model, handleAPIError(err)
}

func (client *Client) GetProducts() ([]sw.ProductModel, error) {
	model, response, err := client.apiClient.ProductsApi.GetProducts(client.GetAuthContext())
	defer response.Body.Close()
	return model, handleAPIError(err)
}

func (client *Client) GetConfigs(productId string) ([]sw.ConfigModel, error) {
	model, response, err := client.apiClient.ConfigsApi.GetConfigs(client.GetAuthContext(), productId)
	defer response.Body.Close()
	return model, handleAPIError(err)
}

func (client *Client) GetEnvironments(productId string) ([]sw.EnvironmentModel, error) {
	model, response, err := client.apiClient.EnvironmentsApi.GetEnvironments(client.GetAuthContext(), productId)
	defer response.Body.Close()
	return model, handleAPIError(err)
}

func (client *Client) CreateEnvironment(productId, environmentName string) (sw.EnvironmentModel, error) {
	body := sw.CreateEnvironmentModel{}
	body.Name = environmentName
	model, response, err := client.apiClient.EnvironmentsApi.CreateEnvironment(
		client.GetAuthContext(),
		body,
		productId)
	defer response.Body.Close()
	return model, handleAPIError(err)
}

func (client *Client) UpdateEnvironment(environmentId, environmentName string) (sw.EnvironmentModel, error) {
	body := sw.UpdateEnvironmentModel{}
	body.Name = environmentName
	model, response, err := client.apiClient.EnvironmentsApi.UpdateEnvironment(
		client.GetAuthContext(),
		body,
		environmentId)
	defer response.Body.Close()
	return model, handleAPIError(err)
}

func (client *Client) GetSetting(settingId int32) (sw.SettingModel, error) {
	model, response, err := client.apiClient.FeatureFlagsSettingsApi.GetSetting(client.GetAuthContext(), settingId)
	defer response.Body.Close()
	return model, handleAPIError(err)
}

func (client *Client) GetSettings(configId string) ([]sw.SettingModel, error) {
	model, response, err := client.apiClient.FeatureFlagsSettingsApi.GetSettings(client.GetAuthContext(), configId)
	defer response.Body.Close()
	return model, handleAPIError(err)
}

func (client *Client) CreateSetting(configId string, body sw.CreateSettingModel) (sw.SettingModel, error) {
	model, response, err := client.apiClient.FeatureFlagsSettingsApi.CreateSetting(client.GetAuthContext(), body, configId)
	defer response.Body.Close()
	return model, handleAPIError(err)
}

func (client *Client) UpdateSetting(settingId int32, body []sw.Operation) (sw.SettingModel, error) {
	model, response, err := client.apiClient.FeatureFlagsSettingsApi.UpdateSetting(client.GetAuthContext(), body, settingId)
	defer response.Body.Close()
	return model, handleAPIError(err)
}

func (client *Client) DeleteSetting(settingId int32) error {
	response, err := client.apiClient.FeatureFlagsSettingsApi.DeleteSetting(client.GetAuthContext(), settingId)
	defer response.Body.Close()
	return handleAPIError(err)
}

func (client *Client) GetSettingValue(environmentId string, settingId int32) (sw.SettingValueModel, error) {
	model, response, err := client.apiClient.FeatureFlagSettingValuesApi.GetSettingValue(client.GetAuthContext(), environmentId, settingId)
	defer response.Body.Close()
	return model, handleAPIError(err)
}

func (client *Client) ReplaceSettingValue(environmentId string, settingId int32, body sw.UpdateSettingValueModel) (sw.SettingValueModel, error) {
	model, response, err := client.apiClient.FeatureFlagSettingValuesApi.ReplaceSettingValue(client.GetAuthContext(), body, environmentId, settingId, nil)
	defer response.Body.Close()
	return model, handleAPIError(err)
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
