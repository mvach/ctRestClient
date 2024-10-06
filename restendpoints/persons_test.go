package restendpoints_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"ctRestClient/httpclient/httpclientfakes"
	"ctRestClient/restendpoints"
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

			personsEndpoint := restendpoints.NewPersonsEndpoint(httpClient)
			resp, err := personsEndpoint.GetPerson(5)

			Expect(err).NotTo(HaveOccurred())

			var data []map[string]interface{}
			err = json.Unmarshal(resp.Data, &data)
			Expect(err).NotTo(HaveOccurred())

			Expect(data).To(HaveLen(1))
			Expect(int(data[0]["id"].(float64))).To(Equal(5))
			Expect(data[0]["firstName"].(string)).To(Equal("foo"))
			Expect(data[0]["lastName"].(string)).To(Equal("bar"))
		})

		It("returns an error if the request cannot be send", func() {
			httpClient.DoReturns(nil, errors.New("request failed"))

			personsEndpoint := restendpoints.NewPersonsEndpoint(httpClient)
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

			personsEndpoint := restendpoints.NewPersonsEndpoint(httpClient)
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

			personsEndpoint := restendpoints.NewPersonsEndpoint(httpClient)
			_, err := personsEndpoint.GetPerson(5)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("response body is not containing expected json"))
		})
	})
})
