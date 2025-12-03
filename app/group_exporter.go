package app

import (
	"ctRestClient/rest"
	"encoding/json"
	"fmt"
	"slices"
)

type GroupName2IDMap map[string]int

//counterfeiter:generate . GroupExporter
type GroupExporter interface {
	ExportGroupMembers(
		groupName string,
		groupsEndpoint rest.GroupsEndpoint,
		dynamicGroupsEndpoint rest.DynamicGroupsEndpoint,
		personsEndpoint rest.PersonsEndpoint,
	) ([]json.RawMessage, error)
}

type groupExporter struct {
}

func NewGroupExporter() GroupExporter {
	return groupExporter{}
}

func (g groupExporter) ExportGroupMembers(
	groupName string,
	groupsEndpoint rest.GroupsEndpoint,
	dynamicGroupsEndpoint rest.DynamicGroupsEndpoint,
	personsEndpoint rest.PersonsEndpoint,
) ([]json.RawMessage, error) {
	var result []json.RawMessage

	ctGroup, err := groupsEndpoint.GetGroup(groupName)
	if err != nil {
		return nil, fmt.Errorf("failed to get group by name: %v", err)
	}

	dynamicGroupsResponse, err := dynamicGroupsEndpoint.GetAllDynamicGroups()
	if err != nil {
		return nil, fmt.Errorf("failed to get all dynamic groups, %w", err)
	}

	if slices.Contains(dynamicGroupsResponse.GroupIDs, ctGroup.ID) {
		dynamicGroup, err := dynamicGroupsEndpoint.GetGroupStatus(ctGroup.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get dynamic group status, %w", err)
		}

		if *dynamicGroup.Status != "active" {
			return nil, &GroupNotActiveError{GroupName: groupName}
		}
	}

	groupMembers, err := groupsEndpoint.GetGroupMembers(ctGroup.ID)
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
