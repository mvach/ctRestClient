package rest

type DynamicGroupsStatusResponse struct {
	Status *string `json:"dynamicGroupStatus"`
}

type DynamicGroupsResponse struct {
	GroupIDs []int `json:"data"`
}
