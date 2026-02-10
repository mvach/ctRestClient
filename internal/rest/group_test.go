package rest_test

import (
	"errors"
	"io"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"ctRestClient/internal/httpclient/httpclientfakes"
	"ctRestClient/internal/rest"
	"ctRestClient/internal/testutil"
)

var _ = Describe("GroupsEndpoint", func() {

	var (
		httpClient   *httpclientfakes.FakeHTTPClient
		httpResponse *http.Response
	)

	BeforeEach(func() {
		httpClient = &httpclientfakes.FakeHTTPClient{}

		httpResponse = &http.Response{
			StatusCode: 200,
			Body: io.NopCloser(testutil.JsonToBufferString(
				`{
					"data": {
						"id": 6,
						"name": "Dynamische Gruppe"
					}
				}`)),
		}
	})

	var _ = Describe("GetGroupType", func() {

		It("returns a groupType", func() {
			httpClient.DoReturns(httpResponse, nil)

			groupEndpoint := rest.NewGroupEndpoint(httpClient)
			groupType, err := groupEndpoint.GetGroupType(6)

			Expect(err).NotTo(HaveOccurred())
			Expect(groupType.ID).To(Equal(6))
			Expect(groupType.Name).To(Equal("Dynamische Gruppe"))
		})

		It("returns an error if the request cannot be send", func() {
			httpClient.DoReturns(nil, errors.New("request failed"))

			groupEndpoint := rest.NewGroupEndpoint(httpClient)
			_, err := groupEndpoint.GetGroupType(6)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("failed to send request, request failed"))
		})

		It("returns an error if the status code is wrong", func() {
			httpResponse := &http.Response{
				StatusCode: 404,
				Body: io.NopCloser(testutil.JsonToBufferString(
					`{
						"data": [],
						"meta": { "count": 0 }
					}`)),
			}
			httpClient.DoReturns(httpResponse, nil)

			groupEndpoint := rest.NewGroupEndpoint(httpClient)
			_, err := groupEndpoint.GetGroupType(6)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("received non-200 response code: 404"))
		})

		It("returns an error if the response body is not a church tools json response", func() {
			httpResponse := &http.Response{
				StatusCode: 200,
				Body: io.NopCloser(testutil.JsonToBufferString(
					`{
						"foo": [],
					}`)),
			}
			httpClient.DoReturns(httpResponse, nil)

			groupEndpoint := rest.NewGroupEndpoint(httpClient)
			_, err := groupEndpoint.GetGroupType(6)


			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("response body is not containing expected json"))
		})
	})
})
