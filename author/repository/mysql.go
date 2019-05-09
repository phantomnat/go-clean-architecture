package repository

import (
	"database/sql"

	"github.com/phantomnat/go-clean-architecture/author"
	"github.com/sirupsen/logrus"
)

type MysqlAuthorRepo struct {
	DB *sql.DB
}

func NewMysqlAuthorRepository(db *sql.DB) AuthorRepository {
	return &MysqlAuthorRepo{DB: db}
}

func (m *MysqlAuthorRepo) getOne(query string, args ...interface{}) (*author.Author, error) {
	stmt, err := m.DB.Prepare(query)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	row := stmt.QueryRow(args...)
	a := &author.Author{}

	err = row.Scan(
		&a.ID,
		&a.Name,
		&a.CreatedAt,
		&a.UpdatedAt)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return a, nil
}

func (m *MysqlAuthorRepo) GetByID(id int64) (*author.Author, error) {
	query := `SELECT id, name, created_at, updated_at FROM author WHERE id=?`
	return m.getOne(query, id)
}
