package app_test

import (
	"ctRestClient/app"
	"ctRestClient/rest"
	"ctRestClient/rest/restfakes"
	"encoding/json"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("GroupExporter", func() {

	var (
		groupsEndpoint        *restfakes.FakeGroupsEndpoint
		dynamicGroupsEndpoint *restfakes.FakeDynamicGroupsEndpoint
		personsEndpoint       *restfakes.FakePersonsEndpoint
		groupExporter         app.GroupExporter
	)

	BeforeEach(func() {
		groupsEndpoint = &restfakes.FakeGroupsEndpoint{}
		dynamicGroupsEndpoint = &restfakes.FakeDynamicGroupsEndpoint{}
		personsEndpoint = &restfakes.FakePersonsEndpoint{}

		groupExporter = app.NewGroupExporter()
	})

	var _ = Describe("ExportPersonData", func() {
		BeforeEach(func() {
			groupsEndpoint.GetGroupReturns(
				rest.GroupsResponse{ID: 1, GUID: "1234", Name: "group1"}, nil,
			)
			groupsEndpoint.GetGroupMembersReturns(
				[]rest.GroupsMembersResponse{
					{PersonId: 1, GroupId: 1},
					{PersonId: 2, GroupId: 1},
				}, nil,
			)
		})

		It("returns persons", func() {
			dynamicGroupsEndpoint.GetGroupStatusReturns(
				rest.DynamicGroupsStatusResponse{Status: ptr("active")}, nil,
			)

			person1 := `{	
                "id": 1,
                "firstName": "foo_firstname",
                "lastName": "foo_lastname"
            }`
			person2 := `{	
                "id": 2,
                "firstName": "bar_firstname",
                "lastName": "bar_lastname"
            }`
			personsEndpoint.GetPersonReturnsOnCall(0, []json.RawMessage{json.RawMessage(person1)}, nil)
			personsEndpoint.GetPersonReturnsOnCall(1, []json.RawMessage{json.RawMessage(person2)}, nil)

			personData, err := groupExporter.ExportGroupMembers(
				"group1",
				groupsEndpoint,
				dynamicGroupsEndpoint,
				personsEndpoint,
			)

			Expect(err).NotTo(HaveOccurred())
			Expect(personData).To(HaveLen(2))

			Expect(personData[0]).To(MatchJSON(person1))
			Expect(personData[1]).To(MatchJSON(person2))
		})

		var _ = Context("group is a dynamic group", func() {

			BeforeEach(func() {
				// The id 1 of the group "group1" is contained in the list of dynamic group IDs
				// thus "group1" is a dynamic group
				dynamicGroupsEndpoint.GetAllDynamicGroupsReturns(
					rest.DynamicGroupsResponse{GroupIDs: []int{1}}, nil,
				)
			})

			It("returns an error if fetching the dynamic group status fails", func() {
				dynamicGroupsEndpoint.GetGroupStatusReturns(
					rest.DynamicGroupsStatusResponse{}, errors.New("boom"),
				)

				personData, err := groupExporter.ExportGroupMembers(
					"group1",
					groupsEndpoint,
					dynamicGroupsEndpoint,
					personsEndpoint,
				)

				Expect(err.Error()).To(Equal("failed to get dynamic group status, boom"))
				Expect(personData).To(BeNil())
			})

			It("returns an error if the dynamic group status is not active", func() {
				dynamicGroupsEndpoint.GetGroupStatusReturns(
					rest.DynamicGroupsStatusResponse{Status: ptr("not-active")}, nil,
				)

				personData, err := groupExporter.ExportGroupMembers(
					"group1",
					groupsEndpoint,
					dynamicGroupsEndpoint,
					personsEndpoint,
				)

				Expect(err.Error()).To(Equal("dynamic group 'group1' is not active"))
				Expect(personData).To(BeNil())
			})
		})

		It("returns an error if group members cannot be resolved", func() {
			dynamicGroupsEndpoint.GetGroupStatusReturns(
				rest.DynamicGroupsStatusResponse{Status: ptr("active")}, nil,
			)

			groupsEndpoint.GetGroupMembersReturns(nil, errors.New("boom"))

			personData, err := groupExporter.ExportGroupMembers(
				"group1",
				groupsEndpoint,
				dynamicGroupsEndpoint,
				personsEndpoint,
			)

			Expect(err.Error()).To(Equal("failed to resolve group members, boom"))
			Expect(personData).To(BeNil())
		})

		It("returns an error if person cannot be resolved", func() {
			dynamicGroupsEndpoint.GetGroupStatusReturns(
				rest.DynamicGroupsStatusResponse{Status: ptr("active")}, nil,
			)

			personsEndpoint.GetPersonReturnsOnCall(0, nil, errors.New("boom"))

			personData, err := groupExporter.ExportGroupMembers(
				"group1",
				groupsEndpoint,
				dynamicGroupsEndpoint,
				personsEndpoint,
			)

			Expect(err.Error()).To(Equal("failed to resolve person with id 1, boom"))
			Expect(personData).To(BeNil())
		})
	})
})
