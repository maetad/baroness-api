package userservice

type UserCreateRequest struct {
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required"`
	DisplayName string `json:"display_name" binding:"required"`
}

type UserUpdateRequest struct {
	Password    string `json:"password"`
	DisplayName string `json:"display_name" binding:"required"`
}
