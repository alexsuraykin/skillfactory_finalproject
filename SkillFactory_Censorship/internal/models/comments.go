package models

type Comments struct {
	ID              int32   `json:"ID"`
	NewsID          *int32  `json:"NewsID"`
	ParentCommentID *int32  `json:"ParentCommentID"`
	Content         string  `json:"Content"`
	CreatedAt       *string `json:"CreatedAt"`
}
