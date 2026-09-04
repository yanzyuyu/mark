package bookmark

import "time"

type Bookmark struct {
	ID        string    `json:"id"`
	Label     string    `json:"label"`
	File      string    `json:"file"`
	Line      int       `json:"line"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	Author    string    `json:"author"`
}
