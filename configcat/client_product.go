package configcat

import (
	sw "github.com/configcat/configcat-publicapi-go-client"
)

func (client *Client) GetProducts() ([]sw.ProductModel, error) {
	model, response, err := client.apiClient.ProductsApi.GetProducts(client.GetAuthContext())
	defer response.Body.Close()
	return model, handleAPIError(err)
}

func (client *Client) GetProduct(productID string) (sw.ProductModel, error) {
	model, response, err := client.apiClient.ProductsApi.GetProduct(client.GetAuthContext(), productID)
	defer response.Body.Close()
	return model, handleAPIError(err)
}

func (client *Client) CreateProduct(organizationID string, body sw.CreateProductRequest) (sw.ProductModel, error) {
	model, response, err := client.apiClient.ProductsApi.CreateProduct(
		client.GetAuthContext(),
		body,
		organizationID)
	defer response.Body.Close()
	return model, handleAPIError(err)
}

func (client *Client) UpdateProduct(productID string, body sw.UpdateProductRequest) (sw.ProductModel, error) {
	model, response, err := client.apiClient.ProductsApi.UpdateProduct(
		client.GetAuthContext(),
		body,
		productID)
	defer response.Body.Close()
	return model, handleAPIError(err)
}

func (client *Client) DeleteProduct(productID string) error {
	response, err := client.apiClient.ProductsApi.DeleteProduct(
		client.GetAuthContext(),
		productID)
	defer response.Body.Close()
	return handleAPIError(err)
}
