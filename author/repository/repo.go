package repository

import "github.com/phantomnat/go-clean-architecture/author"

type AuthorRepository interface {
	GetByID(id int64) (*author.Author, error)
}
