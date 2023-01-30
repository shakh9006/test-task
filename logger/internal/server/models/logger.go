package models

type Logger struct {
	ID      string `json:"id,omitempty" binding:"required"`
	Type    string `json:"type" binding:"required"`
	Message string `json:"message" binding:"required"`
}
