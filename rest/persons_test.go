package rest_test

import (
    "bytes"
    "encoding/json"
    "errors"
    "io"
    "net/http"

    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"

    "ctRestClient/httpclient/httpclientfakes"
    "ctRestClient/rest"
)

var _ = Describe("PersonsEndpoint", func() {

    var (
        httpClient *httpclientfakes.FakeHTTPClient
    )

    BeforeEach(func() {
        httpClient = &httpclientfakes.FakeHTTPClient{}
    })

    var _ = Describe("GetPerson", func() {

        It("returns groupname  ...", func() {
            httpResponse := &http.Response{
                StatusCode: 200,
                Body: io.NopCloser(bytes.NewBufferString(
                    `{
                        "data": [
                            {
                                "id": 5,
                                "firstName": "foo",
                                "lastName": "bar"
                            }
                        ]
                    }`))}
            httpClient.DoReturns(httpResponse, nil)

            personsEndpoint := rest.NewPersonsEndpoint(httpClient)
            resp, err := personsEndpoint.GetPerson(5)

            Expect(err).NotTo(HaveOccurred())

            var data map[string]interface{}
            err = json.Unmarshal(resp[0], &data)
            Expect(err).NotTo(HaveOccurred())

            Expect(int(data["id"].(float64))).To(Equal(5))
            Expect(data["firstName"].(string)).To(Equal("foo"))
            Expect(data["lastName"].(string)).To(Equal("bar"))
        })

        It("returns an error if the request cannot be send", func() {
            httpClient.DoReturns(nil, errors.New("request failed"))

            personsEndpoint := rest.NewPersonsEndpoint(httpClient)
            _, err := personsEndpoint.GetPerson(5)

            Expect(err).To(HaveOccurred())
            Expect(err.Error()).To(Equal("failed to send request, request failed"))
        })

        It("returns an error if the status code is wrong", func() {
            httpResponse := &http.Response{
                StatusCode: 404,
                Body: io.NopCloser(bytes.NewBufferString(
                    `{
                        "data": []
                    }`))}
            httpClient.DoReturns(httpResponse, nil)

            personsEndpoint := rest.NewPersonsEndpoint(httpClient)
            _, err := personsEndpoint.GetPerson(5)

            Expect(err).To(HaveOccurred())
            Expect(err.Error()).To(Equal("received non-200 response code: 404"))
        })

        It("returns an error if the response body is not a church tools json response", func() {
            httpResponse := &http.Response{
                StatusCode: 200,
                Body: io.NopCloser(bytes.NewBufferString(
                    `{
                        "foo": [],
                    }`))}
            httpClient.DoReturns(httpResponse, nil)

            personsEndpoint := rest.NewPersonsEndpoint(httpClient)
            _, err := personsEndpoint.GetPerson(5)

            Expect(err).To(HaveOccurred())
            Expect(err.Error()).To(ContainSubstring("response body is not containing expected json"))
        })
    })
})
