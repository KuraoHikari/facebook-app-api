package dto

type PostUpdateDTO struct {
	ID          uint64 `json:"id" form:"id" binding:"required"`
	ImageLink   string `json:"image_link" form:"image_link"`
	Description string `json:"description" form:"description" binding:"required"`
	UserID      uint64 `json:"user_id,omitempty"  form:"user_id,omitempty"`
}

type PostCreateDTO struct {
	ImageLink   string `json:"image_link" form:"image_link"`
	Description string `json:"description" form:"description" binding:"required"`
	UserID      uint64 `json:"user_id,omitempty"  form:"user_id,omitempty"`
}