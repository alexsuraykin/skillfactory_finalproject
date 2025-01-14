package models

type Feeds struct {
	Id      int    `json:"id" db:"id" bson:"id, omitempty"`
	Title   string `json:"title" db:"title"`
	Content string `json:"content" db:"content"`
	Link    string `json:"link" db:"link"`
	PubDate string `json:"pub_date" db:"pub_date"`
}

type Comments struct {
	ID              int32   `json:"ID"`
	NewsID          *int32  `json:"NewsID"`
	ParentCommentID *int32  `json:"ParentCommentID"`
	Content         string  `json:"Content"`
	CreatedAt       *string `json:"CreatedAt"`
}

type FeedsById struct {
	Feeds    Feeds      `json:"feeds"`
	Comments []Comments `json:"comments"`
}
