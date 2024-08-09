package client

import (
	sw "github.com/configcat/configcat-publicapi-go-client/v2"
)

func (client *Client) GetProducts() ([]sw.ProductModel, error) {
	model, response, err := client.apiClient.ProductsAPI.GetProducts(client.GetAuthContext()).Execute()
	error := handleAPIError(err)
	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}
	return model, error
}

func (client *Client) GetProduct(productID string) (*sw.ProductModel, error) {
	model, response, err := client.apiClient.ProductsAPI.GetProduct(client.GetAuthContext(), productID).Execute()
	error := handleAPIError(err)
	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}
	return model, error
}

func (client *Client) CreateProduct(organizationID string, body sw.CreateProductRequest) (*sw.ProductModel, error) {
	model, response, err := client.apiClient.ProductsAPI.CreateProduct(
		client.GetAuthContext(),
		organizationID).CreateProductRequest(body).Execute()
	error := handleAPIError(err)
	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}
	return model, error
}

func (client *Client) UpdateProduct(productID string, body sw.UpdateProductRequest) (*sw.ProductModel, error) {
	model, response, err := client.apiClient.ProductsAPI.UpdateProduct(
		client.GetAuthContext(),
		productID).UpdateProductRequest(body).Execute()
	error := handleAPIError(err)
	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}
	return model, error
}

func (client *Client) DeleteProduct(productID string) error {
	response, err := client.apiClient.ProductsAPI.DeleteProduct(
		client.GetAuthContext(),
		productID).Execute()
	error := handleAPIError(err)
	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}
	return error
}

func (client *Client) GetProductPreferences(productID string) (*sw.PreferencesModel, error) {
	model, response, err := client.apiClient.ProductsAPI.GetProductPreferences(client.GetAuthContext(), productID).Execute()
	error := handleAPIError(err)
	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}
	return model, error
}

func (client *Client) UpdateProductPreferences(productID string, body sw.UpdatePreferencesRequest) (*sw.PreferencesModel, error) {
	model, response, err := client.apiClient.ProductsAPI.UpdateProductPreferences(
		client.GetAuthContext(),
		productID).UpdatePreferencesRequest(body).Execute()
	error := handleAPIError(err)
	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}
	return model, error
}
