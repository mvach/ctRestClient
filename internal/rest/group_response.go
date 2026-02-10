package rest

type GroupTypeResponseJson struct {
	Data GroupTypeResponse `json:"data"`
}

type GroupTypeResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (r *GroupTypeResponse) IsDynamicGroup() bool {
	return r.Name == "Dynamische Gruppe"
}
