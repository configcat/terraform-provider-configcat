package configcat

//
type Client struct {
	basePath          string
	basicAuthUsername string
	basicAuthPassword string
	client            *APIClient
}

//
func NewClient(basePath, basicAuthUsername, basicAuthPassword) (*Client, error) {

}
