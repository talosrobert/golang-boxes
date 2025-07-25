package models

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	Id      uuid.UUID
	Name    string
	Email   string
	Pswhash string
	Created time.Time
}

type UserModel struct {
	DB *pgxpool.Pool
}

func (m *UserModel) Insert(name string, email string, psw string) error {
	query := fmt.Sprintf("INSERT INTO users (name, email, pswhash) VALUES ('%s', '%s', crypt('%s', gen_salt('bf'))) RETURNING id", name, email, psw)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var id uuid.UUID
	err := m.DB.QueryRow(ctx, query).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr); pgErr.Code == "23505" && strings.Contains(pgErr.Message, "users_email_key") {
			return ErrDuplicateEmail
		}
		return err
	}

	return nil
}

func (m *UserModel) Get(id uuid.UUID) (User, error) {
	query := fmt.Sprintf("SELECT id, name, email, created FROM users WHERE id = '%s'", id)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var user User
	err := m.DB.QueryRow(ctx, query).Scan(&user.Id, &user.Name, &user.Email, &user.Created)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return user, ErrNoRows
		}
		return user, err
	}

	return user, nil
}

func (m *UserModel) Authenticate(email string, psw string) (uuid.UUID, error) {
	query := fmt.Sprintf("SELECT id FROM users WHERE email = '%s' AND pswhash = crypt('%s', pswhash);", email, psw)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var id uuid.UUID
	err := m.DB.QueryRow(ctx, query).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return id, ErrWrongCredentials
		} else {
			return id, err
		}
	}
	return id, nil
}

func (m *UserModel) UpdatePsw(email string, psw string) (uuid.UUID, error) {
	// https://www.postgresql.org/docs/16/pgcrypto.html#PGCRYPTO-PASSWORD-HASHING-FUNCS
	// SELECT (pswhash = crypt('entered password', pswhash)) AS pswmatch FROM ... ;
	var id uuid.UUID
	return id, nil
}
