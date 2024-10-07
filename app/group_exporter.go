package app

import (
    "ctRestClient/rest"
    "encoding/json"
    "fmt"
)

type GroupName2IDMap map[string]int

//counterfeiter:generate . GroupExporter
type GroupExporter interface {
    ExportPersonData(
        groupname string,
        dynamicGroupsEndpoint rest.DynamicGroupsEndpoint,
        groupsEndpoint rest.GroupsEndpoint,
        personsEndpoint rest.PersonsEndpoint,
    ) ([]json.RawMessage, error)
}

type groupExporter struct {
}

func NewGroupExporter() GroupExporter {
    return groupExporter{}
}

func (g groupExporter) ExportPersonData(
    groupname string,
    dynamicGroupsEndpoint rest.DynamicGroupsEndpoint,
    groupsEndpoint rest.GroupsEndpoint,
    personsEndpoint rest.PersonsEndpoint,
) ([]json.RawMessage, error) {

    var result []json.RawMessage

    groupName2IDMap, err := g.resolveGroupNames2IDs(dynamicGroupsEndpoint, groupsEndpoint)
    if err != nil {
        return nil, fmt.Errorf("failed to resolve groupnames to ids, %w", err)
    }

    groupMembers, err := groupsEndpoint.GetGroupMembers(groupName2IDMap[groupname])
    if err != nil {
        return nil, fmt.Errorf("failed to resolve group members, %w", err)
    }

    for _, groupMember := range groupMembers {
        personsJson, err := personsEndpoint.GetPerson(groupMember.PersonId)
        if err != nil {
            return nil, fmt.Errorf("failed to resolve person with id %d, %w", groupMember.PersonId, err)
        }

        result = append(result, personsJson...)
    }

    return result, nil
}

func (g groupExporter) resolveGroupNames2IDs(
    dynamicGroupsEndpoint rest.DynamicGroupsEndpoint,
    groupsEndpoint rest.GroupsEndpoint,
) (GroupName2IDMap, error) {
    groupName2IDMap := make(GroupName2IDMap)

    groupIds, err := dynamicGroupsEndpoint.GetDynamicGroupIds()
    if err != nil {
        return nil, fmt.Errorf("failed to get dynamic groups, %w", err)
    }

    if len(groupIds) == 0 {
        return nil, fmt.Errorf("no dynamic groups found")
    }

    for _, groupId := range groupIds {
        groupsResponse, err := groupsEndpoint.GetGroupName(groupId)
        if err != nil {
            return nil, fmt.Errorf("failed to retrieve group name, %w", err)
        }

        if len(groupsResponse) == 0 {
            return nil, fmt.Errorf("no group name found")
        }

        for _, group := range groupsResponse {
            groupName2IDMap[group.Name] = group.ID
        }
    }

    return groupName2IDMap, nil
}
