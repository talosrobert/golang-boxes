package models

import (
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Box struct {
	ID      uint
	Title   string
	Content string
	created time.Time
	expires time.Time
}

type BoxModel struct {
	DB *pgxpool.Pool
}

func (m *BoxModel) insert(title, content string, expires int) (int, error) {
	return 0, nil
}

func (m *BoxModel) get(ID uint) (Box, error) {
	return Box{}, nil
}

func (m *BoxModel) latest() ([]Box, error) {
	return nil, nil
}
