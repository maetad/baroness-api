package stateservice

type StateCreateRequest struct {
	WorkflowID uint   `json:"-"`
	Name       string `json:"name" binding:"required"`
}

type StateUpdateRequest struct {
	Name string `json:"name" binding:"required"`
}
