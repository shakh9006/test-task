package models

type Number struct {
	ID     string `json:"id,omitempty" binding:"required"`
	Number string `json:"number" binding:"required"`
}
