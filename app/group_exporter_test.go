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
        dynamicGroupsEndpoint *restfakes.FakeDynamicGroupsEndpoint
        groupsEndpoint        *restfakes.FakeGroupsEndpoint
        personsEndpoint       *restfakes.FakePersonsEndpoint
        groupExporter         app.GroupExporter
    )

    BeforeEach(func() {
        dynamicGroupsEndpoint = &restfakes.FakeDynamicGroupsEndpoint{}
        groupsEndpoint = &restfakes.FakeGroupsEndpoint{}
        personsEndpoint = &restfakes.FakePersonsEndpoint{}

        groupExporter = app.NewGroupExporter()
    })

    var _ = Describe("ExportPersonData", func() {
        BeforeEach(func() {
            groupsEndpoint.GetGroupMembersReturns(
                []rest.GroupsMembersResponse{
                    {PersonId: 1, GroupId: 1},
                    {PersonId: 2, GroupId: 1},
                }, nil,
            )
        })

        
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
            personsEndpoint.GetPersonReturnsOnCall(0, []json.RawMessage{json.RawMessage(person1)}, nil)
            personsEndpoint.GetPersonReturnsOnCall(1, []json.RawMessage{json.RawMessage(person2)}, nil)

            personData, err := groupExporter.ExportPersonData(
                1,
                groupsEndpoint,
                personsEndpoint,
            )

            Expect(err).NotTo(HaveOccurred())
            Expect(personData).To(HaveLen(2))

            Expect(personData[0]).To(MatchJSON(person1))
            Expect(personData[1]).To(MatchJSON(person2))
        })

        It("returns an error if group members cannot be resolved", func() {
            groupsEndpoint.GetGroupMembersReturns(nil, errors.New("boom"))

            personData, err := groupExporter.ExportPersonData(
                1,
                groupsEndpoint,
                personsEndpoint,
            )

            Expect(err.Error()).To(Equal("failed to resolve group members, boom"))
            Expect(personData).To(BeNil())
        })

        It("returns an error if person cannot be resolved", func() {
            personsEndpoint.GetPersonReturnsOnCall(0, nil, errors.New("boom"))

            personData, err := groupExporter.ExportPersonData(
                1,
                groupsEndpoint,
                personsEndpoint,
            )

            Expect(err.Error()).To(Equal("failed to resolve person with id 1, boom"))
            Expect(personData).To(BeNil())
        })
    })

    var _ = Describe("GetGroupNames2IDMapping", func() {
        BeforeEach(func() {
            dynamicGroupsEndpoint.GetDynamicGroupIdsReturns([]int{1, 2}, nil)
    
            groupsEndpoint.GetGroupNamesReturns(
                []rest.GroupsResponse{
                    {ID: 1, Name: "foo_group"},
                    {ID: 2, Name: "bar_group"},
                }, nil,
            )
        })

        It("returns group id to name mapping", func() {
            group2IDMap, err := groupExporter.GetGroupNames2IDMapping(
                dynamicGroupsEndpoint,
                groupsEndpoint,
            )

            Expect(err).NotTo(HaveOccurred())
            Expect(group2IDMap["foo_group"]).To(Equal(1))
            Expect(group2IDMap["bar_group"]).To(Equal(2))
        })

        It("returns an error if dynamic groups cannot be resolved", func() {
            dynamicGroupsEndpoint.GetDynamicGroupIdsReturns(nil, errors.New("boom"))

            personData, err := groupExporter.GetGroupNames2IDMapping(
                dynamicGroupsEndpoint,
                groupsEndpoint,
            )

            Expect(err.Error()).To(Equal("failed to get dynamic groups, boom"))
            Expect(personData).To(BeNil())
        })

        It("returns an error if dynamic groups are empty", func() {
            dynamicGroupsEndpoint.GetDynamicGroupIdsReturns([]int{}, nil)

            personData, err := groupExporter.GetGroupNames2IDMapping(
                dynamicGroupsEndpoint,
                groupsEndpoint,
            )

            Expect(err.Error()).To(Equal("no dynamic groups found"))
            Expect(personData).To(BeNil())
        })

        It("returns an error if group names cannot be resolved", func() {
            groupsEndpoint.GetGroupNamesReturns(nil, errors.New("boom"))

            personData, err := groupExporter.GetGroupNames2IDMapping(
                dynamicGroupsEndpoint,
                groupsEndpoint,
            )

            Expect(err.Error()).To(Equal("failed to retrieve group name, boom"))
            Expect(personData).To(BeNil())
        })

        It("returns an error if group names are empty", func() {
            groupsEndpoint.GetGroupNamesReturns(
                []rest.GroupsResponse{}, nil,
            )

            personData, err := groupExporter.GetGroupNames2IDMapping(
                dynamicGroupsEndpoint,
                groupsEndpoint,
            )

            Expect(err.Error()).To(Equal("no group name found"))
            Expect(personData).To(BeNil())
        })

    })
})
