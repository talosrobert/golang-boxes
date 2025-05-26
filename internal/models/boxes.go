package models

import (
	"context"
	"fmt"
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

func (m *BoxModel) Insert(title string, content string, expires int) (int64, error) {
	query := fmt.Sprintf("INSERT INTO boxes (title, content, created, expires) VALUES ('%s', '%s', now(), now() + interval '%d days')", title, content, expires)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	tag, err := m.DB.Exec(ctx, query)
	if err != nil {
		return 0, err
	}

	return tag.RowsAffected(), nil
}

func (m *BoxModel) Get(ID uint) (Box, error) {
	return Box{}, nil
}

func (m *BoxModel) Latest() ([]Box, error) {
	return nil, nil
}
