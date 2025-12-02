package rest_test

import (
	"bytes"
	"errors"
	"io"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"ctRestClient/httpclient/httpclientfakes"
	"ctRestClient/rest"
)

var _ = Describe("DynamicGroupsEndpoint", func() {

	var (
		httpClient   *httpclientfakes.FakeHTTPClient
		httpResponse *http.Response
	)

	BeforeEach(func() {
		httpClient = &httpclientfakes.FakeHTTPClient{}

		httpResponse = &http.Response{
			StatusCode: 200,
			Body: io.NopCloser(bytes.NewBufferString(
				`{  
  					"dynamicGroupStatus": "active"
                }`)),
		}
	})

	var _ = Describe("GetGroupStatus", func() {

		It("returns a group status", func() {

			httpClient.DoReturns(httpResponse, nil)

			endpoint := rest.NewDynamicGroupsEndpoint(httpClient)
			group, err := endpoint.GetGroupStatus(1)

			Expect(err).NotTo(HaveOccurred())
			Expect(*group.Status).To(Equal("active"))
			request := httpClient.DoArgsForCall(0)
			Expect(request.URL.Path).To(Equal("/api/dynamicgroups/1/status"))
		})

		It("returns an error if the request cannot be send", func() {
			httpClient.DoReturns(nil, errors.New("request failed"))

			endpoint := rest.NewDynamicGroupsEndpoint(httpClient)
			_, err := endpoint.GetGroupStatus(1)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("failed to send request, request failed"))
		})

		It("returns an error if the status code is wrong", func() {
			httpResponse := &http.Response{
				StatusCode: 404,
				Body:       io.NopCloser(bytes.NewBufferString(`{}`))}
			httpClient.DoReturns(httpResponse, nil)

			endpoint := rest.NewDynamicGroupsEndpoint(httpClient)
			_, err := endpoint.GetGroupStatus(1)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("received non-200 response code: 404"))
		})

		It("returns an error if the response body is an invalid json response", func() {
			httpResponse := &http.Response{
				StatusCode: 200,
				Body: io.NopCloser(bytes.NewBufferString(
					``))}
			httpClient.DoReturns(httpResponse, nil)

			endpoint := rest.NewDynamicGroupsEndpoint(httpClient)
			_, err := endpoint.GetGroupStatus(1)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("response body is not containing expected json"))
		})

		It("returns an error if the response body is missing the dynamicGroupStatus field", func() {
			httpResponse := &http.Response{
				StatusCode: 200,
				Body: io.NopCloser(bytes.NewBufferString(
					`{"foo": []}`))}
			httpClient.DoReturns(httpResponse, nil)

			endpoint := rest.NewDynamicGroupsEndpoint(httpClient)
			_, err := endpoint.GetGroupStatus(1)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("response body is missing dynamicGroupStatus field"))
		})
	})
})
