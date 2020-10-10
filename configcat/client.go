package configcat

import configcatpublicapi "github.com/configcat/configcat-publicapi-go-client/configcat-publicapi-go-client"

//
type Client struct {
	basePath           string
	basicAuthUsername string
	basicAuthPassword string
	client            *configcatpublicapi.APIClient
}

//
func NewClient(basePath, basicAuthUsername, basicAuthPassword: string) (*Client, error) {
	publicApiClient, error := &configcatpublicapi.NewAPIClient(
		&configcatpublicapi.NewConfiguration{
			BasePath: strings.TrimSpace(basePath),
			UserAgent: "terraform-provider-configcat/1.0.0"
		}
	)


	newClient, error := &Client{
		basePath: strings.TrimSpace(basePath),
		basicAuthUsername: strings.TrimSpace(basicAuthUsername),
		basicAuthPassword:  strings.TrimSpace(basicAuthPassword),
		client: &configcatpublicapi.NewAPIClient(
			&configcatpublicapi.NewConfiguration{
				BasePath: strings.TrimSpace(basePath),
				UserAgent: "terraform-provider-configcat/1.0.0"
			}
		)
	}, nil;

	return newClient, nil
}
