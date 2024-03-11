package bitpanda

type ApiClient interface {
}

type identityApiClient struct {
	apiKey string
}

func NewApiClient(apiKey string) ApiClient {
	return &identityApiClient{apiKey: apiKey}
}
