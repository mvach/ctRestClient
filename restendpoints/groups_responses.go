package restendpoints

type GroupsResponseJson struct {
    Data []DataItem `json:"data"`
}

type DataItem struct {
    ID                    int       `json:"id"`
    GUID                  string    `json:"guid"`
    Name                  string    `json:"name"`
}
