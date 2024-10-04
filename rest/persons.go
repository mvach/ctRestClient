package rest

import (
    "ctRestClient/httpclient"
    "encoding/json"
    "fmt"
    "io"

    "net/http"
    "net/url"
)

//counterfeiter:generate . PersonsEndpoint
type PersonsEndpoint interface {
    GetPerson(personId int) ([]json.RawMessage, error)
}

type personsEndpoint struct {
    httpclient httpclient.HTTPClient
}

func NewPersonsEndpoint(httpclient httpclient.HTTPClient) PersonsEndpoint {
    return personsEndpoint{
        httpclient: httpclient,
    }
}

func (c personsEndpoint) GetPerson(personId int) ([]json.RawMessage, error) {

    req, err := http.NewRequest("GET", "", nil)
    if err != nil {
        return nil, fmt.Errorf("failed to create request, %w", err)
    }

    params := url.Values{}
    params.Add("ids[]", fmt.Sprintf("%d", personId))
    encodedQueryParam := params.Encode()

    req.URL.Path = "/api/persons"
    req.URL.RawQuery = encodedQueryParam

    resp, err := c.httpclient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to send request, %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("received non-200 response code: %d", resp.StatusCode)
    }

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read response body, %w", err)
    }

    var response PersonResponseJson
    if err := json.Unmarshal(body, &response); err != nil {
        return nil, fmt.Errorf("response body is not containing expected json, %w", err)
    }

    return response.Data, nil
}
