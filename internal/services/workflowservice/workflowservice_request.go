package workflowservice

type WorkflowCreateRequest struct {
	EventID uint   `json:"-"`
	Name    string `json:"name" binding:"required"`
}

type WorkflowUpdateRequest struct {
	Name string `json:"name" binding:"required"`
}
