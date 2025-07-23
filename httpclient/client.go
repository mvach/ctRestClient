package httpclient

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
	"strings"
)

//counterfeiter:generate . HTTPClient
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type httpClient struct {
	client    *http.Client
	hostname  string
	authToken string
}

func NewHTTPClient(hostname string, authToken string) HTTPClient {
	client := &http.Client{}

	if os.Getenv("ALLOW_SELF_SIGNED_CERTS") == "true" {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				// Skip the default verification so we can do custom verification
				InsecureSkipVerify: true,
				VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
					// Only allow self-signed certificates for localhost and 127.0.0.1
					if strings.HasPrefix(hostname, "127.0.0.1") || strings.HasPrefix(hostname, "localhost") {
						return nil
					}
					// For production hosts, reject self-signed certificates
					return fmt.Errorf("self-signed certificates not allowed for host: %s", hostname)
				},
				MinVersion: tls.VersionTLS12,
			},
		}
	}

	return httpClient{
		client:    client,
		authToken: authToken,
		hostname:  hostname,
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
