package httpclient

import "net/http"

//counterfeiter:generate . HTTPClient
type HTTPClient interface {
    Do(req *http.Request) (*http.Response, error)
}

type httpClient struct {
    client    *http.Client
    hostname   string
    authToken string
}

func NewHTTPClient(hostname string, authToken string) HTTPClient {
    return httpClient{
        client:    &http.Client{},
        authToken: authToken,
        hostname:   hostname,
    }
}

func (c httpClient) Do(req *http.Request) (*http.Response, error) {
    // Set common headers
    req.Header.Set("Accept", "application/json")
    req.Header.Set("Authorization", "Login "+c.authToken)

    // Construct the full URL
    req.URL.Scheme = "https"
    req.URL.Host = c.hostname

    // Perform the request
    return c.client.Do(req)
}
