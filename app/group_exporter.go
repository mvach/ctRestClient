package app

import (
	"ctRestClient/rest"
	"encoding/json"
	"fmt"
)

type GroupName2IDMap map[string]int

//counterfeiter:generate . GroupExporter
type GroupExporter interface {
	ExportGroupMembers(
		groupName string,
		groupsEndpoint rest.GroupsEndpoint,
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
	personsEndpoint rest.PersonsEndpoint,
) ([]json.RawMessage, error) {
	var result []json.RawMessage

	ctGroup, err := groupsEndpoint.GetGroup(groupName)
	if err != nil {
		return nil, fmt.Errorf("failed to get group by name: %v", err)
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
