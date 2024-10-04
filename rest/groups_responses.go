package rest

type GroupsResponseJson struct {
    Data []GroupsResponse `json:"data"`
}

type GroupsResponse struct {
    ID   int    `json:"id"`
    GUID string `json:"guid"`
    Name string `json:"name"`
}

type GroupsMembersResponseJson struct {
    Data []GroupsMembersResponse `json:"data"`
}

type GroupsMembersResponse struct {
    PersonId          int    `json:"personId"`
    GroupId           int    `json:"groupId"`
    GroupTypeRoleId   int    `json:"groupTypeRoleId"`
    GroupMemberStatus string `json:"groupMemberStatus"`
    Deleted           bool   `json:"deleted"`
}
