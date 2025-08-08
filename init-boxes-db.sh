#!/usr/bin/env bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
	CREATE TABLE boxes (
	    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	    title varchar(100) NOT NULL,
	    content text NOT NULL,
	    created timestamp NOT NULL DEFAULT now(),
	    expires timestamp NOT NULL
	);

	CREATE INDEX boxes_idx ON boxes(created);

	INSERT INTO boxes (title, content, expires) VALUES 
	    ('box1', 'aaaaaaa', now() + interval '365 days'),
	    ('box2', 'bbbbbbb', now() + interval '365 days'),
	    ('box3', 'ccccccc', now() + interval '365 days');

	CREATE TABLE sessions (
		token text PRIMARY KEY,
		data bytea NOT NULL,
		expiry timestamp NOT NULL
	);

	CREATE INDEX sessions_expiry_idx ON sessions (expiry);

	CREATE EXTENSION pgcrypto;

	CREATE TABLE users (
	    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	    name varchar(100) NOT NULL,
	    email text NOT NULL UNIQUE,
	    pswhash text NOT NULL, 
	    created timestamp NOT NULL DEFAULT now()
	);
EOSQL
