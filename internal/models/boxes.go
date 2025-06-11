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
	Created time.Time
	Expires time.Time
}

type BoxModel struct {
	DB *pgxpool.Pool
}

func (m *BoxModel) Insert(title string, content string, expires int) (int, error) {
	query := fmt.Sprintf("INSERT INTO boxes (title, content, created, expires) VALUES ('%s', '%s', now(), now() + interval '%d days') RETURNING id", title, content, expires)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var id int
	err := m.DB.QueryRow(ctx, query).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (m *BoxModel) Get(id int) (Box, error) {
	query := fmt.Sprintf("SELECT id, title, content, created, expires FROM boxes WHERE id = %d", id)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var box Box
	err := m.DB.QueryRow(ctx, query).Scan(&box.ID, &box.Title, &box.Content, &box.Created, &box.Expires)
	if err != nil {
		return Box{}, err
	}

	return box, nil
}

func (m *BoxModel) Latest() ([]Box, error) {
	query := "SELECT id, title, content, created, expires FROM boxes ORDER BY created DESC LIMIT 10"
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var boxes []Box
	rows, err := m.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var box Box
		err := rows.Scan(&box.ID, &box.Title, &box.Content, &box.Created, &box.Expires)
		if err != nil {
			return nil, err
		}
		boxes = append(boxes, box)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return boxes, nil
}
