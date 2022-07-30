package transitionservice

type TransitionCreateRequest struct {
	Name     string `json:"name" binding:"required"`
	ParentID uint   `json:"parent_id"`
	TargetID uint   `json:"target_id"`
}

type TransitionUpdateRequest struct {
	Name     string `json:"name" binding:"required"`
	ParentID uint   `json:"parent_id"`
	TargetID uint   `json:"target_id"`
}
