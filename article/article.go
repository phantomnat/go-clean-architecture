package article

import (
	"time"

	"github.com/phantomnat/go-clean-architecture/author"
)

type Article struct {
	ID        int64         `json:"id"`
	Title     string        `json:"title" validate:"required"`
	Content   string        `json:"content" validate:"required"`
	Author    author.Author `json:"author"`
	UpdatedAt time.Time     `json:"updated_at"`
	CreatedAt time.Time     `json:"created_at"`
}
