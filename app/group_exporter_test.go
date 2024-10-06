package app_test

import (
    "ctRestClient/app"
    "ctRestClient/restendpoints"
    "ctRestClient/restendpoints/restendpointsfakes"
    "encoding/json"
    "errors"

    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
)

var _ = Describe("GroupExporter", func() {

    var (
        dynamicGroupsEndpoint *restendpointsfakes.FakeDynamicGroupsEndpoint
        groupsEndpoint        *restendpointsfakes.FakeGroupsEndpoint
        personsEndpoint       *restendpointsfakes.FakePersonsEndpoint
    )

    BeforeEach(func() {
        dynamicGroupsEndpoint = &restendpointsfakes.FakeDynamicGroupsEndpoint{}
        groupsEndpoint = &restendpointsfakes.FakeGroupsEndpoint{}
        personsEndpoint = &restendpointsfakes.FakePersonsEndpoint{}

        dynamicGroupsEndpoint.GetDynamicGroupIdsReturns([]int{1, 2}, nil)
        groupsEndpoint.GetGroupNameReturns(
            []restendpoints.GroupsResponse{
                {ID: 1, Name: "foo_group"},
                {ID: 1, Name: "bar_group"},
            }, nil,
        )
        groupsEndpoint.GetGroupMembersReturns(
            []restendpoints.GroupsMembersResponse{
                {PersonId: 1, GroupId: 1},
                {PersonId: 2, GroupId: 1},
            }, nil,
        )
    })

    var _ = Describe("ExportPersonData", func() {
        It("returns persons", func() {
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
            personsEndpoint.GetPersonReturnsOnCall(0, json.RawMessage(person1), nil)
            personsEndpoint.GetPersonReturnsOnCall(1, json.RawMessage(person2), nil)

            personData, err := app.NewGroupExporter().ExportPersonData(
                "foo_group_name",
                dynamicGroupsEndpoint,
                groupsEndpoint,
                personsEndpoint)

            Expect(err).NotTo(HaveOccurred())
            Expect(personData).To(HaveLen(2))

            Expect(personData[0]).To(MatchJSON(person1))
            Expect(personData[1]).To(MatchJSON(person2))
        })

        It("returns an error if dynamic groups cannot be resolved", func() {
            dynamicGroupsEndpoint.GetDynamicGroupIdsReturns(nil, errors.New("boom"))

            personData, err := app.NewGroupExporter().ExportPersonData(
                "foo_group_name",
                dynamicGroupsEndpoint,
                groupsEndpoint,
                personsEndpoint)

            Expect(err.Error()).To(Equal("failed to resolve groupnames to ids, failed to get dynamic groups, boom"))
            Expect(personData).To(BeNil())
        })

        It("returns an error if dynamic groups are empty", func() {
            dynamicGroupsEndpoint.GetDynamicGroupIdsReturns([]int{}, nil)

            personData, err := app.NewGroupExporter().ExportPersonData(
                "foo_group_name",
                dynamicGroupsEndpoint,
                groupsEndpoint,
                personsEndpoint)

            Expect(err.Error()).To(Equal("failed to resolve groupnames to ids, no dynamic groups found"))
            Expect(personData).To(BeNil())
        })

        It("returns an error if group names cannot be resolved", func() {
            groupsEndpoint.GetGroupNameReturns(nil, errors.New("boom"))
                
            personData, err := app.NewGroupExporter().ExportPersonData(
                "foo_group_name",
                dynamicGroupsEndpoint,
                groupsEndpoint,
                personsEndpoint)

            Expect(err.Error()).To(Equal("failed to resolve groupnames to ids, failed to retrieve group name, boom"))
            Expect(personData).To(BeNil())
        })

        It("returns an error if group names are empty", func() {
            groupsEndpoint.GetGroupNameReturns(
                []restendpoints.GroupsResponse{}, nil,
            )
                
            personData, err := app.NewGroupExporter().ExportPersonData(
                "foo_group_name",
                dynamicGroupsEndpoint,
                groupsEndpoint,
                personsEndpoint)

            Expect(err.Error()).To(Equal("failed to resolve groupnames to ids, no group name found"))
            Expect(personData).To(BeNil())
        })

        It("returns an error if group members cannot be resolved", func() {
            groupsEndpoint.GetGroupMembersReturns(nil, errors.New("boom"))
                
            personData, err := app.NewGroupExporter().ExportPersonData(
                "foo_group_name",
                dynamicGroupsEndpoint,
                groupsEndpoint,
                personsEndpoint)

            Expect(err.Error()).To(Equal("failed to resolve group members, boom"))
            Expect(personData).To(BeNil())
        })

        It("returns an error if person cannot be resolved", func() {
            personsEndpoint.GetPersonReturnsOnCall(0, nil, errors.New("boom"))

            personData, err := app.NewGroupExporter().ExportPersonData(
                "foo_group_name",
                dynamicGroupsEndpoint,
                groupsEndpoint,
                personsEndpoint)

            Expect(err.Error()).To(Equal("failed to resolve person with id 1, boom"))
            Expect(personData).To(BeNil())
        })
    })
})
