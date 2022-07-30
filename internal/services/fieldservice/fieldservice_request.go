package fieldservice

type FieldCreateRequest struct {
	WorkflowID uint   `json:"-"`
	Name       string `json:"name" binding:"required"`
}

type FieldUpdateRequest struct {
	Name string `json:"name" binding:"required"`
}
