package restendpoints

type GroupsResponseJson struct {
	Data []GroupsResponseItem `json:"data"`
}

type GroupsResponseItem struct {
	ID   int    `json:"id"`
	GUID string `json:"guid"`
	Name string `json:"name"`
}

type GroupsMembersResponseJson struct {
	Data []GroupsMembersResponseItem `json:"data"`
}

type GroupsMembersResponseItem struct {
	PersonId          int    `json:"personId"`
	GroupId           int    `json:"groupId"`
	GroupTypeRoleId   int    `json:"groupTypeRoleId"`
	GroupMemberStatus string `json:"groupMemberStatus"`
	Deleted           bool   `json:"deleted"`
}
